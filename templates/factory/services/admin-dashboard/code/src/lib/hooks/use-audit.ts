import { useQuery } from "@tanstack/react-query";
import * as api from "../api-client";
import type { AuditLogParams } from "../api-client";

export function useAuditLog(params: AuditLogParams) {
  return useQuery({
    queryKey: ["audit-log", params],
    queryFn: () => api.queryAuditLog(params),
    select: (data) => ({
      events: data.events ?? [],
      totalCount: data.total_count ?? 0,
    }),
  });
}
