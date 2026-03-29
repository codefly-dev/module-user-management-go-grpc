package adapters

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/codefly-dev/core/wool"

	"backend/pkg/business"
	"backend/pkg/gen"
	"backend/pkg/infra"
)

var service *business.Service

func WithService(s *business.Service) {
	service = s
}

// ============================================================================
// UserService RPCs (on UserServer)
// ============================================================================

func (s *UserServer) GetSelf(ctx context.Context, _ *gen.GetSelfRequest) (*gen.GetSelfResponse, error) {
	w := wool.Get(ctx).In("GetSelf")
	w.GRPC().Inject()

	// X-Auth-Id is set by the gateway/sidecar
	userID, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found in headers")
	}
	_ = userID
	return nil, status.Error(codes.Unimplemented, "GetSelf not yet implemented")
}

func (s *UserServer) RegisterUser(ctx context.Context, req *gen.RegisterUserRequest) (*gen.RegisterUserResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.RegisterUser(ctx, req)
}

func (s *UserServer) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.User, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "GetUser not yet implemented")
}

func (s *UserServer) ListUsers(ctx context.Context, _ *gen.ListUsersRequest) (*gen.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListUsers not yet implemented")
}

func (s *UserServer) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*gen.User, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "UpdateUser not yet implemented")
}

func (s *UserServer) DeleteUser(ctx context.Context, req *gen.GetUserRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "DeleteUser not yet implemented")
}

func (s *UserServer) AddIdentity(ctx context.Context, req *gen.AddIdentityRequest) (*gen.UserIdentity, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "AddIdentity not yet implemented")
}

func (s *UserServer) FindUserByIdentity(ctx context.Context, req *gen.FindUserByIdentityRequest) (*gen.User, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "FindUserByIdentity not yet implemented")
}

func (s *UserServer) ListUserIdentities(ctx context.Context, req *gen.ListUserIdentitiesRequest) (*gen.ListUserIdentitiesResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "ListUserIdentities not yet implemented")
}

// ============================================================================
// OrganizationService RPCs (on OrgServer)
// ============================================================================

func (s *OrgServer) CreateOrganization(ctx context.Context, req *gen.CreateOrganizationRequest) (*gen.CreateOrganizationResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("CreateOrganization")
	w.GRPC().Inject()
	userID, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	return service.CreateOrganization(ctx, userID, req)
}

func (s *OrgServer) GetOrganization(ctx context.Context, req *gen.GetOrganizationRequest) (*gen.Organization, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "GetOrganization not yet implemented")
}

func (s *OrgServer) ListOrganizations(ctx context.Context, _ *gen.ListOrganizationsRequest) (*gen.ListOrganizationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListOrganizations not yet implemented")
}

func (s *OrgServer) AddMember(ctx context.Context, req *gen.AddOrgMemberRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "AddMember not yet implemented")
}

func (s *OrgServer) RemoveMember(ctx context.Context, req *gen.RemoveOrgMemberRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "RemoveMember not yet implemented")
}

func (s *OrgServer) ListMembers(ctx context.Context, req *gen.ListOrgMembersRequest) (*gen.ListOrgMembersResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "ListMembers not yet implemented")
}

// ============================================================================
// TeamService RPCs (on TeamServer)
// ============================================================================

func (s *TeamServer) CreateTeam(ctx context.Context, req *gen.CreateTeamRequest) (*gen.CreateTeamResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.CreateTeam(ctx, req)
}

func (s *TeamServer) ListTeams(ctx context.Context, req *gen.ListTeamsRequest) (*gen.ListTeamsResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "ListTeams not yet implemented")
}

func (s *TeamServer) AddMember(ctx context.Context, req *gen.AddTeamMemberRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "AddMember not yet implemented")
}

func (s *TeamServer) RemoveMember(ctx context.Context, req *gen.RemoveTeamMemberRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "RemoveMember not yet implemented")
}

func (s *TeamServer) ListMembers(ctx context.Context, req *gen.ListTeamMembersRequest) (*gen.ListTeamMembersResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "ListMembers not yet implemented")
}

// ============================================================================
// PermissionService RPCs (on PermServer)
// ============================================================================

func (s *PermServer) CreateRole(ctx context.Context, req *gen.CreateRoleRequest) (*gen.CreateRoleResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.CreateRole(ctx, req)
}

func (s *PermServer) ListRoles(ctx context.Context, _ *gen.ListRolesRequest) (*gen.ListRolesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListRoles not yet implemented")
}

func (s *PermServer) DeleteRole(ctx context.Context, req *gen.DeleteRoleRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "DeleteRole not yet implemented")
}

func (s *PermServer) AssignRole(ctx context.Context, req *gen.AssignRoleRequest) (*gen.AssignRoleResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.AssignRole(ctx, req)
}

func (s *PermServer) RevokeRole(ctx context.Context, req *gen.RevokeRoleRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "RevokeRole not yet implemented")
}

func (s *PermServer) CheckPermission(ctx context.Context, req *gen.CheckPermissionRequest) (*gen.CheckPermissionResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.CheckPermission(ctx, req)
}

// ============================================================================
// IdentityService RPCs (on IdentServer)
// ============================================================================

