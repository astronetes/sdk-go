package aws

import (
	"context"

	"github.com/astronetes/sdk-go/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53Client interface {
	CreateResource(ctx context.Context, name string, value string, ttl int64, zoneID string) error
	DeleteResource(ctx context.Context, name string, value string, ttl int64, zoneID string) error
}

type route53Client struct {
	client route53.Client
}

func (c *route53Client) changeResourceRecourSet(ctx context.Context, action route53types.ChangeAction, name string, value string, ttl int64, zoneID string) error {
	logger := log.FromContext(ctx)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action: action,
					ResourceRecordSet: &route53types.ResourceRecordSet{
						Name: aws.String(name),
						Type: route53types.RRTypeCname,
						ResourceRecords: []route53types.ResourceRecord{
							{
								Value: aws.String(value),
							},
						},
						TTL: aws.Int64(ttl),
					},
				},
			},
		},
		HostedZoneId: aws.String(zoneID),
	}
	response, err := c.client.ChangeResourceRecordSets(ctx, input)
	if err != nil {
		logger.V(log.Error).Info("Unable to change resource record set: '%w'", err)
		return err
	}
	logger.V(log.Info).Info("action over resource record was performed successfully, current status: '%s'", response.ChangeInfo.Status)
	return nil
}

func (c *route53Client) CreateResource(ctx context.Context, name string, value string, ttl int64, zoneID string) error {
	return c.changeResourceRecourSet(ctx, route53types.ChangeActionCreate, name, value, ttl, zoneID)
}

func (c *route53Client) DeleteResource(ctx context.Context, name string, value string, ttl int64, zoneID string) error {
	return c.changeResourceRecourSet(ctx, route53types.ChangeActionDelete, name, value, ttl, zoneID)
}
