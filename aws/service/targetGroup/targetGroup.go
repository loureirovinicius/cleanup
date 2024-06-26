package targetGroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/loureirovinicius/cleanup/helpers/logger"
)

type TargetGroup struct {
	Client *aws.Config
}

func (r *TargetGroup) List(ctx context.Context) []string {
	var tgArns []string

	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	tgs, err := client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{})
	if err != nil {
		msg := fmt.Sprintf("error retrieving Target Groups: %v", err)
		logger.Log(ctx, "error", msg)
		return nil
	}

	for _, tg := range tgs.TargetGroups {
		tgArns = append(tgArns, *tg.TargetGroupArn)
	}

	return tgArns
}

func (r *TargetGroup) Validate(ctx context.Context, name string) bool {
	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	tgs, err := client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{TargetGroupArns: []string{name}})
	if err != nil {
		msg := fmt.Sprintf("error retrieving Target Groups: %v", err)
		logger.Log(ctx, "error", msg)
		return false
	}

	lbs := len(tgs.TargetGroups[0].LoadBalancerArns)
	return lbs == 0
}

func (r *TargetGroup) Delete(ctx context.Context, arn string) string {
	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	_, err := client.DeleteTargetGroup(ctx, &elasticloadbalancingv2.DeleteTargetGroupInput{TargetGroupArn: &arn})
	if err != nil {
		msg := fmt.Sprintf("error retrieving Target Groups: %v", err)
		logger.Log(ctx, "error", msg)
		return ""
	}

	return "Target Group deleted successfully!"
}
