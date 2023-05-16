package provider

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	ID     string
	Status int32
)

const (
	Uncreated = iota
	OnCreation
	Ready
	OnDeletion
	Deleted
	Unknown
)

type Provider[T any, C any] interface {
	SetUp(ctx context.Context, runtimeClient client.Client, cfg C) error
	Status(ctx context.Context, obj T) (Status, error)
	Create(ctx context.Context, obj T) error
	Delete(ctx context.Context, obj T) error
}

type Manager[T any, C any] struct {
	providers map[ID]Provider[T, C]
}

func (m Manager[T, C]) WithProvider(providerID ID, provider Provider[T, C]) Manager[T, C] {
	if m.providers == nil {
		m.providers = make(map[ID]Provider[T, C])
	}
	m.providers[providerID] = provider
	return m
}

func (m Manager[T, C]) Get(ctx context.Context, runtimeClient client.Client, cfg C, providerID ID) (Provider[T, C], error) {
	provider, ok := m.providers[providerID]
	if !ok {
		return nil, fmt.Errorf("unsupported provider id '%v' for handling this resource", providerID)
	}
	if err := provider.SetUp(ctx, runtimeClient, cfg); err != nil {
		return nil, err
	}
	return provider, nil
}
