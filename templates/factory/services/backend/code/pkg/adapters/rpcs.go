package adapters

import (
	"context"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/codefly-dev/core/wool"

	"backend/pkg/business"
	"backend/pkg/gen"
)

var service *business.Service

func WithService(s *business.Service) {
	service = s
}

func (s *GrpcServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("Login").GRPC().Inject()
	id, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	user, err := service.GetUserByAuthID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get user: %v", err)
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &gen.LoginResponse{User: user}, nil
}

func (s *GrpcServer) GetOrganization(ctx context.Context, req *gen.GetOrganizationRequest) (*gen.Organization, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("GetOrganization").GRPC().Inject()

	id, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	user, err := service.GetUserByAuthID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user: %v", err)
	}
	return service.GetOrganizationForOwner(ctx, user)
}

func (s *GrpcServer) UpdateOrganization(ctx context.Context, req *gen.UpdateOrganizationRequest) (*gen.UpdateOrganizationResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	//w := wool.Get(ctx).In("CreateOrganization")
	//email, w.UserEmail()
	//organization, err := service.CreateOrganization(ctx, request, request)
	return nil, nil
}

func (s *GrpcServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterUserResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("CreateSelf").GRPC().Inject()
	w.Info("in register")
	id, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	w.Info("got user id", wool.Field("id", id))
	email, found := w.UserEmail()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user email not found")
	}
	// Optional
	name, _ := w.UserName()
	givenName, _ := w.UserGivenName()

	var user *gen.User
	var err error

	if user, err = service.GetUserByAuthID(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user: %v", err)
	}
	if user != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}
	user = &gen.User{
		SignupAuthId: id,
		Email:        email,
		Profile: &gen.UserProfile{
			Name:      name,
			GivenName: givenName,
		},
	}
	resp, err := service.RegisterUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}
	return resp, nil
}

func (s *GrpcServer) CreateTeam(ctx context.Context, req *gen.CreateTeamRequest) (*gen.CreateTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GrpcServer) AddUserToTeam(ctx context.Context, req *gen.AddUserToTeamRequest) (*gen.AddUserToTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GrpcServer) RemoveUserFromTeam(ctx context.Context, req *gen.RemoveUserFromTeamRequest) (*gen.RemoveUserFromTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GrpcServer) AssignPermission(ctx context.Context, req *gen.AssignPermissionRequest) (*gen.AssignPermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GrpcServer) RevokePermission(ctx context.Context, req *gen.RevokePermissionRequest) (*gen.RevokePermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}
