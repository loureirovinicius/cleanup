package loadBalancer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type LoadBalancer struct {
	Client *aws.Config
}

func (r *LoadBalancer) List(ctx context.Context) ([]string, error) {
	var lbArns []string

	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	lbs, err := client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}

	for _, lb := range lbs.LoadBalancers {
		lbArns = append(lbArns, *lb.LoadBalancerArn)
	}

	return lbArns, nil
}

func (r *LoadBalancer) Validate(ctx context.Context, arn string) (bool, error) {
	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	listeners, err := client.DescribeListeners(ctx, &elasticloadbalancingv2.DescribeListenersInput{LoadBalancerArn: &arn})
	if err != nil {
		return false, err
	}

	return len(listeners.Listeners) == 0, nil
}

func (r *LoadBalancer) Delete(ctx context.Context, arn string) error {
	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	_, err := client.DeleteLoadBalancer(ctx, &elasticloadbalancingv2.DeleteLoadBalancerInput{LoadBalancerArn: &arn})
	if err != nil {
		return err
	}

	return nil
}
