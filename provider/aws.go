package provider

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/loureirovinicius/cleanup/aws/service/targetGroup"
	"github.com/loureirovinicius/cleanup/cmd/cleaner"
	"github.com/spf13/viper"
)

type AWS struct {
	config    Config
	client    *aws.Config
	Resources map[string]cleaner.Cleanable
}

type Config struct {
	// AWS Access and Secret Key credentials. DO NOT USE THIS FIELD FOR PRODUCTION PURPOSES.
	Credentials Credentials

	// AWS Configuration Profile to be used. DO NOT USE THIS FIELD FOR PRODUCTION PURPOSES.
	Profile Profile

	// AWS Region
	Region string
}

type Profile struct {
	// AWS Configuration Profile name.
	Name string

	// AWS Configuration Profile file path.
	Path string
}

type Credentials struct {
	// AWS User's Access Key.
	AccessKey string

	// AWS User's Secret Key.
	SecretKey string
}

func (p *AWS) Initialize(ctx context.Context, cfg *ProviderConfig) {
	p.loadConfig()
	p.createClient(ctx)
	p.Resources = map[string]cleaner.Cleanable{
		"targetGroup": &targetGroup.TargetGroup{Client: p.client},
	}
	cfg.AWS = *p
}

func (p *AWS) createClient(ctx context.Context) {
	credentials := credentials.NewStaticCredentialsProvider(p.config.Credentials.AccessKey, p.config.Credentials.SecretKey, "")

	config, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(p.config.Region),
		config.WithSharedConfigFiles([]string{p.config.Profile.Path}),
		config.WithSharedConfigProfile(p.config.Profile.Name),
		config.WithCredentialsProvider(credentials),
	)
	if err != nil {
		log.Fatalf("error loading AWS configs: %v", err)
	}

	p.client = &config
}

func (p *AWS) loadConfig() {
	region := viper.GetString("aws.region")
	if region == "" {
		log.Fatalln("AWS region can't be empty.")
	}
	p.config.Region = region

	profile := viper.GetStringMapString("aws.authentication.profile")
	if profile["name"] != "" || profile["path"] != "" {
		p.config.Profile = Profile{
			Name: profile["name"],
			Path: profile["path"],
		}
	}

	credentials := viper.GetStringMapString("aws.authentication.credentials")
	if credentials["access_key"] != "" && credentials["secret_key"] != "" {
		p.config.Credentials = Credentials{
			AccessKey: credentials["access_key"],
			SecretKey: credentials["secret_key"],
		}
	}
}