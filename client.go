package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Client SecretsManager client interface
type Client struct {
	secretsmanageriface.SecretsManagerAPI
	secretName  string
	secretValue *secretsmanager.GetSecretValueOutput
}

// New instanciates a Client struct
func New() (*Client, error) {
	// Instanciate a new aws session
	awssession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Instanciate a new SecretsManager client with an aws session
	svc = secretsmanager.New(awssession)

	// Prepare input to retrieve secret's value string
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(*secretName),
	}

	// Retrieve secret's value'
	output, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	return &Client{svc, *secretName, output}, nil
}

// DescribeSecret will describe the given aws secret from AWS SecretsManager
// Minimum IAM permission required:
// * secretsmanager:DescribeSecret
func (c *Client) DescribeSecret() error {
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(*secretName),
	}
	output, err := svc.DescribeSecret(input)
	if err != nil {
		return err
	}
	log.Println(output)
	return nil
}

// GetSecretString will retrieve the value of the encrypted key of the given secret
//
// Minimum IAM permission:
//
// * secretsmanager:GetSecretValue
//
// * kms:Decrypt
//
//
// It returns a string containing encrypted value and any error encountered
func (c *Client) GetSecretString() string {
	return *c.secretValue.SecretString
}

// KeyExists will verify whether the given key exists or not
// It returns false if the key doesn't exists, and true if the key exists
func (c *Client) KeyExists() bool {
	secretString := c.GetSecretString()

	// Lookup if the key already exists.
	// gjson.Get() return the value if the key exists
	return gjson.Get(secretString, *secretKey).Exists()
}

// UpdateSecretString will update the value of the SecretString of the given secret in json format
//
// Minimum IAM permission:
//
// * secretsmanager:UpdateSecret
//
//
// It returns any error encountered
func (c *Client) UpdateSecretString(payload string) error {
	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(*secretName),
		SecretString: &payload,
	}

	_, err := svc.UpdateSecret(input)
	if err != nil {
		return err
	}

	return nil
}

// ReadSecretKey will read the given key of a given secret
// It returns the unencrypted value of the given key and any error encountered
func (c *Client) ReadSecretKey() (string, error) {
	if !c.KeyExists() {
		return "", ErrKeyDoesNotExists
	}

	secretString := c.GetSecretString()

	return gjson.Get(secretString, *secretKey).String(), nil
}

// AddSecretKeyValue will add a new key/value to the given secret
// It returns any error encountered
func (c *Client) AddSecretKeyValue() error {
	// Lookup if the key already exists.
	if c.KeyExists() {
		return ErrKeyAlreadyExists
	}

	secretString := c.GetSecretString()

	// Build the new secret json by append the new key: value
	newSecretString, err := sjson.Set(secretString, *secretKey, *secretValue)
	if err != nil {
		return err
	}

	// Update the AWS Secret
	err = c.UpdateSecretString(newSecretString)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSecretValue will update the given value of a given key for a given secret
// It returns any error encountered
func (c *Client) UpdateSecretValue() error {
	secretString := c.GetSecretString()

	// Lookup if the key exists.
	if !c.KeyExists() {
		return ErrKeyDoesNotExists
	}

	// Build the new secret json by updating the key: value
	newSecretString, err := sjson.Set(secretString, *secretKey, *secretValue)
	if err != nil {
		return err
	}

	// Update the AWS Secret
	err = c.UpdateSecretString(newSecretString)
	if err != nil {
		return err
	}
	return nil
}

// RemoveSecretKeyValue will remove the given key of a given secret
// It returns any error encountered
func (c *Client) RemoveSecretKeyValue() error {
	secretString := c.GetSecretString()

	// Lookup if the key exists.
	if !c.KeyExists() {
		return ErrKeyDoesNotExists
	}

	// Build the new secret json by removing the key
	newSecretString, err := sjson.Delete(secretString, *secretKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Update the AWS Secret
	err = c.UpdateSecretString(newSecretString)
	if err != nil {
		return err
	}
	return nil
}
