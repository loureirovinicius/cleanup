package elasticblockstorage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

	ebs, err := r.API.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, err
	}

	for _, ebs := range ebs.Volumes {
		ebsIds = append(ebsIds, *ebs.VolumeId)
	}

	return ebsIds, nil
}

func (r *ElasticBlockStorage) Validate(ctx context.Context, id string) (bool, error) {
	var tagged bool

	ebs, err := r.API.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{VolumeIds: []string{id}})
	if err != nil {
		return false, err
	}

	state := ebs.Volumes[0].State
	tags := ebs.Volumes[0].Tags

	for _, v := range tags {
		tagged = *v.Key == "cleanup-ignore" && *v.Value == "true"
	}

	return state == "available" && !tagged, nil
}

func (r *ElasticBlockStorage) Delete(ctx context.Context, id string) error {
	_, err := r.API.DeleteVolume(ctx, &ec2.DeleteVolumeInput{VolumeId: &id})
	if err != nil {
		return err
	}

	return nil
}
