package loadbalancer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/loureirovinicius/cleanup/helpers/logger"
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

	logger.Log(ctx, "debug", "Starting to list all the LBs")
	lbs, err := r.API.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, fmt.Errorf("error calling the AWS DescribeLoadBalancers API: %w", err)
	}

	for _, lb := range lbs.LoadBalancers {
		lbArns = append(lbArns, *lb.LoadBalancerArn)
	}

	logger.Log(ctx, "debug", "Finished listing all the LBs")
	return lbArns, nil
}

func (r *LoadBalancer) Validate(ctx context.Context, arn string) (bool, error) {
	logger.Log(ctx, "debug", fmt.Sprintf("Validating LB: %v", arn))
	listeners, err := r.API.DescribeListeners(ctx, &elasticloadbalancingv2.DescribeListenersInput{LoadBalancerArn: &arn})
	if err != nil {
		return false, fmt.Errorf("error calling the AWS DescribeLoadBalancers API: %w", err)
	}

	logger.Log(ctx, "debug", fmt.Sprintf("LB listeners count: %v", len(listeners.Listeners)))

	logger.Log(ctx, "debug", "Finished validating the LB")
	return len(listeners.Listeners) == 0, nil
}

func (r *LoadBalancer) Delete(ctx context.Context, arn string) error {
	logger.Log(ctx, "debug", "Deleting LB: %v", arn)
	_, err := r.API.DeleteLoadBalancer(ctx, &elasticloadbalancingv2.DeleteLoadBalancerInput{LoadBalancerArn: &arn})
	if err != nil {
		return fmt.Errorf("error calling the AWS DeleteLoadBalancer API: %w", err)
	}

	logger.Log(ctx, "debug", "Finished deleting the LB")
	return nil
}
