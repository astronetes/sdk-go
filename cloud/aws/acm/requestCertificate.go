package acm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func RequestCertificate(ctx context.Context, client *acm.Client, req RequestCertificateRequest) RequestCertificateResponse {
	domainParts := strings.Split(req.domain, ".")
	if len(domainParts) < 2 {
		return RequestCertificateResponse{
			error: fmt.Errorf("invalid domain, It should have at least two blocks"),
		}
	}
	validationDomain := strings.Join(domainParts[len(domainParts)-2:], ".")
	request := &acm.RequestCertificateInput{
		DomainName: aws.String(req.domain),
		DomainValidationOptions: []acmTypes.DomainValidationOption{
			{
				DomainName:       aws.String(req.domain),
				ValidationDomain: aws.String(validationDomain),
			},
		},
		ValidationMethod: acmTypes.ValidationMethodDns,
	}
	response, err := client.RequestCertificate(ctx, request)
	return RequestCertificateResponse{
		response: response,
		error:    err,
	}
}

type RequestCertificateRequest struct {
	domain string
}

func NewRequestCertificateRequest(domain string) RequestCertificateRequest {
	return RequestCertificateRequest{domain: domain}
}

type RequestCertificateResponse struct {
	response *acm.RequestCertificateOutput
	error
}

func (r RequestCertificateResponse) Error() error {
	return r.error
}

func (r RequestCertificateResponse) CertificateARN() string {
	if r.response == nil {
		return ""
	}
	return *r.response.CertificateArn
}

func (r RequestCertificateResponse) RequestCertificateOutput() *acm.RequestCertificateOutput {
	if r.response == nil {
		return nil
	}
	return r.response
}
