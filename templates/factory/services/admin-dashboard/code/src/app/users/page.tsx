"use client";

import { useState } from "react";
import { useUsers, useSuspendUser, useUnsuspendUser } from "@/lib/hooks";
import { formatDate } from "@/lib/utils";

export default function UsersPage() {
  const [query, setQuery] = useState("");
  const { data: users = [], isLoading } = useUsers(query);
  const suspendUser = useSuspendUser();
  const unsuspendUser = useUnsuspendUser();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Users</h2>
        <input
          type="text"
          placeholder="Search by email or name..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="px-4 py-2 border border-gray-300 rounded-lg w-80 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Email</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Status</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Created</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : users.length === 0 ? (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">No users found</td></tr>
            ) : (
              users.map((user) => (
                <tr key={user.uuid} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium">{user.primaryEmail}</td>
                  <td className="px-6 py-4">
                    <StatusBadge status={user.status} />
                  </td>
                  <td className="px-6 py-4 text-gray-500">{formatDate(user.createdAt)}</td>
                  <td className="px-6 py-4 space-x-2">
                    {user.status === "USER_STATUS_ACTIVE" && (
                      <button
                        onClick={() => suspendUser.mutate({ userId: user.uuid!, reason: "admin action" })}
                        className="text-red-600 hover:text-red-800 text-sm font-medium"
                      >
                        Suspend
                      </button>
                    )}
                    {user.status === "USER_STATUS_SUSPENDED" && (
                      <button
                        onClick={() => unsuspendUser.mutate(user.uuid!)}
                        className="text-green-600 hover:text-green-800 text-sm font-medium"
                      >
                        Unsuspend
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

function StatusBadge({ status }: { status?: string }) {
  const label = status?.replace("USER_STATUS_", "").toLowerCase() ?? "unknown";
  const color = status === "USER_STATUS_ACTIVE" ? "bg-green-100 text-green-800" :
    status === "USER_STATUS_SUSPENDED" ? "bg-red-100 text-red-800" :
    "bg-gray-100 text-gray-800";
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${color}`}>
      {label}
    </span>
  );
}
