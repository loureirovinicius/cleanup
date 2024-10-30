package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elasticblockstorage "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticBlockStorage"
	elasticip "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticIp"
	elasticnetworkinterface "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticNetworkInterface"
	loadbalancer "github.com/loureirovinicius/cleanup/aws/service/ec2/loadBalancer"
	targetgroup "github.com/loureirovinicius/cleanup/aws/service/ec2/targetGroup"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/spf13/viper"
)

type AWS struct {
	config    Config
	client    *aws.Config
	Resources map[string]Cleanable
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
	logger.Log(ctx, "debug", "Loading AWS configurations...")
	err := p.loadConfig()
	if err != nil {
		return err
	}
	logger.Log(ctx, "debug", "AWS configs were loaded successfully!")

	logger.Log(ctx, "debug", "Creating AWS client...")
	err = p.createClient(ctx)
	if err != nil {
		return err
	}
	logger.Log(ctx, "debug", "AWS client was created successfully!")

	p.Resources = p.loadServices(ctx)

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

func (p *AWS) loadServices(ctx context.Context) map[string]Cleanable {
	ec2API := ec2.NewFromConfig(*p.client)
	logger.Log(ctx, "debug", "EC2 Configs initialized...")
	elbAPI := elasticloadbalancingv2.NewFromConfig(*p.client)
	logger.Log(ctx, "debug", "ELB Configs initialized...")

	return map[string]Cleanable{
		"targetGroup":  &targetgroup.TargetGroup{API: elbAPI},
		"loadBalancer": &loadbalancer.LoadBalancer{API: elbAPI},
		"eni":          &elasticnetworkinterface.ElasticNetworkInterface{API: ec2API},
		"eip":          &elasticip.ElasticIP{API: ec2API},
		"ebs":          &elasticblockstorage.ElasticBlockStorage{API: ec2API},
	}
}