func (s *IdentServer) ResolveIdentity(ctx context.Context, req *gen.ResolveIdentityRequest) (*gen.ResolveIdentityResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.ResolveIdentity(ctx, req)
}

// ============================================================================
// APIKeyService RPCs (on APIKeyServer)
// ============================================================================

func (s *APIKeyServer) CreateAPIKey(ctx context.Context, req *gen.CreateAPIKeyRequest) (*gen.CreateAPIKeyResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("CreateAPIKey")
	w.GRPC().Inject()
	userID, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	return service.CreateAPIKey(ctx, userID, req)
}

func (s *APIKeyServer) ListAPIKeys(ctx context.Context, req *gen.ListAPIKeysRequest) (*gen.ListAPIKeysResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.ListAPIKeys(ctx, req)
}

func (s *APIKeyServer) RevokeAPIKey(ctx context.Context, req *gen.RevokeAPIKeyRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	if err := service.RevokeAPIKey(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *APIKeyServer) ValidateAPIKey(ctx context.Context, req *gen.ValidateAPIKeyRequest) (*gen.ValidateAPIKeyResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.ValidateAPIKey(ctx, req.KeyHash)
}

// ============================================================================
// AuthService RPCs (on AuthServer)
// ============================================================================

func (s *AuthServer) Authenticate(ctx context.Context, req *gen.AuthenticateRequest) (*gen.AuthenticateResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.Authenticate(ctx, req)
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.RefreshToken(ctx, req)
}

func (s *AuthServer) Logout(ctx context.Context, req *gen.LogoutRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	if err := service.Logout(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthServer) GetJWKS(ctx context.Context, _ *emptypb.Empty) (*gen.JWKSResponse, error) {
	jwks, err := service.GetJWKS(ctx)
	if err != nil {
		return nil, err
	}
	return &gen.JWKSResponse{KeysJson: jwks}, nil
}

// ============================================================================
// AuditService RPCs (on AuditServer)
// ============================================================================

func (s *AuditServer) QueryAuditLog(ctx context.Context, req *gen.QueryAuditLogRequest) (*gen.QueryAuditLogResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}

	var from, to *time.Time
	if req.From != nil {
		t := req.From.AsTime()
		from = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		to = &t
	}

	entries, nextToken, totalCount, err := service.QueryAuditLog(ctx,
		req.OrgId, req.ActorId, req.Action, req.Resource, req.ResourceId,
		from, to, req.PageSize, req.PageToken)
	if err != nil {
		return nil, err
	}

	var events []*gen.AuditEvent
	for _, e := range entries {
		events = append(events, infra.AuditEntryToProto(e))
	}

	return &gen.QueryAuditLogResponse{
		Events:        events,
		NextPageToken: nextToken,
		TotalCount:    totalCount,
	}, nil
}

// ============================================================================
// InvitationService RPCs (on InvitationServer)
// ============================================================================

func (s *InvitationServer) CreateInvitation(ctx context.Context, req *gen.CreateInvitationRequest) (*gen.CreateInvitationResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("CreateInvitation")
	w.GRPC().Inject()
	userID, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	return service.CreateInvitation(ctx, userID, req)
}

func (s *InvitationServer) AcceptInvitation(ctx context.Context, req *gen.AcceptInvitationRequest) (*gen.AcceptInvitationResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("AcceptInvitation")
	w.GRPC().Inject()
	userID, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	return service.AcceptInvitation(ctx, userID, req)
}

func (s *InvitationServer) ListInvitations(ctx context.Context, req *gen.ListInvitationsRequest) (*gen.ListInvitationsResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return service.ListInvitations(ctx, req)
}

func (s *InvitationServer) RevokeInvitation(ctx context.Context, req *gen.RevokeInvitationRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	w := wool.Get(ctx).In("RevokeInvitation")
	w.GRPC().Inject()
	userID, found := w.UserAuthID()
	if !found {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}
	if err := service.RevokeInvitation(ctx, userID, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ============================================================================
// AdminService RPCs (on AdminServer — requires admin role)
// ============================================================================

func (s *AdminServer) SearchUsers(ctx context.Context, req *gen.SearchUsersRequest) (*gen.SearchUsersResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "SearchUsers not yet implemented")
}

func (s *AdminServer) SuspendUser(ctx context.Context, req *gen.SuspendUserRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "SuspendUser not yet implemented")
}

func (s *AdminServer) UnsuspendUser(ctx context.Context, req *gen.UnsuspendUserRequest) (*emptypb.Empty, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "UnsuspendUser not yet implemented")
}

func (s *AdminServer) ImpersonateUser(ctx context.Context, req *gen.ImpersonateUserRequest) (*gen.ImpersonateUserResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "ImpersonateUser not yet implemented")
}

func (s *AdminServer) ListActiveSessions(ctx context.Context, req *gen.ListActiveSessionsRequest) (*gen.ListActiveSessionsResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "ListActiveSessions not yet implemented")
}

func (s *AdminServer) GetOrgEntitlements(ctx context.Context, req *gen.GetOrgEntitlementsRequest) (*gen.GetOrgEntitlementsResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "GetOrgEntitlements not yet implemented")
}

func (s *AdminServer) OverrideEntitlement(ctx context.Context, req *gen.OverrideEntitlementRequest) (*gen.OverrideEntitlementResponse, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "OverrideEntitlement not yet implemented")
}
