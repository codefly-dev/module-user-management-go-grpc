package business

import (
	"backend/pkg/gen"
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type Status string

const (
	Active  Status = "active"
	Pending Status = "pending"
)

const DefaultOrganizationName = "My Organization"

type Service struct {
	store Store
}

func NewService(store Store) (*Service, error) {
	return &Service{store: store}, nil
}

func (s *Service) RegisterUser(ctx context.Context, user *gen.User) (*gen.RegisterUserResponse, error) {
	// Invitation only
	if slices.Contains(invited, user.Email) {
		user.Status = string(Active)
	} else {
		user.Status = string(Pending)
	}
	// If already exists, fails
	if u, err := s.store.GetUserByAuthID(ctx, user.SignupAuthId); err != nil {
		return nil, err
	} else if u != nil {
		return nil, fmt.Errorf("user already exists: %s", user.Email)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating user id: %w", err)
	}
	user.Id = id.String()
	u, err := s.store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	// Create organization
	org, err := s.CreateOrganization(ctx, u)
	if err != nil {
		return nil, err
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
func (s *Service) CreateOrganization(ctx context.Context, u *gen.User) (*gen.Organization, error) {
	// Check if we already have a organization owned by this user
	org, err := s.store.GetOrganizationForOwner(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error getting organization: %w", err)
	}
	if org != nil {
		return nil, fmt.Errorf("organization already exists: %s", u.Email)
	}
	orgID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating organization id: %w", err)
	}
	org = &gen.Organization{
		Id:   orgID.String(),
		Name: DefaultOrganizationName,
	}
	org, err = s.store.CreateOrganization(ctx, u, org)
	if err != nil {
		return nil, fmt.Errorf("error creating organization: %w", err)
	}
	return org, nil
}

func (s *Service) SetStore(store Store) {
	s.store = store
}

func (s *Service) DeleteUser(ctx context.Context, authSignupId string) (*gen.User, error) {
	u, err := s.store.DeleteUser(ctx, authSignupId)
	if err != nil {
		return nil, fmt.Errorf("error deleting user: %w", err)
	}
	if u == nil {
		return nil, nil
	}
	// Delete organization for the user
	org, err := s.GetOrganizationForOwner(ctx, u)
	if err != nil {
		return u, fmt.Errorf("error getting organization: %w", err)
	}
	if org == nil {
		return u, nil
	}
	err = s.DeleteOrganization(ctx, org)
	return u, err
}

func (s *Service) DeleteOrganization(ctx context.Context, org *gen.Organization) error {
	return s.store.DeleteOrganization(ctx, org)
}

// Right now hardcode email of invited users
var invited = []string{"antoine.toussaint@codefly.ai"}
