import type { Metadata } from "next";
import Link from "next/link";
import { Providers } from "@/lib/providers";
import "./globals.css";

export const metadata: Metadata = {
  title: "Admin Dashboard",
  description: "User management administration",
};

const navigation = [
  { name: "Users", href: "/users" },
  { name: "Organizations", href: "/organizations" },
  { name: "Invitations", href: "/invitations" },
  { name: "API Keys", href: "/api-keys" },
  { name: "Sessions", href: "/sessions" },
  { name: "Audit Log", href: "/audit-log" },
  { name: "Entitlements", href: "/entitlements" },
];

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-gray-50 text-gray-900 antialiased">
        <Providers>
          <div className="flex h-screen">
            <aside className="w-56 bg-white border-r border-gray-200 flex flex-col">
              <div className="p-5 border-b border-gray-200">
                <h1 className="text-lg font-semibold">Admin</h1>
                <p className="text-xs text-gray-500">User Management</p>
              </div>
              <nav className="flex-1 p-3 space-y-0.5">
                {navigation.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    className="block px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100 transition-colors"
                  >
                    {item.name}
                  </Link>
                ))}
              </nav>
            </aside>
            <main className="flex-1 overflow-auto">
              <div className="p-8">{children}</div>
            </main>
          </div>
        </Providers>
      </body>
    </html>
  );
}
