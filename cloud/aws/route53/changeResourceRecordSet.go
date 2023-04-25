package route53

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func UpsertResourceRecordSet(ctx context.Context, client *route53.Client, req ChangeResourceRecordSetRequest) ChangeResourceRecordSetResponse {
	return changeResourceRecordSetResponse(ctx, client, req)
}

func DeleteResourceRecordSet(ctx context.Context, client *route53.Client, req ChangeResourceRecordSetRequest) ChangeResourceRecordSetResponse {
	return changeResourceRecordSetResponse(ctx, client, req)
}

func changeResourceRecordSetResponse(ctx context.Context, client *route53.Client, req ChangeResourceRecordSetRequest) ChangeResourceRecordSetResponse {
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action: req.action,
					ResourceRecordSet: &route53types.ResourceRecordSet{
						Name: aws.String(req.name),
						Type: route53types.RRTypeCname,
						ResourceRecords: []route53types.ResourceRecord{
							{
								Value: aws.String(req.value),
							},
						},
						TTL: aws.Int64(req.ttl),
					},
				},
			},
		},
		HostedZoneId: aws.String(req.zoneID),
	}
	response, err := client.ChangeResourceRecordSets(ctx, input)
	return ChangeResourceRecordSetResponse{
		response: response,
		err:      err,
	}
}

type ChangeResourceRecordSetRequest struct {
	action route53types.ChangeAction
	name   string
	value  string
	ttl    int64
	zoneID string
}

func NewChangeResourceRecordSetRequest(action route53types.ChangeAction, name string, value string, ttl int64, zoneID string) ChangeResourceRecordSetRequest {
	return ChangeResourceRecordSetRequest{action: action, name: name, value: value, ttl: ttl, zoneID: zoneID}
}

type ChangeResourceRecordSetResponse struct {
	err      error
	response *route53.ChangeResourceRecordSetsOutput
}

func (r ChangeResourceRecordSetResponse) Error() error {
	return r.err
}

func (r ChangeResourceRecordSetResponse) Response() *route53.ChangeResourceRecordSetsOutput {
	return r.response
}

func (r ChangeResourceRecordSetResponse) ChangeInfoID() (string, error) {
	if r.response == nil {
		return "", fmt.Errorf("missing response")
	}
	if r.response.ChangeInfo == nil {
		return "", fmt.Errorf("missing change info from the response")
	}
	return *r.response.ChangeInfo.Id, nil
}
