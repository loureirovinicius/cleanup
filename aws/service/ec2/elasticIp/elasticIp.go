package elasticip

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/loureirovinicius/cleanup/helpers/logger"
)

type ElasticIP struct {
	API ElasticIPAPI
}

type ElasticIPAPI interface {
	DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAddressesOutput, error)
	ReleaseAddress(ctx context.Context, params *ec2.ReleaseAddressInput, optFns ...func(*ec2.Options)) (*ec2.ReleaseAddressOutput, error)
}

func (r *ElasticIP) List(ctx context.Context) ([]string, error) {
	var eipsIds []string

	logger.Log(ctx, "debug", "Starting to list all the EIPs")
	eips, err := r.API.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, fmt.Errorf("error calling the AWS DescribeAddresses API: %v", err)
	}

	for _, eip := range eips.Addresses {
		eipsIds = append(eipsIds, *eip.AllocationId)
	}

	logger.Log(ctx, "debug", "Finished listing all the EIPs")
	return eipsIds, nil
}

func (r *ElasticIP) Validate(ctx context.Context, id string) (bool, error) {

	logger.Log(ctx, "debug", fmt.Sprintf("Starting the call to the DescribeAddresses API for EIP: %v", id))
	eips, err := r.API.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{AllocationIds: []string{id}})
	if err != nil {
		return false, fmt.Errorf("error calling the AWS DescribeAddresses API: %v", err)
	}

	if len(eips.Addresses) == 0 {
		logger.Log(ctx, "info", "No EIP found for ID: %v", id)
		return false, nil
	}

	status := eips.Addresses[0].AssociationId
	logger.Log(ctx, "debug", fmt.Sprintf("EIP address association ID: %v", status))

	logger.Log(ctx, "debug", "Finished validating the EIP")
	return status == nil, nil
}

func (r *ElasticIP) Delete(ctx context.Context, id string) error {
	logger.Log(ctx, "debug", "Releasing EIP: %v", id)
	_, err := r.API.ReleaseAddress(ctx, &ec2.ReleaseAddressInput{AllocationId: &id})
	if err != nil {
		return fmt.Errorf("error calling the AWS ReleaseAddress API: %v", err)
	}

	logger.Log(ctx, "debug", "Finished releasing the EIP")
	return nil
}
