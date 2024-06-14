package targetGroup

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/loureirovinicius/cleanup/cmd/cleaner"
)

type TargetGroup struct {
	Client *aws.Config
}

func (r *TargetGroup) List(ctx context.Context) []string {
	var tgNames []string

	client := elasticloadbalancingv2.NewFromConfig(*r.Client)
	tgs, err := client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{})
	if err != nil {
		log.Fatalf("error retrieving Target Groups: %v", err)
	}

	for _, tg := range tgs.TargetGroups {
		tgNames = append(tgNames, *tg.TargetGroupName)
	}

	return tgNames
}

func (TargetGroup) Validate(resource cleaner.Cleanable) bool {
	return true
}

func (TargetGroup) Delete(name string) string {
	return ""
}
