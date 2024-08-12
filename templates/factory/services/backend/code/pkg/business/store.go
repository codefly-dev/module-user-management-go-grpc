package business

import (
	"context"

	"backend/pkg/gen"
)

type Store interface {

	// User manipulation

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *gen.User) (*gen.User, error)

	// GetUserByAuthID returns a user by the auth id
	// - user, nil if found
	// - nil, nil if not found
	// - nil, err if something wrong
	GetUserByAuthID(ctx context.Context, id string) (*gen.User, error)

	// DeleteUser deletes a user by the auth id
	// - user, nil if found
	// - nil, nil if not found
	// - nil, err if something wrong
	DeleteUser(ctx context.Context, authSignupId string) (*gen.User, error)

	// Organization

	GetOrganizationForOwner(ctx context.Context, user *gen.User) (*gen.Organization, error)

	CreateOrganization(ctx context.Context, owner *gen.User, org *gen.Organization) (*gen.Organization, error)

	DeleteOrganization(ctx context.Context, org *gen.Organization) error

	// Teams

	CreateTeam(ctx context.Context, org *gen.Organization, team *gen.Team) (*gen.Team, error)

	GetTeams(ctx context.Context, org *gen.Organization) ([]*gen.Team, error)

	// Add user to team

	AddUserToTeam(ctx context.Context, team *gen.Team, user *gen.User) error
}
