package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	codefly "github.com/codefly-dev/sdk-go"
	"github.com/jackc/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/jackc/pgx/v4"

	"backend/pkg/gen"

	"github.com/codefly-dev/core/wool"

	"github.com/jackc/pgx/v4/pgxpool"

	"backend/pkg/business"
)

type Close func()

type PostgresStore struct {
	Close
	pool *pgxpool.Pool
}

var _ business.Store = (*PostgresStore)(nil)

func (s *PostgresStore) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Begin transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails
	defer tx.Rollback(ctx)

	// Create a new context with the transaction
	txCtx := context.WithValue(ctx, "tx", tx)

	// Run the provided function
	if err := fn(txCtx); err != nil {
		// If there's an error, rollback and return the error
		return err
	}

	// If everything succeeded, commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type QueryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (s *PostgresStore) getQueryExecutor(ctx context.Context) QueryExecutor {
	tx, ok := ctx.Value("tx").(pgx.Tx)
	if ok {
		return tx
	}
	return s.pool
}

func (s *PostgresStore) CreateUser(ctx context.Context, user *gen.User) error {
	w := wool.Get(ctx).In("CreateUser")

	now := time.Now().UTC()

	sql := `
        INSERT INTO users (id, status, signed_up_at, last_login_at, email, profile)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	profileJSON, err := json.Marshal(user.Profile)
	if err != nil {
		return w.Wrapf(err, "error marshaling user profile")
	}

	args := []interface{}{
		user.Id,
		user.Status,
		now,
		now,
		user.Email,
		profileJSON,
	}

	_, err = s.getQueryExecutor(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return w.Wrapf(err, "error creating user")
	}

	return nil
}

func (s *PostgresStore) LinkUserWithAuth(ctx context.Context, id string, authID string) error {
	w := wool.Get(ctx).In("LinkUserWithAuth")

	sql := `
		INSERT INTO users_auth (user_id, auth_id)
		VALUES ($1, $2)
	`

	_, err := s.getQueryExecutor(ctx).Exec(ctx, sql, id, authID)
	if err != nil {
		return w.Wrapf(err, "error linking user with auth")
	}

	return nil
}

func (s *PostgresStore) GetUserByAuthId(ctx context.Context, authID string) (*gen.User, error) {
	w := wool.Get(ctx).In("GetUserByAuthId")

	// First, get the user_id from users_auth table
	sql := `
        SELECT user_id
        FROM users_auth
        WHERE auth_id = $1
    `

	var id string
	err := s.getQueryExecutor(ctx).QueryRow(ctx, sql, authID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found, return nil without error
		}
		return nil, w.Wrapf(err, "error querying users_auth table")
	}

	sql = `
        SELECT id, status, signed_up_at, last_login_at, email, profile
        FROM users
        WHERE id = $1
    `

	var user gen.User
	var profileJSON []byte
	var signedUpAt, lastLoginAt time.Time

	err = s.getQueryExecutor(ctx).QueryRow(ctx, sql, id).Scan(
		&user.Id,
		&user.Status,
		&signedUpAt,
		&lastLoginAt,
		&user.Email,
		&profileJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found, return nil without error
		}
		return nil, w.Wrapf(err, "error querying user")
	}

	user.SignedUpAt = timestamppb.New(signedUpAt)
	user.LastLoginAt = timestamppb.New(lastLoginAt)

	// Unmarshal profile JSON
	if len(profileJSON) > 0 {
		var profile gen.UserProfile
		if err := json.Unmarshal(profileJSON, &profile); err != nil {
			return nil, w.Wrapf(err, "error unmarshaling user profile")
		}
		user.Profile = &profile
	}

	return &user, nil
}

func (s *PostgresStore) GetUserById(ctx context.Context, id string) (*gen.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) DeleteUser(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) UpdateUser(ctx context.Context, user *gen.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) CreateOrganization(ctx context.Context, org *gen.Organization) error {
	w := wool.Get(ctx).In("CreateOrganization")

	sql := `
        INSERT INTO organizations (id, name, domain)
        VALUES ($1, $2, $3)
    `

	_, err := s.getQueryExecutor(ctx).Exec(ctx, sql, org.Id, org.Name, org.Domain)
	if err != nil {
		return w.Wrapf(err, "error creating organization")
	}

	return nil
}

func (s *PostgresStore) GetOrganization(ctx context.Context, id string) (*gen.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) DeleteOrganization(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) UpdateOrganization(ctx context.Context, org *gen.Organization) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) CreateTeam(ctx context.Context, orgID string, team *gen.Team) error {
	w := wool.Get(ctx).In("CreateTeam")

	sql := `
        INSERT INTO teams (id, name, organization_id)
        VALUES ($1, $2, $3)
    `

	_, err := s.getQueryExecutor(ctx).Exec(ctx, sql, team.Id, team.Name, orgID)
	if err != nil {
		return w.Wrapf(err, "error creating team")
	}

	// If the team has members, add them to the team
	if len(team.MemberIds) > 0 {
		insertMemberSQL := `
            INSERT INTO user_teams (team_id, user_id)
            VALUES ($1, $2)
        `

		for _, memberID := range team.MemberIds {
			_, err := s.getQueryExecutor(ctx).Exec(ctx, insertMemberSQL, team.Id, memberID)
			if err != nil {
				return w.Wrapf(err, "error adding member to team")
			}
		}
	}

	return nil
}

// CreateRole creates a new role in the database
func (s *PostgresStore) CreateRole(ctx context.Context, role *gen.Role) error {
	w := wool.Get(ctx).In("CreateRole")

	sql := `
        INSERT INTO roles (id, name)
        VALUES ($1, $2)
    `

	_, err := s.getQueryExecutor(ctx).Exec(ctx, sql, role.Id, role.Name)
	if err != nil {
		return w.Wrapf(err, "error creating role")
	}

	return nil
}

// CreatePermission creates a new permission in the database
func (s *PostgresStore) CreatePermission(ctx context.Context, permission *gen.Permission) error {
	w := wool.Get(ctx).In("CreatePermission")

	sql := `
        INSERT INTO permissions (id, name, resource, access)
        VALUES ($1, $2, $3, $4)
    `

	_, err := s.getQueryExecutor(ctx).Exec(ctx, sql, permission.Id, permission.Name, permission.Resource, permission.Access)
	if err != nil {
		return w.Wrapf(err, "error creating permission")
	}

	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (s *PostgresStore) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	w := wool.Get(ctx).In("AssignPermissionToRole")

	// Check if the assignment already exists
	checkSQL := `
        SELECT EXISTS (
            SELECT 1 FROM role_permissions 
            WHERE role_id = $1 AND permission_id = $2
        )
    `
	var exists bool
	err := s.getQueryExecutor(ctx).QueryRow(ctx, checkSQL, roleID, permissionID).Scan(&exists)
	if err != nil {
		return w.Wrapf(err, "error checking existing permission assignment")
	}
	if exists {
		return w.Wrapf(err, "permission is already assigned to role")
	}

	// If not exists, proceed with insertion
	insertSQL := `
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
    `

	_, err = s.getQueryExecutor(ctx).Exec(ctx, insertSQL, roleID, permissionID)
	if err != nil {
		return w.Wrapf(err, "error assigning permission to role")
	}
	return nil
}

func (s *PostgresStore) GetTeams(ctx context.Context, orgID string) ([]*gen.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetTeamByID(ctx context.Context, teamID string) (*gen.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) DeleteTeam(ctx context.Context, teamID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) UpdateTeam(ctx context.Context, team *gen.Team) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) AddUserToOrganization(ctx context.Context, orgID string, userID string, roleID string) error {
	w := wool.Get(ctx).In("AddUserToOrganization")
	// Check if the user is already in the organization
	checkSQL := `
        SELECT EXISTS (
            SELECT 1 FROM organization_users 
            WHERE organization_id = $1 AND user_id = $2
        )
    `
	var exists bool
	err := s.getQueryExecutor(ctx).QueryRow(ctx, checkSQL, orgID, userID).Scan(&exists)
	if err != nil {
		return w.Wrapf(err, "error checking existing user in organization")
	}
	if exists {
		return w.Wrapf(err, "user is already in the organization")
	}

	sql := `
        INSERT INTO organization_users (organization_id, user_id, role, joined_at)
        VALUES ($1, $2, $3, $4)
    `

	joinedAt := time.Now().UTC()
	_, err = s.getQueryExecutor(ctx).Exec(ctx, sql, orgID, userID, roleID, joinedAt)
	if err != nil {
		return w.Wrapf(err, "error adding user to organization")
	}
	return nil
}

func (s *PostgresStore) RemoveUserFromOrganization(ctx context.Context, orgID string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetUsersInOrganization(ctx context.Context, orgID string) ([]*gen.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) AddUserToTeam(ctx context.Context, teamID string, userID string, roleID string) error {
	w := wool.Get(ctx).In("AddUserToTeam")

	// Validate team exists
	if err := s.validateEntityExists(ctx, "teams", teamID); err != nil {
		return w.Wrapf(err, "invalid team ID")
	}

	// Validate user exists
	if err := s.validateEntityExists(ctx, "users", userID); err != nil {
		return w.Wrapf(err, "invalid user ID")
	}

	// Validate role exists
	if err := s.validateEntityExists(ctx, "roles", roleID); err != nil {
		return w.Wrapf(err, "invalid role ID")
	}

	// First, check if the user is already in the team
	checkSQL := `
        SELECT EXISTS (
            SELECT 1 FROM team_members 
            WHERE team_id = $1 AND user_id = $2
        )
    `
	var exists bool
	err := s.getQueryExecutor(ctx).QueryRow(ctx, checkSQL, teamID, userID).Scan(&exists)
	if err != nil {
		return w.Wrapf(err, "error checking existing user in team")
	}
	if exists {
		return w.Wrapf(nil, "user is already in the team")
	}

	// If the user is not in the team, add them
	insertSQL := `
        INSERT INTO team_members (team_id, user_id, role_id, joined_at)
        VALUES ($1, $2, $3, $4)
    `

	joinedAt := time.Now().UTC()
	_, err = s.getQueryExecutor(ctx).Exec(ctx, insertSQL, teamID, userID, roleID, joinedAt)
	if err != nil {
		return w.Wrapf(err, "error adding user to team")
	}
	return nil

}

func (s *PostgresStore) validateEntityExists(ctx context.Context, table, id string) error {
	sql := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE id = $1)", table)
	var exists bool
	err := s.getQueryExecutor(ctx).QueryRow(ctx, sql, id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s with id %s does not exist", table, id)
	}
	return nil
}

func (s *PostgresStore) RemoveUserFromTeam(ctx context.Context, teamID string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetUsersInTeam(ctx context.Context, teamID string) ([]*gen.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetRoleByID(ctx context.Context, roleID string) (*gen.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) UpdateRole(ctx context.Context, role *gen.Role) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) DeleteRole(ctx context.Context, roleID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetAllRoles(ctx context.Context) ([]*gen.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetPermissionByID(ctx context.Context, permissionID string) (*gen.Permission, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) UpdatePermission(ctx context.Context, permission *gen.Permission) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) DeletePermission(ctx context.Context, permissionID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetAllPermissions(ctx context.Context) ([]*gen.Permission, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetPermissionsForRole(ctx context.Context, roleID string) ([]*gen.Permission, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) HasPermission(ctx context.Context, userID string, resourceType string, resourceID string, permissionName string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStore) GetUserRoles(ctx context.Context, userID string, resourceType string, resourceID string) ([]*gen.Role, error) {
	//TODO implement me
	panic("implement me")
}

//	func (p *PostgresStore) GetOrganization(ctx context.Context, id string) (*gen.Organization, error) {
//		w := wool.Get(ctx).In("GetOrganization")
//		sql := `SELECT id, name FROM organizations WHERE id = $1`
//		row := p.pool.QueryRow(ctx, sql, id)
//		var id string
//		var name string
//		err := row.Scan(&id, &name)
//		if err != nil {
//			if errors.Is(err, pgx.ErrNoRows) {
//				return nil, nil
//			}
//			return nil, w.Wrapf(err, "cannot get organization")
//		}
//		return &gen.Organization{
//			Id:   id,
//			Name: name,
//		}, nil
//	}
//
//	func (p *PostgresStore) CreateOrganization(ctx context.Context, owner *gen.User, org *gen.Organization) error {
//		w := wool.Get(ctx).In("CreateOrganization")
//		sql := `INSERT INTO organizations (id, name, owner) VALUES ($1, $2, $3)`
//		args := []any{org.Id, org.Name, owner.Id}
//		_, err := p.pool.Exec(ctx, sql, args...)
//		if err != nil {
//			return w.Wrapf(err, "cannot create organization")
//		}
//		return nil
//	}
//
//	func (p *PostgresStore) DeleteOrganization(ctx context.Context, id string) error {
//		w := wool.Get(ctx).In("DeleteOrganization")
//		sql := `DELETE FROM organizations WHERE id = $1`
//		_, err := p.pool.Exec(ctx, sql, id)
//		if err != nil {
//			return w.Wrapf(err, "cannot delete organization")
//		}
//		return nil
//	}
//
//	func (p *PostgresStore) DeleteUser(ctx context.Context, id string) error {
//		sql := `DELETE FROM users WHERE id = $1`
//		_, err := p.pool.Exec(ctx, sql, id)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func (p *PostgresStore) CreateTeam(ctx context.Context, team *gen.Team) error {
//		sql := `INSERT INTO teams (id, name, organization_id) VALUES ($1, $2, $3)`
//		args := []any{team.Id, team.Name, org.Id}
//		_, err := p.pool.Exec(ctx, sql, args...)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func (p *PostgresStore) GetTeams(ctx context.Context, org *gen.Organization) ([]*gen.Team, error) {
//		sql := `SELECT id, name FROM teams WHERE organization_id = $1`
//		rows, err := p.pool.Query(ctx, sql, org.Id)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//
//		var teams []*gen.Team
//		for rows.Next() {
//			var team gen.Team
//			err := rows.Scan(&team.Id, &team.Name)
//			if err != nil {
//				return nil, err
//			}
//			teams = append(teams, &team)
//		}
//		if err = rows.Err(); err != nil {
//			return nil, err
//		}
//		return teams, nil
//	}
//
//	func (p *PostgresStore) DeleteTeam(ctx context.Context, team *gen.Team) error {
//		sql := `DELETE FROM teams WHERE id = $1`
//		_, err := p.pool.Exec(ctx, sql, team.Id)
//		if err != nil {
//			return fmt.Errorf("error deleting team: %w", err)
//		}
//		return nil
//	}
//
//	func (p *PostgresStore) AddUserToTeam(ctx context.Context, team *gen.Team, user *gen.User) error {
//		sql := `INSERT INTO team_users (team_id, user_id) VALUES ($1, $2)`
//		args := []any{team.Id, user.Id}
//		_, err := p.pool.Exec(ctx, sql, args...)
//		return err
//	}
//
//	func (p *PostgresStore) assignPermission(ctx context.Context, table string, entityId string, permission *gen.Permission) error {
//		// Insert the permission into the permissions table if it does not exist
//		sql := `INSERT INTO permissions (id, name, resource, access) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`
//		args := []any{permission.Id, permission.Name, permission.Resource, permission.Access}
//		_, err := p.pool.Exec(ctx, sql, args...)
//		if err != nil {
//			return fmt.Errorf("error inserting permission: %w", err)
//		}
//
//		// Insert the association into the appropriate table
//		sql = fmt.Sprintf(`INSERT INTO %s_permissions (entity_id, permission_id) VALUES ($1, $2)`, table)
//		args = []any{entityId, permission.Id}
//		_, err = p.pool.Exec(ctx, sql, args...)
//		if err != nil {
//			return fmt.Errorf("error assigning permission: %w", err)
//		}
//		return nil
//	}
//
//	func (p *PostgresStore) AssignTeamPermission(ctx context.Context, team *gen.Team, permission *gen.Permission) error {
//		return p.assignPermission(ctx, "team", team.Id, permission)
//	}
//
//	func (p *PostgresStore) AssignUserPermission(ctx context.Context, user *gen.User, permission *gen.Permission) error {
//		return p.assignPermission(ctx, "user", user.Id, permission)
//	}
func NewPostgresStore(ctx context.Context) (*PostgresStore, error) {
	w := wool.Get(ctx).In("NewPostgresStore")
	connection, err := codefly.For(ctx).Service("store").Secret("postgres", "connection")
	if err != nil {
		return nil, w.Wrapf(err, "failed to get connection string")
	}

	poolConfig, err := pgxpool.ParseConfig(connection)
	if err != nil {
		return nil, w.Wrapf(err, "failed to parse connection string")
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, w.Wrapf(err, "failed to connect to database")
	}
	store := &PostgresStore{
		pool:  pool,
		Close: pool.Close,
	}
	return store, nil
}
