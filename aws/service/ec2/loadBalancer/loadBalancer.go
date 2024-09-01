package loadbalancer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type LoadBalancer struct {
	API LoadBalancerAPI
}

type LoadBalancerAPI interface {
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error)
	DescribeListeners(ctx context.Context, params *elasticloadbalancingv2.DescribeListenersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeListenersOutput, error)
	DeleteLoadBalancer(ctx context.Context, params *elasticloadbalancingv2.DeleteLoadBalancerInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteLoadBalancerOutput, error)
}

func (r *LoadBalancer) List(ctx context.Context) ([]string, error) {
	var lbArns []string

	lbs, err := r.API.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}

	for _, lb := range lbs.LoadBalancers {
		lbArns = append(lbArns, *lb.LoadBalancerArn)
	}

	return lbArns, nil
}

func (r *LoadBalancer) Validate(ctx context.Context, arn string) (bool, error) {
	listeners, err := r.API.DescribeListeners(ctx, &elasticloadbalancingv2.DescribeListenersInput{LoadBalancerArn: &arn})
	if err != nil {
		return false, err
	}

	return len(listeners.Listeners) == 0, nil
}

func (r *LoadBalancer) Delete(ctx context.Context, arn string) error {
	_, err := r.API.DeleteLoadBalancer(ctx, &elasticloadbalancingv2.DeleteLoadBalancerInput{LoadBalancerArn: &arn})
	if err != nil {
		return err
	}

	return nil
}
