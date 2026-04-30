"use client";

import { useRouter } from "next/navigation";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { api } from "@/lib/api";
import {
  clearSession,
  getStoredUser,
  getToken,
  saveSession,
  updateStoredUser,
} from "@/lib/auth";
import type { User } from "@/types/api";

type AuthContextValue = {
  user: User | null;
  token: string | null;
  loading: boolean;
  setSession: (token: string, user: User) => void;
  refreshUser: () => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const setSession = useCallback((nextToken: string, nextUser: User) => {
    saveSession(nextToken, nextUser);
    setToken(nextToken);
    setUser(nextUser);
  }, []);

  const refreshUser = useCallback(async () => {
    const response = await api.me();
    if (response.user) {
      updateStoredUser(response.user);
      setUser(response.user);
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      if (getToken()) {
        await api.logout();
      }
    } finally {
      clearSession();
      setUser(null);
      setToken(null);
      router.push("/login");
    }
  }, [router]);

  useEffect(() => {
    const storedToken = getToken();
    setToken(storedToken);
    setUser(getStoredUser());

    if (!storedToken) {
      setLoading(false);
      return;
    }

    api
      .me()
      .then((response) => {
        if (response.user) {
          updateStoredUser(response.user);
          setUser(response.user);
        }
      })
      .catch(() => {
        clearSession();
        setToken(null);
        setUser(null);
      })
      .finally(() => setLoading(false));
  }, []);

  const value = useMemo(
    () => ({ user, token, loading, setSession, refreshUser, logout }),
    [user, token, loading, setSession, refreshUser, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
}
