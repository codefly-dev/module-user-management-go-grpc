"use client";

import { useState } from "react";
import { useOrganizations, useOrgEntitlements } from "@/lib/hooks";
import { formatLimit } from "@/lib/utils";

export default function EntitlementsPage() {
  const { data: orgs = [] } = useOrganizations();
  const [selectedOrg, setSelectedOrg] = useState("");
  const { data, isLoading } = useOrgEntitlements(selectedOrg || null);
  const entitlements = data?.entitlements ?? [];
  const planName = data?.planName ?? "";

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-4">
          <h2 className="text-2xl font-bold">Entitlements</h2>
          {planName && (
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-purple-100 text-purple-800">
              {planName}
            </span>
          )}
        </div>
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
              <th className="text-left px-6 py-3 font-medium text-gray-500">Feature</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Limit</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Used</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Usage</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Override</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {!selectedOrg ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">Select an organization</td></tr>
            ) : isLoading ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : entitlements.length === 0 ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400">No entitlements</td></tr>
            ) : (
              entitlements.map((ent) => {
                const limit = Number(ent.limit ?? 0);
                const used = Number(ent.used ?? 0);
                const pct = limit > 0 ? Math.min(100, Math.round((used / limit) * 100)) : 0;
                return (
                  <tr key={ent.feature} className="hover:bg-gray-50">
                    <td className="px-6 py-4 font-medium">{ent.feature}</td>
                    <td className="px-6 py-4">{formatLimit(limit)}</td>
                    <td className="px-6 py-4">{used.toLocaleString()}</td>
                    <td className="px-6 py-4 w-48">
                      {limit > 0 && (
                        <div className="flex items-center gap-2">
                          <div className="flex-1 bg-gray-200 rounded-full h-2">
                            <div
                              className={`h-2 rounded-full ${pct > 80 ? "bg-red-500" : pct > 50 ? "bg-yellow-500" : "bg-green-500"}`}
                              style={{ width: `${pct}%` }}
                            />
                          </div>
                          <span className="text-xs text-gray-500">{pct}%</span>
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      {ent.hasOverride && (
                        <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                          Override
                        </span>
                      )}
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
