package provider

import (
	"context"

	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"github.com/astronetes/sdk-go/k8s/operator/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	ID      string
	GroupID string
	Status  int32
)

const (
	Uncreated = iota
	OnCreation
	Ready
	OnDeletion
	Deleted
	Unknown
)

type Provider[T v1.Resource] interface {
	SetUp(ctx context.Context, runtimeClient client.Client, cfg Config) error
	Status(ctx context.Context, obj T) (Status, error)
	Create(ctx context.Context, obj T) error
	Delete(ctx context.Context, obj T) error
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

func (m Manager[T]) Get(ctx context.Context, runtimeClient client.Client, cfg Config, providerID ID) (Provider[T], error) {
	provider, ok := m.providers[providerID]
	if !ok {
		return nil, errors.ProviderError("unsupported provider id '%v' for handling this resource", providerID)
	}
	if err := provider.SetUp(ctx, runtimeClient, cfg); err != nil {
		return nil, err
	}
	return provider, nil
}
