package elasticnetworkinterface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type ElasticNetworkInterface struct {
	API ElasticNetworkInterfaceAPI
}

type ElasticNetworkInterfaceAPI interface {
	DescribeNetworkInterfaces(ctx context.Context, params *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error)
	DeleteNetworkInterface(ctx context.Context, params *ec2.DeleteNetworkInterfaceInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInterfaceOutput, error)
}

func (r *ElasticNetworkInterface) List(ctx context.Context) ([]string, error) {
	var eniIds []string

	enis, err := r.API.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return nil, err
	}

	for _, eni := range enis.NetworkInterfaces {
		eniIds = append(eniIds, *eni.NetworkInterfaceId)
	}

	return eniIds, nil
}

func (r *ElasticNetworkInterface) Validate(ctx context.Context, id string) (bool, error) {
	enis, err := r.API.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{NetworkInterfaceIds: []string{id}})
	if err != nil {
		return false, err
	}

	status := enis.NetworkInterfaces[0].Status

	return status == "available", nil
}

func (r *ElasticNetworkInterface) Delete(ctx context.Context, id string) error {
	_, err := r.API.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{NetworkInterfaceId: &id})
	if err != nil {
		return err
	}

	return nil
}
