package business

import (
	"backend/pkg/gen"
	"context"
	"fmt"
	"github.com/codefly-dev/core/wool"
	"slices"

	"github.com/google/uuid"
)

type Status string

func (s *Service) SetStore(store Store) {
	s.store = store
}

const (
	Active  Status = "active"
	Pending Status = "pending"
)

const DefaultOrganizationName = "My Organization"

const AdminTeamName = "Administrators"

type Service struct {
	store Store
}

func NewService(store Store) (*Service, error) {
	return &Service{store: store}, nil
}

func (s *Service) RegisterUser(ctx context.Context, user *gen.User) (*gen.RegisterUserResponse, error) {
	w := wool.Get(ctx).In("RegisterUser")
	// Invitation only
	if slices.Contains(invited, user.Email) {
		user.Status = string(Active)
	} else {
		user.Status = string(Pending)
	}
	// If already exists, fails
	if u, err := s.store.GetUserByAuthID(ctx, user.SignupAuthId); err != nil {
		return nil, w.Wrapf(err, "error getting user")
	} else if u != nil {
		return nil, w.NewError("user already exists: %s", user.Email)
	}

	u, err := s.store.CreateUser(ctx, user)
	if err != nil {
		return nil, w.Wrapf(err, "error creating user")
	}
	if u == nil {
		return nil, w.NewError("error creating user")
	}

	// Create organization
	org, err := s.CreateOrganization(ctx, u, &gen.Organization{Name: DefaultOrganizationName})
	if err != nil {
		return nil, w.Wrapf(err, "error creating organization")
	}
	if org == nil {
		return nil, w.NewError("error creating organization")
	}

	// Create Admin Team
	team, err := s.store.CreateTeam(ctx, org, &gen.Team{Name: AdminTeamName})
	if err != nil {
		return nil, w.Wrapf(err, "error creating team")
	}
	// Add Admin permissions this team
	// TODO

	// Adding the User to this team
	err = s.store.AddUserToTeam(ctx, team, u)
	if err != nil {
		return nil, w.Wrapf(err, "error adding user to team")
	}
	return &gen.RegisterUserResponse{
		User:         u,
		Organization: org,
	}, nil
}

func (s *Service) GetUserByAuthID(ctx context.Context, id string) (*gen.User, error) {
	return s.store.GetUserByAuthID(ctx, id)
}

func (s *Service) GetOrganizationForOwner(ctx context.Context, u *gen.User) (*gen.Organization, error) {
	return s.store.GetOrganizationForOwner(ctx, u)
}

// CreateOrganization creates an organization for the user
func (s *Service) CreateOrganization(ctx context.Context, u *gen.User, org *gen.Organization) (*gen.Organization, error) {
	w := wool.Get(ctx).In("CreateOrganization")
	// Check if we already have an organization owned by this user
	exists, err := s.store.GetOrganizationForOwner(ctx, u)
	if err != nil {
		return nil, w.Wrapf(err, "error getting organization")
	}
	if exists != nil {
		return nil, w.NewError("organization already exists")
	}
	// Otherwise, create it
	orgID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating organization id: %w", err)
	}
	org.Id = orgID.String()
	org, err = s.store.CreateOrganization(ctx, u, org)
	if err != nil {
		return nil, fmt.Errorf("error creating organization: %w", err)
	}
	return org, nil
}

func (s *Service) DeleteOwner(ctx context.Context, authSignupId string) (*gen.User, error) {
	w := wool.Get(ctx).In("DeleteOwner")
	u, err := s.GetUserByAuthID(ctx, authSignupId)
	if err != nil {
		return nil, w.Wrapf(err, "error getting user")
	}
	if u == nil {
		return nil, nil
	}
	org, err := s.GetOrganizationForOwner(ctx, u)
	if err != nil {
		return u, fmt.Errorf("error getting organization: %w", err)
	}
	if org != nil {
		err = s.DeleteOrganization(ctx, org)
		if err != nil {
			return u, fmt.Errorf("error deleting organization: %w", err)
		}
	}
	return s.store.DeleteUser(ctx, authSignupId)
}

func (s *Service) DeleteOrganization(ctx context.Context, org *gen.Organization) error {
	return s.store.DeleteOrganization(ctx, org)
}

func (s *Service) GetTeams(ctx context.Context, org *gen.Organization) ([]*gen.Team, error) {
	return s.store.GetTeams(ctx, org)

}

// Right now hardcode email of invited users
var invited = []string{"antoine.toussaint@codefly.ai"}
