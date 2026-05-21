"use client";

import { useTheme } from "@/components/ThemeProvider";

export function ThemeToggle({ className = "" }: { className?: string }) {
  const { theme, toggle } = useTheme();
  const isDark = theme === "dark";
  return (
    <button
      type="button"
      onClick={toggle}
      aria-label={isDark ? "Switch to light theme" : "Switch to dark theme"}
      title={isDark ? "Switch to light theme" : "Switch to dark theme"}
      className={`inline-flex h-8 w-8 items-center justify-center rounded-md border border-slate-200 bg-white text-base text-slate-600 transition hover:bg-slate-100 ${className}`}
    >
      <span aria-hidden>{isDark ? "☀️" : "🌙"}</span>
    </button>
  );
}
