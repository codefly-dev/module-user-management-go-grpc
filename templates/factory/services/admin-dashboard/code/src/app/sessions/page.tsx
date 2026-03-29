"use client";

import { useActiveSessions } from "@/lib/hooks";
import { formatDate } from "@/lib/utils";

export default function SessionsPage() {
  const { data: sessions = [], isLoading } = useActiveSessions();

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Active Sessions</h2>
      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-6 py-3 font-medium text-gray-500">User</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">IP</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Last Active</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Expires</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : sessions.length === 0 ? (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">No active sessions</td></tr>
            ) : (
              sessions.map((s) => (
                <tr key={s.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-mono text-xs">{s.userId?.slice(0, 8)}...</td>
                  <td className="px-6 py-4">{s.ipAddress || "-"}</td>
                  <td className="px-6 py-4 text-gray-500">{formatDate(s.lastActiveAt)}</td>
                  <td className="px-6 py-4 text-gray-500">{formatDate(s.expiresAt)}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
