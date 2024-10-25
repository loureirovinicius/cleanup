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

	logger.Log(ctx, "debug", "Starting the call to the DescribeAddresses API")
	eips, err := r.API.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	logger.Log(ctx, "debug", "Starting to loop through all the EIPs returned by API")
	for _, eip := range eips.Addresses {
		logger.Log(ctx, "debug", fmt.Sprintf("Appending Allocation ID (%v) to list of EIPs", *eip.AllocationId))
		eipsIds = append(eipsIds, *eip.AllocationId)
	}

	logger.Log(ctx, "debug", "Finished listing all the EIPs")
	return eipsIds, nil
}

func (r *ElasticIP) Validate(ctx context.Context, id string) (bool, error) {

	logger.Log(ctx, "debug", fmt.Sprintf("Starting the call to the DescribeAddresses API for EIP: %v", id))
	eips, err := r.API.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{AllocationIds: []string{id}})
	if err != nil {
		return false, err
	}

	status := eips.Addresses[0].AssociationId
	logger.Log(ctx, "debug", fmt.Sprintf("EIP address association ID: %v", status))

	logger.Log(ctx, "debug", fmt.Sprintf("Can the resource be deleted?: %v", status == nil))
	logger.Log(ctx, "debug", "Finished validating the EIP")
	return status == nil, nil
}

func (r *ElasticIP) Delete(ctx context.Context, id string) error {

	logger.Log(ctx, "debug", "Starting to call the ReleaseAddress API")
	_, err := r.API.ReleaseAddress(ctx, &ec2.ReleaseAddressInput{AllocationId: &id})
	if err != nil {
		return err
	}

	logger.Log(ctx, "debug", "Finished releasing the EIP")
	return nil
}
