import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as api from "../api-client";

export function useAPIKeys(orgId: string | null) {
  return useQuery({
    queryKey: ["api-keys", orgId],
    queryFn: () => api.listAPIKeys(orgId!),
    enabled: !!orgId,
    select: (data) => data.keys ?? [],
  });
}

export function useRevokeAPIKey() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.revokeAPIKey(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["api-keys"] }),
  });
}
