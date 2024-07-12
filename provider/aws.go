package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	elasticblockstorage "github.com/loureirovinicius/cleanup/aws/service/elasticBlockStorage"
	elasticip "github.com/loureirovinicius/cleanup/aws/service/elasticIp"
	elasticnetworkinterface "github.com/loureirovinicius/cleanup/aws/service/elasticNetworkInterface"
	loadbalancer "github.com/loureirovinicius/cleanup/aws/service/loadBalancer"
	targetgroup "github.com/loureirovinicius/cleanup/aws/service/targetGroup"
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

func (p *AWS) Initialize(ctx context.Context, cfg *ProviderConfig) error {
	err := p.loadConfig()
	if err != nil {
		return err
	}

	err = p.createClient(ctx)
	if err != nil {
		return err
	}

	p.Resources = map[string]cleaner.Cleanable{
		"targetGroup":  &targetgroup.TargetGroup{Client: p.client},
		"loadBalancer": &loadbalancer.LoadBalancer{Client: p.client},
		"eni":          &elasticnetworkinterface.ElasticNetworkInterface{Client: p.client},
		"eip":          &elasticip.ElasticIP{Client: p.client},
		"ebs":          &elasticblockstorage.ElasticBlockStorage{Client: p.client},
	}
	cfg.AWS = *p

	return nil
}

func (p *AWS) createClient(ctx context.Context) error {
	credentials := credentials.NewStaticCredentialsProvider(p.config.Credentials.AccessKey, p.config.Credentials.SecretKey, "")

	config, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(p.config.Region),
		config.WithSharedConfigFiles([]string{p.config.Profile.Path}),
		config.WithSharedConfigProfile(p.config.Profile.Name),
		config.WithCredentialsProvider(credentials),
	)
	if err != nil {
		return fmt.Errorf("error creating AWS client: %v", err)
	}

	p.client = &config

	return nil
}

func (p *AWS) loadConfig() error {
	region := viper.GetString("aws.region")
	if region == "" {
		return errors.New("AWS region can't be empty")
	}
	p.config.Region = region

	// This doesn't need to raise an error if empty because it's not a required configuration.
	profile := viper.GetStringMapString("aws.authentication.profile")
	if profile["name"] != "" || profile["path"] != "" {
		p.config.Profile = Profile{
			Name: profile["name"],
			Path: profile["path"],
		}
	}

	// This doesn't need to raise an error if empty because it's not a required configuration.
	credentials := viper.GetStringMapString("aws.authentication.credentials")
	if credentials["access_key"] != "" && credentials["secret_key"] != "" {
		p.config.Credentials = Credentials{
			AccessKey: credentials["access_key"],
			SecretKey: credentials["secret_key"],
		}
	}

	return nil
}
