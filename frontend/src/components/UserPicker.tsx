"use client";

import { useEffect, useRef, useState } from "react";
import { Avatar } from "@/components/Avatar";
import { api } from "@/lib/api";
import type { UserSummary } from "@/types/api";

interface Props {
  onSelect: (user: UserSummary) => void;
  placeholder?: string;
  clearOnSelect?: boolean;
}

export function UserPicker({ onSelect, placeholder = "Search users…", clearOnSelect = true }: Props) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<UserSummary[]>([]);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (!containerRef.current?.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  useEffect(() => {
    const trimmed = query.trim();
    if (!open) return;
    if (trimmed.length === 0 && !open) return;
    const ctrl = new AbortController();
    setLoading(true);
    const t = setTimeout(async () => {
      try {
        const res = await api.users.list({ q: trimmed || undefined, limit: 10 });
        setResults(res.users ?? []);
      } catch {
        setResults([]);
      } finally {
        setLoading(false);
      }
    }, 200);
    return () => {
      ctrl.abort();
      clearTimeout(t);
    };
  }, [query, open]);

  const pick = (u: UserSummary) => {
    onSelect(u);
    if (clearOnSelect) setQuery("");
    setOpen(false);
  };

  return (
    <div ref={containerRef} className="relative">
      <input
        type="text"
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        placeholder={placeholder}
        className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
      />

      {open ? (
        <div className="absolute z-10 mt-1 max-h-72 w-full overflow-y-auto rounded-lg border border-slate-200 bg-white shadow-lg">
          {loading ? (
            <p className="px-3 py-2 text-xs text-slate-500">Searching…</p>
          ) : results.length === 0 ? (
            <p className="px-3 py-2 text-xs text-slate-500">No users found.</p>
          ) : (
            <ul>
              {results.map((u) => (
                <li key={u.id}>
                  <button
                    type="button"
                    onClick={() => pick(u)}
                    className="flex w-full items-center gap-3 px-3 py-2 text-left text-sm hover:bg-slate-50"
                  >
                    <Avatar
                      src={u.avatar_path}
                      firstName={u.first_name}
                      lastName={u.last_name}
                      nickname={u.nickname}
                      size={32}
                    />
                    <div className="min-w-0">
                      <div className="truncate font-medium text-slate-900">
                        {u.first_name} {u.last_name}
                      </div>
                      <div className="truncate text-xs text-slate-500">
                        @{u.nickname}
                        {u.is_online ? " · online" : ""}
                      </div>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      ) : null}
    </div>
  );
}
