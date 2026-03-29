import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { request, APIError, configureAPI } from "../api-client";

// Mock global fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

beforeEach(() => {
  configureAPI({ baseUrl: "http://test:8080" });
  mockFetch.mockReset();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("request", () => {
  it("makes a GET request with correct URL", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => JSON.stringify({ users: [] }),
    });

    const result = await request<{ users: unknown[] }>("/v1/users");
    expect(result).toEqual({ users: [] });
    expect(mockFetch).toHaveBeenCalledWith(
      "http://test:8080/v1/users",
      expect.objectContaining({
        headers: expect.objectContaining({
          "Content-Type": "application/json",
        }),
      }),
    );
  });

  it("includes auth token when configured", async () => {
    configureAPI({
      baseUrl: "http://test:8080",
      getToken: () => "test-token-123",
    });

    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => "{}",
    });

    await request("/v1/users");
    expect(mockFetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: "Bearer test-token-123",
        }),
      }),
    );
  });

  it("does not include auth header when no token", async () => {
    configureAPI({ baseUrl: "http://test:8080", getToken: () => null });

    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => "{}",
    });

    await request("/v1/users");
    const headers = mockFetch.mock.calls[0][1].headers;
    expect(headers).not.toHaveProperty("Authorization");
  });

  it("throws APIError on non-ok response", async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 404,
      text: async () => "not found",
    });

    await expect(request("/v1/missing")).rejects.toThrow(APIError);
    await expect(request("/v1/missing")).rejects.toThrow("API 404");
  });

  it("passes POST body correctly", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => JSON.stringify({ id: "123" }),
    });

    await request("/v1/auth/authenticate", {
      method: "POST",
      body: JSON.stringify({ provider: "google" }),
    });

    expect(mockFetch).toHaveBeenCalledWith(
      "http://test:8080/v1/auth/authenticate",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ provider: "google" }),
      }),
    );
  });

  it("handles empty response body", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => "",
    });

    const result = await request("/v1/auth/logout");
    expect(result).toEqual({});
  });
});

describe("APIError", () => {
  it("contains status and body", () => {
    const err = new APIError(403, "forbidden");
    expect(err.status).toBe(403);
    expect(err.body).toBe("forbidden");
    expect(err.message).toBe("API 403: forbidden");
    expect(err.name).toBe("APIError");
  });
});

describe("API functions", () => {
  it("searchUsers builds correct URL with query", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => JSON.stringify({ users: [{ uuid: "1", primary_email: "a@b.com" }] }),
    });

    const { searchUsers } = await import("../api-client");
    const result = await searchUsers("alice", 25);
    expect(result.users).toHaveLength(1);

    const url = mockFetch.mock.calls[0][0] as string;
    expect(url).toContain("query=alice");
    expect(url).toContain("page_size=25");
  });

  it("queryAuditLog builds query string from params", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => JSON.stringify({ events: [], total_count: 0 }),
    });

    const { queryAuditLog } = await import("../api-client");
    await queryAuditLog({ org_id: "org-1", action: "user.registered", page_size: 10 });

    const url = mockFetch.mock.calls[0][0] as string;
    expect(url).toContain("org_id=org-1");
    expect(url).toContain("action=user.registered");
    expect(url).toContain("page_size=10");
  });
});
