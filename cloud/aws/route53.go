package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/route53"
)

type Route53Client interface {
	ChangeResourceRecordsSets(ctx context.Context) error
}

type route53Client struct {
	route53 route53.Route53
}
