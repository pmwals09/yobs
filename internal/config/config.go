package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("No .env file found. Loading from environment instead.")
	} else {
		fmt.Println("Loading from .env file")
	}
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3BucketName    string
}

func (a AWSConfig) IsIncomplete() bool {
	return a.Region == "" && a.AccessKeyID == "" && a.SecretAccessKey == "" && a.S3BucketName == ""
}

type Config struct {
	AWSConfig AWSConfig
}

func New() (*Config, error) {
	var allErr error
	awsConfig, err := newAWSConfig()
	allErr = errors.Join(allErr, err)
	if awsConfig.IsIncomplete() {
		allErr = errors.Join(allErr, errors.New("awsConfig is incomplete"))
	}
	c := Config{
		AWSConfig: awsConfig,
	}
	return &c, err
}

func newAWSConfig() (AWSConfig, error) {
	var allErr error
	var awsConfig AWSConfig
	var ok bool
	awsConfig.SecretAccessKey, ok = os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		allErr = errors.Join(allErr, errors.New("no AWS Access Key ID detected"))
	}
	awsConfig.Region, ok = os.LookupEnv("AWS_REGION")
	if !ok {
		allErr = errors.Join(allErr, errors.New("no AWS region detected"))
	}
	awsConfig.SecretAccessKey, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		allErr = errors.Join(allErr, errors.New("no AWS secret access key detected"))
	}
	awsConfig.S3BucketName, ok = os.LookupEnv("AWS_S3_BUCKET")
	if !ok {
		allErr = errors.Join(allErr, errors.New("no AWS S3 bucket name detected"))
	}

	return awsConfig, allErr
}
