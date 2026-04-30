"use client";

import { AuthProvider } from "@/components/AuthProvider";

export default function AuthPagesLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AuthProvider>{children}</AuthProvider>;
}
