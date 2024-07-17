package main

import (
	"backend/fixtures"
	"backend/pkg/business"
	"context"

	codefly "github.com/codefly-dev/sdk-go"

	"backend/pkg/infra"

	"backend/pkg/adapters"
)

func init() {
	WithWork(doWork)
}

func doWork(ctx context.Context) (Clean, error) {
	store, err := infra.NewPostgresStore(ctx)
	if err != nil {
		return nil, err
	}
	service := &business.Service{}
	service.SetStore(store)

	adapters.WithService(service)

	if codefly.WithFixture("simple") {
		err = fixtures.Simple(ctx, service)
		if err != nil {
			return nil, err
		}
	}

	return func() {
		store.Close()
	}, nil
}
