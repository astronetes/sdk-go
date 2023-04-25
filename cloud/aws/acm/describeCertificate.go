package acm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func DescribeCertificate(ctx context.Context, client *acm.Client, req DescribeCertificateRequest) DescribeCertificateResponse {
	if req.certificateARN == "" {
		return DescribeCertificateResponse{
			error: fmt.Errorf("missing required certificate ARN"),
		}
	}
	input := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(req.certificateARN),
	}
	response, err := client.DescribeCertificate(ctx, input)
	return DescribeCertificateResponse{
		response: response,
		error:    err,
	}
}

type DescribeCertificateRequest struct {
	certificateARN string
}

func NewDescribeCertificateRequest(certificateARN string) DescribeCertificateRequest {
	return DescribeCertificateRequest{certificateARN: certificateARN}
}

type DescribeCertificateResponse struct {
	response *acm.DescribeCertificateOutput
	error
}

func (r DescribeCertificateResponse) Error() error {
	return r.error
}

func (r DescribeCertificateResponse) Response() *acm.DescribeCertificateOutput {
	return r.response
}

func (r *DescribeCertificateResponse) DomainValidations() []types.DomainValidation {
	return r.response.Certificate.DomainValidationOptions
}

func (r *DescribeCertificateResponse) IsStatus(expected types.CertificateStatus) bool {
	return r.response.Certificate.Status == expected
}

func (r *DescribeCertificateResponse) ResourceRecordForDefaultValidationOption() (*types.ResourceRecord, error) {
	if r.response == nil {
		return nil, fmt.Errorf("missing response")
	}
	if r.response.Certificate == nil {
		return nil, fmt.Errorf("missing certificate")
	}
	if len(r.response.Certificate.DomainValidationOptions) == 0 {
		return nil, fmt.Errorf("missing domain validation options")
	}
	record := r.response.Certificate.DomainValidationOptions[0].ResourceRecord
	if record == nil {
		return nil, fmt.Errorf("missing resource record for domain validation options")
	}
	return record, nil
}
