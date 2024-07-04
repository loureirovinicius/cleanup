package elasticNetworkInterface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type ElasticNetworkInterface struct {
	Client *aws.Config
}

func (r *ElasticNetworkInterface) List(ctx context.Context) ([]string, error) {
	var eniIds []string

	client := ec2.NewFromConfig(*r.Client)
	enis, err := client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return nil, err
	}

	for _, eni := range enis.NetworkInterfaces {
		eniIds = append(eniIds, *eni.NetworkInterfaceId)
	}

	return eniIds, nil
}

func (r *ElasticNetworkInterface) Validate(ctx context.Context, id string) (bool, error) {
	client := ec2.NewFromConfig(*r.Client)
	enis, err := client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{NetworkInterfaceIds: []string{id}})
	if err != nil {
		return false, err
	}

	status := enis.NetworkInterfaces[0].Status

	return status == "available", nil
}

func (r *ElasticNetworkInterface) Delete(ctx context.Context, id string) error {
	client := ec2.NewFromConfig(*r.Client)
	_, err := client.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{NetworkInterfaceId: &id})
	if err != nil {
		return err
	}

	return nil
}
