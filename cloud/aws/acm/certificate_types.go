package acm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

type Certificate struct {
	Raw *types.CertificateDetail
}

func (r *Certificate) DomainValidations() []types.DomainValidation {
	return r.Raw.DomainValidationOptions
}

func (r *Certificate) IsStatus(expected types.CertificateStatus) bool {
	return r.Raw.Status == expected
}

func (r *Certificate) ResourceRecordForDefaultValidationOption() (*types.ResourceRecord, error) {
	if r.Raw == nil {
		return nil, fmt.Errorf("missing certificate")
	}
	if len(r.Raw.DomainValidationOptions) == 0 {
		return nil, fmt.Errorf("missing domain validation options")
	}
	record := r.Raw.DomainValidationOptions[0].ResourceRecord
	if record == nil {
		return nil, fmt.Errorf("missing resource record for domain validation options")
	}
	return record, nil
}

type RequestCertificateRequest struct {
	Domain string
}

type RequestCertificateResponse struct {
	Arn *string
}
