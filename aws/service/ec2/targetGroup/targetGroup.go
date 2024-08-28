package targetgroup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type TargetGroup struct {
	API TargetGroupAPI
}

type TargetGroupAPI interface {
	DescribeTargetGroups(ctx context.Context, params *elasticloadbalancingv2.DescribeTargetGroupsInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeTargetGroupsOutput, error)
	DeleteTargetGroup(ctx context.Context, params *elasticloadbalancingv2.DeleteTargetGroupInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteTargetGroupOutput, error)
}

func (r *TargetGroup) List(ctx context.Context) ([]string, error) {
	var tgArns []string

	tgs, err := r.API.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, err
	}

	for _, tg := range tgs.TargetGroups {
		tgArns = append(tgArns, *tg.TargetGroupArn)
	}

	return tgArns, nil
}

func (r *TargetGroup) Validate(ctx context.Context, arn string) (bool, error) {
	tgs, err := r.API.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{TargetGroupArns: []string{arn}})
	if err != nil {
		return false, err
	}

	lbs := len(tgs.TargetGroups[0].LoadBalancerArns)
	return lbs == 0, nil
}

func (r *TargetGroup) Delete(ctx context.Context, arn string) error {
	_, err := r.API.DeleteTargetGroup(ctx, &elasticloadbalancingv2.DeleteTargetGroupInput{TargetGroupArn: &arn})
	if err != nil {
		return err
	}

	return nil
}
