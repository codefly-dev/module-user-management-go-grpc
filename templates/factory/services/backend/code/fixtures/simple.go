package fixtures

import (
	"backend/pkg/business"
	"backend/pkg/gen"
	"context"

	"github.com/codefly-dev/core/wool"
)

func Simple(ctx context.Context, service *business.Service) error {
	w := wool.Get(ctx).In("fixtures.Simple")
	w.Info("Applying simple fixtures")
	authID := "primary"
	email := "primary@test.com"

	u := &gen.User{
		SignupAuthId: authID,
		Email:        email,
	}

	_, err := service.DeleteOwner(ctx, authID)
	if err != nil {
		return w.Wrapf(err, "cannot delete user")
	}

	_, err = service.RegisterUser(ctx, u)
	if err != nil {
		return w.Wrapf(err, "cannot register user")
	}
	return nil
}
