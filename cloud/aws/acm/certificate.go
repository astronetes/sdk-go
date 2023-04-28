package acm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// TODO See issue https://github.com/astronetes/sdk-go/issues/6
func RequestCertificate(ctx context.Context, client *acm.Client, req RequestCertificateRequest) (*RequestCertificateResponse, error) {
	domainParts := strings.Split(req.Domain, ".")
	if len(domainParts) < 2 {
		return nil, fmt.Errorf("invalid domain, It should have at least two blocks")
	}

	validationDomain := strings.Join(domainParts[len(domainParts)-2:], ".")
	request := &acm.RequestCertificateInput{
		DomainName: aws.String(req.Domain),
		DomainValidationOptions: []types.DomainValidationOption{
			{
				DomainName:       aws.String(req.Domain),
				ValidationDomain: aws.String(validationDomain),
			},
		},
		ValidationMethod: types.ValidationMethodDns,
	}
	response, err := client.RequestCertificate(ctx, request)
	if err != nil {
		return nil, err
	}

	return &RequestCertificateResponse{
		Arn: response.CertificateArn,
	}, nil
}

func GetCertificate(ctx context.Context, client *acm.Client, arn string) (*types.CertificateDetail, error) {
	if arn == "" {
		return nil, fmt.Errorf("missing required certificate ARN")
	}
	input := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	}

	response, err := client.DescribeCertificate(ctx, input)
	if err != nil {
		return nil, err
	}
	return response.Certificate, nil
}

func DeleteCertificate(ctx context.Context, client *acm.Client, arn string) error {
	if arn == "" {
		return fmt.Errorf("missing required certificate ARN")
	}

	_, err := GetCertificate(ctx, client, arn)
	if err != nil {
		if IsNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}

	input := &acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	}

	_, err = client.DeleteCertificate(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
