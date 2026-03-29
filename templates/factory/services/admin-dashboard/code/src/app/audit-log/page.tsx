"use client";

import { useState } from "react";
import { useAuditLog } from "@/lib/hooks";
import { formatDate } from "@/lib/utils";

const ACTION_TYPES = [
  "user.registered", "auth.login", "auth.logout",
  "api_key.created", "api_key.revoked", "role.assigned",
  "invitation.created", "invitation.accepted", "org.created",
];

export default function AuditLogPage() {
  const [actionFilter, setActionFilter] = useState("");
  const { data, isLoading } = useAuditLog({
    action: actionFilter || undefined,
    page_size: 100,
  });
  const events = data?.events ?? [];

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Audit Log</h2>
        <select
          value={actionFilter}
          onChange={(e) => setActionFilter(e.target.value)}
          className="px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">All actions</option>
          {ACTION_TYPES.map((action) => (
            <option key={action} value={action}>{action}</option>
          ))}
        </select>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Time</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Action</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Actor</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Resource</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : events.length === 0 ? (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">No events</td></tr>
            ) : (
              events.map((event) => (
                <tr key={event.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-gray-500 whitespace-nowrap">{formatDate(event.createdAt)}</td>
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      {event.action}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-gray-500 font-mono text-xs">{event.actorId?.slice(0, 8)}</td>
                  <td className="px-6 py-4 text-gray-500">{event.resource}/{event.resourceId?.slice(0, 8)}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
