package provider

import (
	"context"
	"fmt"
)

type ID string
type Status int32

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

func (m *Manager[T, C]) Register(providerID ID, provider Provider[T, C]) {
	m.providers[providerID] = provider
}

func (m *Manager[T, C]) Get(providerID ID, ctx context.Context, cfg C) (Provider[T, C], error) {
	provider, ok := m.providers[providerID]
	if !ok {
		return nil, fmt.Errorf("unsupported provider id '%v' for handling this resource", providerID)
	}
	if err := provider.SetUp(ctx, cfg); err != nil {
		return nil, err
	}
	return provider, nil
}
