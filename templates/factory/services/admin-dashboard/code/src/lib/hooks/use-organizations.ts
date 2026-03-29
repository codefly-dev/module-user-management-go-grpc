import { useQuery } from "@tanstack/react-query";
import * as api from "../api-client";

export function useOrganizations() {
  return useQuery({
    queryKey: ["organizations"],
    queryFn: () => api.listOrganizations(),
    select: (data) => data.organizations ?? [],
  });
}

export function useOrgEntitlements(orgId: string | null) {
  return useQuery({
    queryKey: ["entitlements", orgId],
    queryFn: () => api.getOrgEntitlements(orgId!),
    enabled: !!orgId,
  });
}
