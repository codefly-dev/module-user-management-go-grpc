package business

import (
	"backend/pkg/gen"
	"context"
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

type RegisterUserInput struct {
	Email        string
	SignupAuthId string
}

// RegisterUser creates a new user with a new organization
// Input is the AuthID of the user that we get from the context
func (s *Service) RegisterUser(ctx context.Context, input *RegisterUserInput) (*gen.RegisterUserResponse, error) {
	w := wool.Get(ctx).In("RegisterUser")

	// If already exists, fails
	if u, err := s.store.GetUserByAuthId(ctx, input.SignupAuthId); err != nil {
		return nil, w.Wrapf(err, "error getting user")
	} else if u != nil {
		return nil, w.NewError("user already exists: %s", input.Email)
	}

	user := &gen.User{
		Id:    uuid.New().String(),
		Email: input.Email,
	}
	// Invitation only
	if slices.Contains(invited, user.Email) {
		user.Status = string(Active)
	} else {
		user.Status = string(Pending)
	}

	// Create organization
	org := &gen.Organization{
		Id:   uuid.New().String(),
		Name: DefaultOrganizationName,
	}

	adminTeam := &gen.Team{
		Id:   uuid.New().String(),
		Name: AdminTeamName,
	}

	adminRole := &gen.Role{
		Id:   uuid.New().String(),
		Name: "Admin",
	}

	// Define admin permissions
	adminPermissions := []*gen.Permission{
		{Id: uuid.New().String(), Name: "manage_users", Resource: "*", Access: "write"},
		{Id: uuid.New().String(), Name: "manage_teams", Resource: "*", Access: "write"},
		{Id: uuid.New().String(), Name: "manage_organization", Resource: "*", Access: "write"},
		// Add more permissions as needed
	}

	err := s.store.RunInTransaction(ctx, func(ctx context.Context) error {

		// Create admin role
		if err := s.store.CreateRole(ctx, adminRole); err != nil {
			return w.Wrapf(err, "cannot create admin role")
		}

		// Assign permissions to admin role
		for _, perm := range adminPermissions {
			if err := s.store.CreatePermission(ctx, perm); err != nil {
				return w.Wrapf(err, "cannot create permission")
			}
			if err := s.store.AssignPermissionToRole(ctx, adminRole.Id, perm.Id); err != nil {
				return w.Wrapf(err, "cannot assign permission to role")
			}
		}

		// Create user
		if err := s.store.CreateUser(ctx, user); err != nil {
			return w.Wrapf(err, "cannot create user")
		}

		// Link user to auth id
		if err := s.store.LinkUserWithAuth(ctx, user.Id, input.SignupAuthId); err != nil {
			return w.Wrapf(err, "cannot link user with auth id")
		}

		// Create organization
		if err := s.store.CreateOrganization(ctx, org); err != nil {
			return w.Wrapf(err, "cannot create organization")
		}

		// Add user to organization
		if err := s.store.AddUserToOrganization(ctx, org.Id, user.Id, adminRole.Id); err != nil {
			return w.Wrapf(err, "cannot add user to organization")
		}

		// Create admin team
		if err := s.store.CreateTeam(ctx, org.Id, adminTeam); err != nil {
			return w.Wrapf(err, "cannot create admin team")
		}

		// Add user to admin team with admin role
		if err := s.store.AddUserToTeam(ctx, adminTeam.Id, user.Id, adminRole.Id); err != nil {
			return w.Wrapf(err, "cannot add user to admin team")
		}
		return nil
	})
	if err != nil {
		return nil, w.Wrapf(err, "error registrating user transaction")
	}
	return &gen.RegisterUserResponse{User: user, Organization: org}, nil
}

func (s *Service) GetUserByAuthId(ctx context.Context, id string) (*gen.User, error) {
	return s.store.GetUserByAuthId(ctx, id)
}

func (s *Service) GetUserById(ctx context.Context, userId string) (*gen.User, error) {
	return s.store.GetUserById(ctx, userId)
}

func (s *Service) CreateUser(ctx context.Context, guest *gen.User) error {
	return s.store.CreateUser(ctx, guest)
}

//	func (s *Service) GetOrganization(ctx context.Context, u *gen.User) (*gen.Organization, error) {
//		return s.store.GetOrganization(ctx, u)
//	}
//
// // CreateOrganization creates an organization for the user
//
//	func (s *Service) CreateOrganization(ctx context.Context, org *gen.Organization) error {
//		w := wool.Get(ctx).In("CreateOrganization")
//		err := s.store.CreateOrganization(ctx, org)
//		if err != nil {
//			return fmt.Errorf("error creating organization: %w", err)
//		}
//		return nil
//	}
//
//	func (s *Service) DeleteUser(ctx context.Context, id string) error {
//		w := wool.Get(ctx).In("DeleteOwner")
//		u, err := s.GetUserByAuthId(ctx, id)
//		if err != nil {
//			return w.Wrapf(err, "error getting user")
//		}
//		if u == nil {
//			return w.NewError("not found: user")
//		}
//		org, err := s.GetOrganization(ctx, u)
//		if err != nil {
//			return w.Wrapf(err, "can't get organization for owner")
//		}
//		if org != nil {
//			err = s.DeleteOrganization(ctx, org)
//			if err != nil {
//				return w.Wrapf(err, "can't delete organization")
//			}
//		}
//		return s.store.DeleteUser(ctx, u)
//	}
//
//	func (s *Service) DeleteOrganization(ctx context.Context, org *gen.Organization) error {
//		w := wool.Get(ctx).In("DeleteOrganization")
//		// Delete teams for this organizations
//		teams, err := s.GetTeams(ctx, org)
//		if err != nil {
//			return w.Wrapf(err, "can't get teams for organization")
//		}
//		for _, team := range teams {
//			err = s.DeleteTeam(ctx, team)
//			if err != nil {
//				return w.Wrapf(err, "can't delete team")
//			}
//		}
//		return s.store.DeleteOrganization(ctx, org)
//	}
//
//	func (s *Service) GetTeams(ctx context.Context, org *gen.Organization) ([]*gen.Team, error) {
//		return s.store.GetTeams(ctx, org)
//
// }
//
//	func (s *Service) CreateTeam(ctx context.Context, org *gen.Organization, team *gen.Team) error {
//		return s.store.CreateTeam(ctx, org, team)
//	}
//
//	func (s *Service) DeleteTeam(ctx context.Context, team *gen.Team) error {
//		return s.store.DeleteTeam(ctx, team)
//	}
//
//	func (s *Service) AddUserToTeam(ctx context.Context, team *gen.Team, user *gen.User) error {
//		return s.store.AddUserToTeam(ctx, team, user)
//	}
//
// Right now hardcode email of invited users
var invited = []string{"antoine.toussaint@codefly.ai"}
