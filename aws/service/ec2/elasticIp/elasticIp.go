package elasticip

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

	eips, err := r.API.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	for _, eip := range eips.Addresses {
		eipsIds = append(eipsIds, *eip.AllocationId)
	}

	return eipsIds, nil
}

func (r *ElasticIP) Validate(ctx context.Context, id string) (bool, error) {
	eips, err := r.API.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{AllocationIds: []string{id}})
	if err != nil {
		return false, err
	}

	status := eips.Addresses[0].AssociationId

	return status == nil, nil
}

func (r *ElasticIP) Delete(ctx context.Context, id string) error {
	_, err := r.API.ReleaseAddress(ctx, &ec2.ReleaseAddressInput{AllocationId: &id})
	if err != nil {
		return err
	}

	return nil
}
