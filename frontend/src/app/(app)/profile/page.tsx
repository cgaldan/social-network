"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/AuthProvider";

export default function ProfileRedirectPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (loading) return;
    if (user) router.replace(`/users/${user.id}`);
    else router.replace("/login");
  }, [user, loading, router]);

  return (
    <div className="flex min-h-[50vh] items-center justify-center text-sm text-slate-500">
      Loading profile…
    </div>
  );
}
