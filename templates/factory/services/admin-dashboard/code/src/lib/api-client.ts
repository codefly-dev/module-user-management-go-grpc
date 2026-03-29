/**
 * Typed API client for the backend gRPC-gateway REST endpoints.
 * All types are auto-generated from proto definitions via the codefly companion.
 *
 * This module is fully testable — no React, no DOM, no side effects.
 */

import type { components, operations } from "./api-types";

// Re-export commonly used types from generated schema
export type User = components["schemas"]["customersUser"];
export type Organization = components["schemas"]["customersOrganization"];
export type APIKey = components["schemas"]["customersAPIKey"];
export type AuditEvent = components["schemas"]["customersAuditEvent"];
export type Invitation = components["schemas"]["customersInvitation"];
export type SessionInfo = components["schemas"]["customersSessionInfo"];
export type EntitlementInfo = components["schemas"]["customersEntitlementInfo"];
export type Permission = components["schemas"]["customersPermission"];

// Auth types
export type AuthenticateResponse = components["schemas"]["customersAuthenticateResponse"];
export type RefreshTokenResponse = components["schemas"]["customersRefreshTokenResponse"];

// Config
export interface APIConfig {
  baseUrl: string;
  getToken?: () => string | null;
}

let config: APIConfig = {
  baseUrl: process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:39042",
};

export function configureAPI(c: Partial<APIConfig>) {
  config = { ...config, ...c };
}

// Core fetch wrapper — testable, no React dependency
export class APIError extends Error {
  constructor(
    public status: number,
    public body: string,
  ) {
    super(`API ${status}: ${body}`);
    this.name = "APIError";
  }
}

export async function request<T>(
  path: string,
  options?: RequestInit,
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  const token = config.getToken?.();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${config.baseUrl}${path}`, {
    ...options,
    headers: { ...headers, ...options?.headers },
  });

  if (!res.ok) {
    const body = await res.text();
    throw new APIError(res.status, body);
  }

  const text = await res.text();
  return text ? JSON.parse(text) : ({} as T);
}

// ============================================================================
// Admin API
// ============================================================================

export function searchUsers(query: string, pageSize = 50) {
  const params = new URLSearchParams({ page_size: String(pageSize) });
  if (query) params.set("query", query);
  return request<{
    users?: User[];
    next_page_token?: string;
    total_count?: number;
  }>(`/v1/admin/users?${params}`);
}

export function suspendUser(userId: string, reason: string) {
  return request<Record<string, never>>(
    `/v1/admin/users/${userId}:suspend`,
    { method: "POST", body: JSON.stringify({ reason }) },
  );
}

export function unsuspendUser(userId: string) {
  return request<Record<string, never>>(
    `/v1/admin/users/${userId}:unsuspend`,
    { method: "POST", body: JSON.stringify({}) },
  );
}

export function listActiveSessions(userId?: string, pageSize = 100) {
  const params = new URLSearchParams({ page_size: String(pageSize) });
  if (userId) params.set("user_id", userId);
  return request<{ sessions?: SessionInfo[]; next_page_token?: string }>(
    `/v1/admin/sessions?${params}`,
  );
}

export function getOrgEntitlements(orgId: string) {
  return request<{
    planName?: string;
    entitlements?: EntitlementInfo[];
  }>(`/v1/admin/organizations/${orgId}/entitlements`);
}

export function overrideEntitlement(
  orgId: string,
  feature: string,
  limitValue: number,
  reason: string,
) {
  return request<{ id?: string }>(
    `/v1/admin/organizations/${orgId}/entitlements`,
    {
      method: "POST",
      body: JSON.stringify({
        feature,
        limit_value: limitValue,
        reason,
      }),
    },
  );
}

// ============================================================================
// Auth API
// ============================================================================

export function authenticate(
  provider: string,
  providerId: string,
  email: string,
) {
  return request<AuthenticateResponse>("/v1/auth/authenticate", {
    method: "POST",
    body: JSON.stringify({
      provider,
      provider_id: providerId,
      provider_email: email,
    }),
  });
}

export function refreshToken(token: string) {
  return request<RefreshTokenResponse>("/v1/auth/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token: token }),
  });
}

export function logout(token: string) {
  return request<Record<string, never>>("/v1/auth/logout", {
    method: "POST",
    body: JSON.stringify({ refresh_token: token }),
  });
}

// ============================================================================
// Audit API
// ============================================================================

export interface AuditLogParams {
  org_id?: string;
  action?: string;
  actor_id?: string;
  page_size?: number;
}

export function queryAuditLog(params: AuditLogParams) {
  const qs = new URLSearchParams();
  if (params.org_id) qs.set("org_id", params.org_id);
  if (params.action) qs.set("action", params.action);
  if (params.actor_id) qs.set("actor_id", params.actor_id);
  qs.set("page_size", String(params.page_size || 50));
  return request<{
    events?: AuditEvent[];
    next_page_token?: string;
    total_count?: number;
  }>(`/v1/audit-log?${qs}`);
}

// ============================================================================
// Invitations API
// ============================================================================

export function listInvitations(orgId: string) {
  return request<{ invitations?: Invitation[] }>(
    `/v1/invitations?org_id=${orgId}`,
  );
}

export function createInvitation(orgId: string, email: string, role: string) {
  return request<{ invitation?: Invitation; invite_token?: string }>(
    "/v1/invitations",
    {
      method: "POST",
      body: JSON.stringify({ org_id: orgId, email, role }),
    },
  );
}

export function revokeInvitation(id: string) {
  return request<Record<string, never>>(`/v1/invitations/${id}`, {
    method: "DELETE",
  });
}

// ============================================================================
// API Keys API
// ============================================================================

export function listAPIKeys(orgId: string, pageSize = 100) {
  return request<{ keys?: APIKey[]; next_page_token?: string }>(
    `/v1/api-keys?organization_id=${orgId}&page_size=${pageSize}`,
  );
}

export function revokeAPIKey(id: string) {
  return request<Record<string, never>>(`/v1/api-keys/${id}`, {
    method: "DELETE",
  });
}

// ============================================================================
// Organizations API
// ============================================================================

export function listOrganizations() {
  return request<{ organizations?: Organization[] }>("/v1/organizations");
}
