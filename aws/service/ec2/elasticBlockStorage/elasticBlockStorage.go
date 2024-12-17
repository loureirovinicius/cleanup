package elasticblockstorage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/loureirovinicius/cleanup/helpers/logger"
)

type ElasticBlockStorage struct {
	API ElasticBlockStorageAPI
}

type ElasticBlockStorageAPI interface {
	DescribeVolumes(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error)
	DeleteVolume(ctx context.Context, params *ec2.DeleteVolumeInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error)
}

func (r *ElasticBlockStorage) List(ctx context.Context) ([]string, error) {
	var ebsIds []string

	logger.Log(ctx, "debug", "Starting to list all the EBS volumes")
	ebs, err := r.API.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, fmt.Errorf("error calling the AWS DescribeVolumes API: %w", err)
	}

	for _, ebs := range ebs.Volumes {
		ebsIds = append(ebsIds, *ebs.VolumeId)
	}

	logger.Log(ctx, "debug", "Finished listing all the EBS volumes")
	return ebsIds, nil
}

func (r *ElasticBlockStorage) Validate(ctx context.Context, id string) (bool, error) {
	var tagged bool

	logger.Log(ctx, "debug", fmt.Sprintf("Validating EBS volume: %v", id))
	ebs, err := r.API.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{VolumeIds: []string{id}})
	if err != nil {
		return false, fmt.Errorf("error calling the AWS DescribeVolumes API: %w", err)
	}

	if len(ebs.Volumes) == 0 {
		logger.Log(ctx, "info", "No volume found for ID: %v", id)
		return false, nil
	}

	state := ebs.Volumes[0].State
	logger.Log(ctx, "debug", fmt.Sprintf("EBS state: %v", state))
	tags := ebs.Volumes[0].Tags
	logger.Log(ctx, "debug", fmt.Sprintf("EBS tags: %v", tags))

	for _, v := range tags {
		if *v.Key == "cleanup-ignore" && *v.Value == "true" {
			tagged = true
			break
		}
	}

	logger.Log(ctx, "debug", "Finished validating the EBS volume")
	return state == "available" && !tagged, nil
}

func (r *ElasticBlockStorage) Delete(ctx context.Context, id string) error {
	logger.Log(ctx, "debug", "Deleting EBS volume: %v", id)
	_, err := r.API.DeleteVolume(ctx, &ec2.DeleteVolumeInput{VolumeId: &id})
	if err != nil {
		return fmt.Errorf("error calling the AWS DeleteVolume API: %w", err)
	}

	logger.Log(ctx, "debug", "Finished deleting the EBS volume")
	return nil
}
