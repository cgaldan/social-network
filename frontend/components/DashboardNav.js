"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/dashboard", label: "Home" },
  { href: "/dashboard/feed", label: "Feed" },
  { href: "/dashboard/groups", label: "Groups" },
  { href: "/dashboard/messages", label: "Messages" },
  { href: "/dashboard/notifications", label: "Notifications" },
  { href: "/dashboard/social", label: "Follow" },
  { href: "/dashboard/profile", label: "Profile" },
  { href: "/dashboard/realtime", label: "Live" },
];

export default function DashboardNav({ unreadCount = 0 }) {
  const pathname = usePathname();

  return (
    <header className="app-header">
      <div className="app-header-inner">
        <Link href="/dashboard" className="app-brand">
          Social
        </Link>
        <nav className="app-nav" aria-label="Main">
          {links.map(({ href, label }) => {
            const active =
              href === "/dashboard"
                ? pathname === "/dashboard"
                : pathname.startsWith(href);
            return (
              <Link
                key={href}
                href={href}
                className={active ? "nav-link nav-link-active" : "nav-link"}
              >
                {label === "Notifications" && unreadCount > 0 ? (
                  <>
                    {label}{" "}
                    <span className="nav-badge">{unreadCount > 99 ? "99+" : unreadCount}</span>
                  </>
                ) : (
                  label
                )}
              </Link>
            );
          })}
        </nav>
      </div>
    </header>
  );
}
