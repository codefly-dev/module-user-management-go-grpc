"use client";

import { useState } from "react";
import { useInvitations, useRevokeInvitation, useOrganizations } from "@/lib/hooks";
import { formatDate } from "@/lib/utils";

export default function InvitationsPage() {
  const { data: orgs = [] } = useOrganizations();
  const [selectedOrg, setSelectedOrg] = useState<string>("");
  const { data: invitations = [], isLoading } = useInvitations(selectedOrg || null);
  const revokeInvitation = useRevokeInvitation();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Invitations</h2>
        <select
          value={selectedOrg}
          onChange={(e) => setSelectedOrg(e.target.value)}
          className="px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select organization...</option>
          {orgs.map((org) => (
            <option key={org.id} value={org.id}>{org.name}</option>
          ))}
        </select>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Email</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Role</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Status</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Expires</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {!selectedOrg ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">Select an organization</td></tr>
            ) : isLoading ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : invitations.length === 0 ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">No invitations</td></tr>
            ) : (
              invitations.map((inv) => (
                <tr key={inv.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium">{inv.email}</td>
                  <td className="px-6 py-4">{inv.role}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      inv.status === "INVITATION_STATUS_PENDING" ? "bg-yellow-100 text-yellow-800" :
                      inv.status === "INVITATION_STATUS_ACCEPTED" ? "bg-green-100 text-green-800" :
                      "bg-gray-100 text-gray-800"
                    }`}>
                      {inv.status?.replace("INVITATION_STATUS_", "").toLowerCase()}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-gray-500">{formatDate(inv.expiresAt)}</td>
                  <td className="px-6 py-4">
                    {inv.status === "INVITATION_STATUS_PENDING" && (
                      <button
                        onClick={() => revokeInvitation.mutate(inv.id!)}
                        className="text-red-600 hover:text-red-800 text-sm font-medium"
                      >
                        Revoke
                      </button>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
