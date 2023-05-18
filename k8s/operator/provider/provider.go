package provider

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"github.com/astronetes/sdk-go/k8s/operator/errors"

	ctrl "sigs.k8s.io/controller-runtime"
)

type (
	ID      string
	GroupID string
)

type Provider[T v1.Resource] interface {
	SetUp(ctx context.Context, mgr manager.Manager, cfg Config) error
	Create(ctx context.Context, obj T) (*ctrl.Result, error)
	Delete(ctx context.Context, obj T) (*ctrl.Result, error)
}

type Manager[T v1.Resource] struct {
	providers map[ID]Provider[T]
}

func (m Manager[T]) WithProvider(providerID ID, provider Provider[T]) Manager[T] {
	if m.providers == nil {
		m.providers = make(map[ID]Provider[T])
	}
	m.providers[providerID] = provider
	return m
}

func (m Manager[T]) Get(ctx context.Context, mgr manager.Manager, cfg Config, providerID ID) (Provider[T], error) {
	provider, ok := m.providers[providerID]
	if !ok {
		return nil, errors.ProviderError("unsupported provider id '%v' for handling this resource", providerID)
	}
	if err := provider.SetUp(ctx, mgr, cfg); err != nil {
		return nil, err
	}
	return provider, nil
}
