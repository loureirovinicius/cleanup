package elasticblockstorage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type ElasticBlockStorage struct {
	Client *aws.Config
}

func (r *ElasticBlockStorage) List(ctx context.Context) ([]string, error) {
	var ebsIds []string

	client := ec2.NewFromConfig(*r.Client)
	ebs, err := client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
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

	client := ec2.NewFromConfig(*r.Client)
	ebs, err := client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{VolumeIds: []string{id}})
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
	client := ec2.NewFromConfig(*r.Client)
	_, err := client.DeleteVolume(ctx, &ec2.DeleteVolumeInput{VolumeId: &id})
	if err != nil {
		return err
	}

	return nil
}
