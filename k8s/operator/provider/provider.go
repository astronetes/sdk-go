package provider

import (
	"context"
	"fmt"
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
	SetUp(ctx context.Context, cfg C) error
	Status(ctx context.Context, obj T) Status
	Create(ctx context.Context, obj T) error
	Delete(ctx context.Context, obj T) error
	CanBeUpdated(ctx context.Context, obj T) (bool, error)
	CanBeDeleted(ctx context.Context, obj T) (bool, error)
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

func (m Manager[T, C]) Get(ctx context.Context, cfg C, providerID ID) (Provider[T, C], error) {
	provider, ok := m.providers[providerID]
	if !ok {
		return nil, fmt.Errorf("unsupported provider id '%v' for handling this resource", providerID)
	}
	if err := provider.SetUp(ctx, cfg); err != nil {
		return nil, err
	}
	return provider, nil
}
