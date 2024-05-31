package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func GetPort() int {
	portStr := getEnvironmentValue("PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("Port: %s Incorrect!", portStr)
	}

	return port
}

func GetRegion() string {
	return getEnvironmentValue("REGION")
}

func GetAccessKey() string {
	return getEnvironmentValue("ACCESS_KEY")
}

func GetSecretAccessKey() string {
	return getEnvironmentValue("SECRET_ACCESS_KEY")
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s Environment Variables not Found or are Missing!", key)
	}

	return os.Getenv(key)
}

func GetAWSConfig() (aws.Config, error){
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(GetRegion()),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, Region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: "http://localhost:4566"}, nil
				},
			),
		),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     GetAccessKey(),
					SecretAccessKey: GetSecretAccessKey(),
					Source:          "",
				},
			},
		),
	)

	return cfg, err
}