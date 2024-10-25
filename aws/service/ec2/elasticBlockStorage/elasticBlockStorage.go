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

	logger.Log(ctx, "debug", "Starting the call to the DescribeVolumes API")
	ebs, err := r.API.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, err
	}

	logger.Log(ctx, "debug", "Starting to loop through all the EBS returned by API")
	for _, ebs := range ebs.Volumes {
		logger.Log(ctx, "debug", fmt.Sprintf("Appending EBS ID (%v) to list of EBS IDs", *ebs.VolumeId))
		ebsIds = append(ebsIds, *ebs.VolumeId)
	}

	logger.Log(ctx, "debug", "Finished listing all the EBS volumes")
	return ebsIds, nil
}

func (r *ElasticBlockStorage) Validate(ctx context.Context, id string) (bool, error) {
	var tagged bool

	logger.Log(ctx, "debug", fmt.Sprintf("Starting the call to the DescribeVolumes API for EBS: %v", id))
	ebs, err := r.API.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{VolumeIds: []string{id}})
	if err != nil {
		return false, err
	}

	state := ebs.Volumes[0].State
	logger.Log(ctx, "debug", fmt.Sprintf("EBS state: %v", state))
	tags := ebs.Volumes[0].Tags
	logger.Log(ctx, "debug", fmt.Sprintf("EBS tags: %v", tags))

	logger.Log(ctx, "debug", "Starting the tag validation to check if volume can be deleted")
	for _, v := range tags {
		tagged = *v.Key == "cleanup-ignore" && *v.Value == "true"
	}

	logger.Log(ctx, "debug", fmt.Sprintf("Can the resource be deleted?: %v", state == "available" && !tagged))
	logger.Log(ctx, "debug", "Finished validating the EBS volume")
	return state == "available" && !tagged, nil
}

func (r *ElasticBlockStorage) Delete(ctx context.Context, id string) error {

	logger.Log(ctx, "debug", "Starting the call to the DeleteVolume API")
	_, err := r.API.DeleteVolume(ctx, &ec2.DeleteVolumeInput{VolumeId: &id})
	if err != nil {
		return err
	}

	logger.Log(ctx, "debug", "Fnished deleting the EBS volume")
	return nil
}
