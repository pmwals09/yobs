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

type awsConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3BucketName    string
}

func (a awsConfig) IsIncomplete() bool {
	return a.Region == "" || a.AccessKeyID == "" || a.SecretAccessKey == "" || a.S3BucketName == ""
}

type dbConfig struct {
	URL   string
	Token string
}

func (d dbConfig) IsIncomplete() bool {
	return d.URL == "" || d.Token == ""
}

type environmentName string

const (
	UnknownEnvironment    environmentName = ""
	LocalEnvironment      environmentName = "Local"
	ProductionEnvironment environmentName = "Production"
)

func EnvFromString(s string) (environmentName, error) {
	switch s {
	case "Local":
		return LocalEnvironment, nil
	case "Production":
		return ProductionEnvironment, nil
	default:
		return UnknownEnvironment, fmt.Errorf("unknown or missing environment name: %s", s)
	}
}

type Config struct {
	AWSConfig       awsConfig
	DBConfig        dbConfig
	EnvironmentName environmentName
}

func getEnvironmentName() (environmentName, error) {
	envRaw, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return LocalEnvironment, errors.New("no ENVIRONMENT detected")
	} else {
		return EnvFromString(envRaw)
	}
}

func New() (*Config, error) {
	var allErr error

	env, err := getEnvironmentName()
	if err != nil {
		allErr = errors.Join(allErr, err)
	}

	awsConfig, err := newAWSConfig()
	if err != nil {
		allErr = errors.Join(allErr, err)
	}
	if awsConfig.IsIncomplete() {
		allErr = errors.Join(allErr, errors.New("awsConfig is incomplete"))
	}

	dbConfig, err := newDBConfig()
	if err != nil {
		allErr = errors.Join(allErr, err)
	}
	// NOTE: We don't need an auth token and a url if we're developing locally
	if dbConfig.IsIncomplete() && env != LocalEnvironment {
		allErr = errors.Join()
	}

	c := Config{
		AWSConfig:       awsConfig,
		DBConfig:        dbConfig,
		EnvironmentName: env,
	}
	return &c, allErr
}

func newAWSConfig() (awsConfig, error) {
	var allErr error
	var awsConfig awsConfig
	var ok bool
	awsConfig.AccessKeyID, ok = os.LookupEnv("AWS_ACCESS_KEY_ID")
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

func newDBConfig() (dbConfig, error) {
	var allErr error
	var dbConfig dbConfig
	var ok bool
	dbConfig.URL, ok = os.LookupEnv("TURSO_DATABASE_URL")
	if !ok {
		allErr = errors.Join(allErr, errors.New("no Turso database url detected"))
	}
	dbConfig.Token, ok = os.LookupEnv("TURSO_AUTH_TOKEN")
	if !ok {
		allErr = errors.Join(allErr, errors.New("no Turso auth token detected"))
	}

	return dbConfig, allErr
}
