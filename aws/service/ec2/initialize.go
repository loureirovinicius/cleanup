package ec2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elasticblockstorage "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticBlockStorage"
	elasticip "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticIp"
	elasticnetworkinterface "github.com/loureirovinicius/cleanup/aws/service/ec2/elasticNetworkInterface"
	loadbalancer "github.com/loureirovinicius/cleanup/aws/service/ec2/loadBalancer"
	targetgroup "github.com/loureirovinicius/cleanup/aws/service/ec2/targetGroup"
	"github.com/loureirovinicius/cleanup/cmd/cleaner"
)

func Initialize(config aws.Config) map[string]cleaner.Cleanable {
	ec2API := ec2.NewFromConfig(config)
	elbAPI := elasticloadbalancingv2.NewFromConfig(config)

	return map[string]cleaner.Cleanable{
		"targetGroup":  &targetgroup.TargetGroup{API: elbAPI},
		"loadBalancer": &loadbalancer.LoadBalancer{API: elbAPI},
		"eni":          &elasticnetworkinterface.ElasticNetworkInterface{API: ec2API},
		"eip":          &elasticip.ElasticIP{API: ec2API},
		"ebs":          &elasticblockstorage.ElasticBlockStorage{API: ec2API},
	}
}
