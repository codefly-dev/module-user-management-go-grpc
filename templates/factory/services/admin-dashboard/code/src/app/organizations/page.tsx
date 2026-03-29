"use client";

import { useOrganizations } from "@/lib/hooks";
import { formatDate } from "@/lib/utils";

export default function OrganizationsPage() {
  const { data: orgs = [], isLoading } = useOrganizations();

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Organizations</h2>
      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Name</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Slug</th>
              <th className="text-left px-6 py-3 font-medium text-gray-500">Created</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={3} className="px-6 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : orgs.length === 0 ? (
              <tr><td colSpan={3} className="px-6 py-8 text-center text-gray-400">No organizations</td></tr>
            ) : (
              orgs.map((org) => (
                <tr key={org.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium">{org.name}</td>
                  <td className="px-6 py-4 text-gray-500">{org.slug}</td>
                  <td className="px-6 py-4 text-gray-500">{formatDate(org.createdAt)}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
