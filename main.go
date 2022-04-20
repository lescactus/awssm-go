package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	// Instanciate cli flags
	operation   = flag.String("op", "", "Operation. Can be one of the following: add/describe/read/remove/show/update")
	secretName  = flag.String("secret", "", "Secret name to describe/read/show/update")
	secretKey   = flag.String("key", "", "Key to add/read/remove/update")
	secretValue = flag.String("value", "", "Value of the key to add/read/update")

	// Instanciate a new aws session
	awssession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Instanciate a new SecretsManager client with an aws session
	svc = secretsmanager.New(awssession)
)

func init() {
	flag.Parse()

	if *secretName == "" {
		log.Fatalln(ErrFlagSecretMissing)
	}
}

func main() {

	// Create a new AWS SecretsManager client
	awssmClient, err := New()
	if err != nil {
		log.Fatalln(err)
	}

	switch *operation {

	case "add":
		if *secretKey == "" {
			log.Fatalln(ErrFlagKeyMissing)
		}
		if *secretValue == "" {
			log.Fatalln(ErrFlagValueMissing)
		}

		err := awssmClient.AddSecretKeyValue()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(*secretName, "updated")

	case "describe":
		err := awssmClient.DescribeSecret()
		if err != nil {
			log.Fatalln(err)
		}

	case "read":
		if *secretKey == "" {
			log.Fatalln(ErrFlagKeyMissing)
		}

		output, err := awssmClient.ReadSecretKey()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(output)

	case "remove":
		if *secretKey == "" {
			log.Fatalln(ErrFlagKeyMissing)
		}

		err := awssmClient.RemoveSecretKeyValue()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(*secretName, "updated")

	case "show":
		output := awssmClient.GetSecretString()
		log.Println(output)

	case "update":
		if *secretKey == "" {
			log.Fatalln(ErrFlagKeyMissing)
		}
		if *secretValue == "" {
			log.Fatalln(ErrFlagValueMissing)
		}

		err := awssmClient.UpdateSecretValue()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(*secretName, "updated")

	default:
		flag.PrintDefaults()
	}
}
