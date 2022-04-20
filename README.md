# awssm-go

This repository contains a simple cli written in go to help update/add/remove/read a key or value of a key/value secret stored in [AWS SecretsManager](https://aws.amazon.com/secrets-manager/).

## Motivations

Since we are intensively using [AWS SecretsManager](https://aws.amazon.com/secrets-manager/) at Moodagent to store key/value micro services runtime's configuration, it can be very time consuming to add a new key or update the value of an existing one for every secrets, especially manually.
This simple cli aims to reduce the time spent in the console clicking "Edit secret". It is very well suitable in scripts or in one-liners command lines.

## Installation

### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go 1.14 only, but should work with go >= 1.11 ). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

```sh
# Clone this repository
git clone https://github.com/lescactus/awssm-go.git && cd awssm-go/

# Build from sources. Use the '-o' flag to change the compiled binary name
go build

# Default compiled binary is awssm-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./awssm-go -h
```

### From source with docker

If you don't have [go](https://golang.org/) installed but have docker, run the following command to build inside a docker container:

```sh
# Build from sources inside a docker container. Use the '-o' flag to change the compiled binary name
# Warning: the compiled binary belongs to root:root
docker run --rm -it -v "$PWD":/app -w /app golang:1.14 go build

# Default compiled binary is awssm-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./awssm-go -h
```

### From source with docker but built inside a docker image

If you don't want to pollute your computer with another program, with cli comes with its own docker image:

```sh
docker build -t awssm-go .
```

### From source with `go get`

You can also add this cli in your `$GOPATH` with `go get`:

```sh
# Download and build this module in $GOPATH.
# $GOPATH can be overriden like:
# GOPATH="$(pwd)/go" go get ...
go get -v github.com/lescactus/awssm-go/

# The compiled binary is installed in $GOPATH/bin (By default: $HOME/go)
ls -l $GOPATH/bin
total 12M
-rwxrwxr-x 1 amaldeme amaldeme 12M Dec 21 12:36 awssm-go
```

## Usage

As behind the scenes an AWS Secret is just a json structure, a key/value secret is simply:
```json
{
    "key1": "value1",
    "key2": "value2",
    ...
}
```
`awssm-go` will only works with simple key/value secrets. Support for more complex json structures (like arrays or nested structures) may be covered in the future.

---
```sh
Usage of ./awssm-go:
  -key string
    	Key to add/read/remove/update
  -op string
    	Operation. Can be one of the following: add/describe/read/remove/show/update
  -secret string
    	Secret name to describe/read/show/update
  -value string
    	Value of the key to add/read/update
```
---
`awssm-go` uses the [AWS SDK for Go API](https://docs.aws.amazon.com/sdk-for-go/api/) to access AWS services. Thus, passing credentials to this program is the same as the standard aws cli (https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)

You can use:
* Environment variables,
* Shared credentials file,
* AWS Profile,
* IAM instance role if running in EC2, ECS, EKS, etc ...,
* ...

Of course you need to have the proper IAM permissions to read/write on [AWS SecretsManager](https://aws.amazon.com/secrets-manager/).

### Usage with docker

If you've build the docker image embedding this cli, you need to pass your aws keys or credentials file to the container:

```sh
# Use environment variables
docker run --rm \
  -it \
  -e AWS_ACCESS_KEY_ID=xxxx \
  -e AWS_SECRET_ACCESS_KEY=xxxx \
  -e AWS_DEFAULT_REGION=xxxx
  awssm-go -secret <your secret> xxxxxx

# Mount your config & credentials files in the container
docker run \
  --rm \
  -it -v ~/.aws:/root/.aws \
  -e AWS_PROFILE=<your profile if needed> \
  awssm-go -secret <your secret> xxxxxx

```

### Examples

#### Describe a secret

```sh
./awssm-go -secret alma/test -op describe
2020/12/21 12:53:42 {
  ARN: "arn:aws:secretsmanager:eu-central-1:123456789123:secret:alma/test-IffD4v",
  CreatedDate: 2020-11-21 20:28:57.312 +0000 UTC,
  Description: "this is a test",
  LastAccessedDate: 2020-12-21 00:00:00 +0000 UTC,
  LastChangedDate: 2020-12-21 10:47:24.885 +0000 UTC,
  Name: "alma/test",
  Tags: [],
  VersionIdsToStages: {
    38445FC3-D1E7-4D22-8089-DA8804E5B6F8: ["AWSPREVIOUS"],
    C7614E92-5B2A-412B-9D44-D53BFAD9D111: ["AWSCURRENT"]
  }
}
```
### Read the value of a given key from a given secret

```sh
./awssm-go -secret alma/test -op read -key mykey
thevalueofthiskey
```

### Add a key to a given secret

```sh
./awssm-go -secret alma/test -op add -key theNewKey -value theValueOfThisNewKey
2020/12/21 12:58:31 alma/test updated

./awssm-go -secret alma/test -op read -key theNewKey
theValueOfThisNewKey
```

### Update the value of a given key of a given secret

```sh
./awssm-go -secret alma/test -op update -key theNewKey -value theNewValueOfThisKey
2020/12/21 12:59:44 alma/test updated

./awssm-go -secret alma/test -op read -key theNewKey 
theNewValueOfThisKey
```

### Remove a given key from a given secret

```sh
./awssm-go -secret alma/test -op remove -key theNewKey
2020/12/21 13:00:20 alma/test updated

./awssm-go -secret alma/test -op read -key theNewKey  
2020/12/21 13:00:22 Error: Key does not exists in secret "alma/test": "theNewKey"
```

### Show the whole secret

```sh
./awssm-go -secret alma/test -op show
2020/12/22 19:08:53 {"key":"azertyuiop","key2":"azertyuiopedede","key22":"ee","key222":"ee","key22200":"ee","00key22200":"ee","mykey2":"aaa","mykey":"azerty"}
```

## Roadmap

* Add tests
* Check whether the given secret is a simple key/value
* Support for more complex json structures