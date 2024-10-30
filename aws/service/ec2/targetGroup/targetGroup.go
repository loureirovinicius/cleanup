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

	logger.Log(ctx, "debug", "Starting the call to the DescribeTargetGroups API")
	tgs, err := r.API.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, err
	}

	logger.Log(ctx, "debug", "Starting to loop through all the TGs returned by API")
	for _, tg := range tgs.TargetGroups {
		logger.Log(ctx, "debug", fmt.Sprintf("Appending TG ARN (%v) to list of TG ARNs", *tg.TargetGroupArn))
		tgArns = append(tgArns, *tg.TargetGroupArn)
	}

	logger.Log(ctx, "debug", "Finished listing all the TGs")
	return tgArns, nil
}

func (r *TargetGroup) Validate(ctx context.Context, arn string) (bool, error) {

	logger.Log(ctx, "debug", fmt.Sprintf("Starting the call to the DescribeTargetGroups API for TG: %v", arn))
	tgs, err := r.API.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{TargetGroupArns: []string{arn}})
	if err != nil {
		return false, err
	}

	lbs := len(tgs.TargetGroups[0].LoadBalancerArns)
	logger.Log(ctx, "debug", fmt.Sprintf("Listeners for LB (%v): %v", arn, lbs))

	logger.Log(ctx, "debug", fmt.Sprintf("Can the resource be deleted?: %v", lbs == 0))
	logger.Log(ctx, "debug", "Finished validating the TG")
	return lbs == 0, nil
}

func (r *TargetGroup) Delete(ctx context.Context, arn string) error {

	logger.Log(ctx, "debug", "Starting to call the DeleteTargetGroup API")
	_, err := r.API.DeleteTargetGroup(ctx, &elasticloadbalancingv2.DeleteTargetGroupInput{TargetGroupArn: &arn})
	if err != nil {
		return err
	}

	logger.Log(ctx, "debug", "Finished deleting the TG")
	return nil
}
