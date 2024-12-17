package elasticnetworkinterface

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/loureirovinicius/cleanup/helpers/logger"
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

	logger.Log(ctx, "debug", "Starting to list all the ENIs")
	enis, err := r.API.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return nil, fmt.Errorf("error calling the AWS DescribeNetworkInterfaces API: %v", err)
	}

	for _, eni := range enis.NetworkInterfaces {
		eniIds = append(eniIds, *eni.NetworkInterfaceId)
	}

	logger.Log(ctx, "debug", "Finished listing all the ENIs")
	return eniIds, nil
}

func (r *ElasticNetworkInterface) Validate(ctx context.Context, id string) (bool, error) {
	logger.Log(ctx, "debug", fmt.Sprintf("Validating ENI: %v", id))
	enis, err := r.API.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{NetworkInterfaceIds: []string{id}})
	if err != nil {
		return false, fmt.Errorf("error calling the AWS DescribeNetworkInterfaces API: %v", err)
	}

	if len(enis.NetworkInterfaces) == 0 {
		logger.Log(ctx, "info", "No ENI found for ID: %v", id)
		return false, nil
	}

	status := enis.NetworkInterfaces[0].Status
	logger.Log(ctx, "debug", fmt.Sprintf("ENI status: %v", status))

	logger.Log(ctx, "debug", "Finished validating the ENI")
	return status == "available", nil
}

func (r *ElasticNetworkInterface) Delete(ctx context.Context, id string) error {
	logger.Log(ctx, "debug", "Deleting ENI: %v", id)
	_, err := r.API.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{NetworkInterfaceId: &id})
	if err != nil {
		return fmt.Errorf("error calling the AWS DeleteNetworkInterface API: %v", err)
	}

	logger.Log(ctx, "debug", "Finished deleting the ENI")
	return nil
}
