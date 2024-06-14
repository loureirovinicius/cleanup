package loadbalancer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type LoadBalancer struct {
	Client *aws.Config
}

func (r *LoadBalancer) List(ctx context.Context) []string {
	return []string{}
}
