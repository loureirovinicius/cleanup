package providers

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
	config  Config
	Service Cleanable
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

// Call all the functions below in order to properly set the required configs for the cloud provider
func (p *AWS) Initialize(ctx context.Context, cfg *ProviderConfig, serviceName string) error {
	logger.Log(ctx, "debug", "Loading AWS configurations...")
	if err := p.loadConfig(); err != nil {
		return err
	}
	logger.Log(ctx, "debug", "AWS configs were loaded successfully!")

	logger.Log(ctx, "debug", "Creating AWS client...")
	client, err := p.createClient(ctx)
	if err != nil {
		return err
	}
	logger.Log(ctx, "debug", "AWS client was created successfully!")

	service, err := p.loadService(ctx, client, serviceName)
	if err != nil {
		return err
	}
	p.Service = service

	cfg.AWS = *p

	return nil
}

// Create a client for performing API calls
func (p *AWS) createClient(ctx context.Context) (*aws.Config, error) {
	credentials := credentials.NewStaticCredentialsProvider(p.config.Credentials.AccessKey, p.config.Credentials.SecretKey, "")

	config, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(p.config.Region),
		config.WithSharedConfigFiles([]string{p.config.Profile.Path}),
		config.WithSharedConfigProfile(p.config.Profile.Name),
		config.WithCredentialsProvider(credentials),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating AWS client: %v", err)
	}

	return &config, nil
}

// Read configs set by Viper
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

// Load the cloud provider's services
func (p *AWS) loadService(ctx context.Context, client *aws.Config, service string) (Cleanable, error) {
	ec2API := ec2.NewFromConfig(*client)
	elbAPI := elasticloadbalancingv2.NewFromConfig(*client)

	logger.Log(ctx, "debug", fmt.Sprintf("Initializing AWS service for: %s", service))

	switch service {
	case "targetGroup":
		return &targetgroup.TargetGroup{API: elbAPI}, nil
	case "loadBalancer":
		return &loadbalancer.LoadBalancer{API: elbAPI}, nil
	case "eni":
		return &elasticnetworkinterface.ElasticNetworkInterface{API: ec2API}, nil
	case "eip":
		return &elasticip.ElasticIP{API: ec2API}, nil
	case "ebs":
		return &elasticblockstorage.ElasticBlockStorage{API: ec2API}, nil
	default:
		return nil, fmt.Errorf("service %s is not supported", service)
	}
}
