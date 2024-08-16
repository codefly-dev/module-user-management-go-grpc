package business

import (
	"context"

	"backend/pkg/gen"
)

type Store interface {

	// Run in transaction

	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	// User manipulation

	CreateUser(ctx context.Context, user *gen.User) error
	LinkUserWithAuth(ctx context.Context, id string, authID string) error
	GetUserByAuthId(ctx context.Context, id string) (*gen.User, error)
	GetUserById(ctx context.Context, id string) (*gen.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *gen.User) error

	// Organization

	CreateOrganization(ctx context.Context, org *gen.Organization) error
	GetOrganization(ctx context.Context, id string) (*gen.Organization, error)
	DeleteOrganization(ctx context.Context, id string) error
	UpdateOrganization(ctx context.Context, org *gen.Organization) error

	// Teams

	CreateTeam(ctx context.Context, orgID string, team *gen.Team) error
	GetTeams(ctx context.Context, orgID string) ([]*gen.Team, error)
	GetTeamByID(ctx context.Context, teamID string) (*gen.Team, error)
	DeleteTeam(ctx context.Context, teamID string) error
	UpdateTeam(ctx context.Context, team *gen.Team) error

	// User-Organization relationships

	AddUserToOrganization(ctx context.Context, orgID string, userID string, roleID string) error
	RemoveUserFromOrganization(ctx context.Context, orgID string, userID string) error
	GetUsersInOrganization(ctx context.Context, orgID string) ([]*gen.User, error)

	// User-Team relationships

	AddUserToTeam(ctx context.Context, teamID string, userID string, roleID string) error
	RemoveUserFromTeam(ctx context.Context, teamID string, userID string) error
	GetUsersInTeam(ctx context.Context, teamID string) ([]*gen.User, error)

	// Roles and Permissions

	CreateRole(ctx context.Context, role *gen.Role) error
	GetRoleByID(ctx context.Context, roleID string) (*gen.Role, error)
	UpdateRole(ctx context.Context, role *gen.Role) error
	DeleteRole(ctx context.Context, roleID string) error
	GetAllRoles(ctx context.Context) ([]*gen.Role, error)

	CreatePermission(ctx context.Context, permission *gen.Permission) error
	GetPermissionByID(ctx context.Context, permissionID string) (*gen.Permission, error)
	UpdatePermission(ctx context.Context, permission *gen.Permission) error
	DeletePermission(ctx context.Context, permissionID string) error
	GetAllPermissions(ctx context.Context) ([]*gen.Permission, error)

	AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
	GetPermissionsForRole(ctx context.Context, roleID string) ([]*gen.Permission, error)
}
