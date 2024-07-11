package elasticip

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type ElasticIP struct {
	Client *aws.Config
}

func (r *ElasticIP) List(ctx context.Context) ([]string, error) {
	var eipsIds []string

	client := ec2.NewFromConfig(*r.Client)
	eips, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	for _, eip := range eips.Addresses {
		eipsIds = append(eipsIds, *eip.AllocationId)
	}

	return eipsIds, nil
}

func (r *ElasticIP) Validate(ctx context.Context, id string) (bool, error) {
	client := ec2.NewFromConfig(*r.Client)
	eips, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{AllocationIds: []string{id}})
	if err != nil {
		return false, err
	}

	status := eips.Addresses[0].AssociationId

	return status == nil, nil
}

func (r *ElasticIP) Delete(ctx context.Context, id string) error {
	client := ec2.NewFromConfig(*r.Client)
	_, err := client.ReleaseAddress(ctx, &ec2.ReleaseAddressInput{AllocationId: &id})
	if err != nil {
		return err
	}

	return nil
}
