import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as api from "../api-client";

export function useUsers(query: string) {
  return useQuery({
    queryKey: ["users", query],
    queryFn: () => api.searchUsers(query),
    select: (data) => data.users ?? [],
  });
}

export function useSuspendUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, reason }: { userId: string; reason: string }) =>
      api.suspendUser(userId, reason),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });
}

export function useUnsuspendUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => api.unsuspendUser(userId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });
}
