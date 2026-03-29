import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as api from "../api-client";

export function useInvitations(orgId: string | null) {
  return useQuery({
    queryKey: ["invitations", orgId],
    queryFn: () => api.listInvitations(orgId!),
    enabled: !!orgId,
    select: (data) => data.invitations ?? [],
  });
}

export function useCreateInvitation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ orgId, email, role }: { orgId: string; email: string; role: string }) =>
      api.createInvitation(orgId, email, role),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["invitations"] }),
  });
}

export function useRevokeInvitation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.revokeInvitation(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["invitations"] }),
  });
}
