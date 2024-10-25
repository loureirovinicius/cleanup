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

	logger.Log(ctx, "debug", "Starting the call to the DescribeLoadBalancers API")
	lbs, err := r.API.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}

	logger.Log(ctx, "debug", "Starting to loop through all the LBs returned by API")
	for _, lb := range lbs.LoadBalancers {
		logger.Log(ctx, "debug", fmt.Sprintf("Appending LB ARN (%v) to list of LB ARNs", *lb.LoadBalancerArn))
		lbArns = append(lbArns, *lb.LoadBalancerArn)
	}

	logger.Log(ctx, "debug", "Finished listing all the LBs")
	return lbArns, nil
}

func (r *LoadBalancer) Validate(ctx context.Context, arn string) (bool, error) {

	logger.Log(ctx, "debug", fmt.Sprintf("Starting the call to the DescribeListeners API for LB: %v", arn))
	listeners, err := r.API.DescribeListeners(ctx, &elasticloadbalancingv2.DescribeListenersInput{LoadBalancerArn: &arn})
	if err != nil {
		return false, err
	}

	logger.Log(ctx, "debug", fmt.Sprintf("LB listeners count: %v", len(listeners.Listeners)))

	logger.Log(ctx, "debug", fmt.Sprintf("Can the resource be deleted?: %v", len(listeners.Listeners) == 0))
	logger.Log(ctx, "debug", "Finished validating the LB")
	return len(listeners.Listeners) == 0, nil
}

func (r *LoadBalancer) Delete(ctx context.Context, arn string) error {

	logger.Log(ctx, "debug", "Starting to call the DeleteLoadBalancer API")
	_, err := r.API.DeleteLoadBalancer(ctx, &elasticloadbalancingv2.DeleteLoadBalancerInput{LoadBalancerArn: &arn})
	if err != nil {
		return err
	}

	logger.Log(ctx, "debug", "Finished deleting the LB")
	return nil
}
