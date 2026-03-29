import { useQuery } from "@tanstack/react-query";
import * as api from "../api-client";

export function useActiveSessions(userId?: string) {
  return useQuery({
    queryKey: ["sessions", userId],
    queryFn: () => api.listActiveSessions(userId),
    select: (data) => data.sessions ?? [],
  });
}
