package targetgroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/loureirovinicius/cleanup/helpers/logger"
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

	logger.Log(ctx, "debug", "Starting to list all the TargetGroups")
	tgs, err := r.API.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, fmt.Errorf("error calling the AWS DescribeTargetGroups API: %v", err)
	}

	for _, tg := range tgs.TargetGroups {
		tgArns = append(tgArns, *tg.TargetGroupArn)
	}

	logger.Log(ctx, "debug", "Finished listing all the TGs")
	return tgArns, nil
}

func (r *TargetGroup) Validate(ctx context.Context, arn string) (bool, error) {
	logger.Log(ctx, "debug", fmt.Sprintf("Validating TG: %v", arn))
	tgs, err := r.API.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{TargetGroupArns: []string{arn}})
	if err != nil {
		return false, fmt.Errorf("error calling the AWS DescribeTargetGroups API: %v", err)
	}

	lbs := len(tgs.TargetGroups[0].LoadBalancerArns)
	logger.Log(ctx, "debug", fmt.Sprintf("LBs for TargetGroup (%v): %v", arn, lbs))

	logger.Log(ctx, "debug", "Finished validating the TG")
	return lbs == 0, nil
}

func (r *TargetGroup) Delete(ctx context.Context, arn string) error {
	logger.Log(ctx, "debug", "Deleting TG: %v", arn)
	_, err := r.API.DeleteTargetGroup(ctx, &elasticloadbalancingv2.DeleteTargetGroupInput{TargetGroupArn: &arn})
	if err != nil {
		return fmt.Errorf("error calling the AWS DeleteTargetGroup API: %v", err)
	}

	logger.Log(ctx, "debug", "Finished deleting the TG")
	return nil
}
