package targetGroup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type TargetGroup struct {
	Client *aws.Config
}

func (r *TargetGroup) List(ctx context.Context) ([]string, error) {
	var tgArns []string

	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	tgs, err := client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, err
	}

	for _, tg := range tgs.TargetGroups {
		tgArns = append(tgArns, *tg.TargetGroupArn)
	}

	return tgArns, nil
}

func (r *TargetGroup) Validate(ctx context.Context, arn string) (bool, error) {
	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	tgs, err := client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{TargetGroupArns: []string{arn}})
	if err != nil {
		return false, err
	}

	lbs := len(tgs.TargetGroups[0].LoadBalancerArns)
	return lbs == 0, nil
}

func (r *TargetGroup) Delete(ctx context.Context, arn string) error {
	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	_, err := client.DeleteTargetGroup(ctx, &elasticloadbalancingv2.DeleteTargetGroupInput{TargetGroupArn: &arn})
	if err != nil {
		return err
	}

	return nil
}
