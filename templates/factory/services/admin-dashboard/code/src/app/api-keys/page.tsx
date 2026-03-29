"use client";

import { useState } from "react";
import { useAPIKeys, useRevokeAPIKey, useOrganizations } from "@/lib/hooks";
import { formatDate } from "@/lib/utils";

export default function APIKeysPage() {
  const { data: orgs = [] } = useOrganizations();
  const [selectedOrg, setSelectedOrg] = useState("");
  const { data: keys = [], isLoading } = useAPIKeys(selectedOrg || null);
  const revokeKey = useRevokeAPIKey();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">API Keys</h2>
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
              <th className="text-left px-6 py-3 font-medium text-gray-500">Name</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Prefix</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Env</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Last Used</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {!selectedOrg ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">Select an organization</td></tr>
            ) : isLoading ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : keys.length === 0 ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">No API keys</td></tr>
            ) : (
              keys.map((key) => (
                <tr key={key.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium">{key.name}</td>
                  <td className="px-6 py-4 font-mono text-xs">{key.prefix}...</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                      key.environment === "API_KEY_ENVIRONMENT_LIVE" ? "bg-green-100 text-green-800" : "bg-orange-100 text-orange-800"
                    }`}>
                      {key.environment?.replace("API_KEY_ENVIRONMENT_", "").toLowerCase()}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-gray-500">{formatDate(key.lastUsedAt)}</td>
                  <td className="px-6 py-4">
                    {!key.revokedAt && (
                      <button
                        onClick={() => revokeKey.mutate(key.id!)}
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
