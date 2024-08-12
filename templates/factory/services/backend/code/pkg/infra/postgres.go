package infra

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	codefly "github.com/codefly-dev/sdk-go"

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

func (p *PostgresStore) CreateUser(ctx context.Context, user *gen.User) (*gen.User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	user.Id = id.String()
	now := time.Now()
	sql := `INSERT INTO users (id, status, auth_signup_id, signed_up_at, email, profile) VALUES ($1, $2, $3, $4, $5, $6)`
	args := []any{user.Id, user.Status, user.SignupAuthId, now, user.Email, user.Profile}
	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresStore) GetUserByAuthID(ctx context.Context, id string) (*gen.User, error) {
	sql := `SELECT id, email  FROM users WHERE auth_signup_id = $1`
	row := p.pool.QueryRow(ctx, sql, id)
	var uid string
	var email string
	err := row.Scan(&uid, &email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &gen.User{
		Id:           uid,
		SignupAuthId: id,
		Email:        email,
	}, nil
}

func (p *PostgresStore) GetOrganizationForOwner(ctx context.Context, user *gen.User) (*gen.Organization, error) {
	sql := `SELECT id, name FROM organizations WHERE owner = $1`
	row := p.pool.QueryRow(ctx, sql, user.Id)
	var id string
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &gen.Organization{
		Id:   id,
		Name: name,
	}, nil
}

func (p *PostgresStore) CreateOrganization(ctx context.Context, owner *gen.User, org *gen.Organization) (*gen.Organization, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	org.Id = id.String()
	sql := `INSERT INTO organizations (id, name, owner) VALUES ($1, $2, $3)`
	args := []any{org.Id, org.Name, owner.Id}
	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (p *PostgresStore) DeleteOrganization(ctx context.Context, org *gen.Organization) error {
	sql := `DELETE FROM organizations WHERE id = $1`
	_, err := p.pool.Exec(ctx, sql, org.Id)
	return err
}

func (p *PostgresStore) DeleteUser(ctx context.Context, authSignupId string) (*gen.User, error) {
	u, err := p.GetUserByAuthID(ctx, authSignupId)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if u == nil {
		return nil, nil
	}
	sql := `DELETE FROM users WHERE auth_signup_id = $1`
	_, err = p.pool.Exec(ctx, sql, u.SignupAuthId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (p *PostgresStore) CreateTeam(ctx context.Context, org *gen.Organization, team *gen.Team) (*gen.Team, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	team.Id = id.String()
	sql := `INSERT INTO teams (id, name, organization_id) VALUES ($1, $2, $3)`
	args := []any{team.Id, team.Name, org.Id}
	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (p *PostgresStore) GetTeams(ctx context.Context, org *gen.Organization) ([]*gen.Team, error) {
	sql := `SELECT id, name FROM teams WHERE organization_id = $1`
	rows, err := p.pool.Query(ctx, sql, org.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*gen.Team
	for rows.Next() {
		var team gen.Team
		err := rows.Scan(&team.Id, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func (p *PostgresStore) AddUserToTeam(ctx context.Context, team *gen.Team, user *gen.User) error {
	sql := `INSERT INTO team_users (team_id, user_id) VALUES ($1, $2)`
	args := []any{team.Id, user.Id}
	_, err := p.pool.Exec(ctx, sql, args...)
	return err
}

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
