"use client";

import { useCallback, useEffect, useRef, useState } from "react";

export type PageFetcherArgs = {
  limit: number;
  offset: number;
  signal: AbortSignal;
};

export type PageFetcher<T> = (args: PageFetcherArgs) => Promise<{
  items: T[];
  hasMore: boolean;
}>;

export interface UsePaginatedListResult<T> {
  items: T[];
  hasMore: boolean;
  loading: boolean;
  error: string | null;
  loadMore: () => void;
  reset: () => void;
  setItems: React.Dispatch<React.SetStateAction<T[]>>;
}

export function usePaginatedList<T>(
  fetcher: PageFetcher<T>,
  pageSize: number = 20,
  getKey?: (item: T) => string | number,
): UsePaginatedListResult<T> {
  const [items, setItems] = useState<T[]>([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const offsetRef = useRef(0);
  const loadingRef = useRef(false);
  const hasMoreRef = useRef(true);
  const abortRef = useRef<AbortController | null>(null);
  const getKeyRef = useRef(getKey);
  getKeyRef.current = getKey;

  loadingRef.current = loading;
  hasMoreRef.current = hasMore;

  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMoreRef.current) return;

    abortRef.current?.abort();
    const ctl = new AbortController();
    abortRef.current = ctl;

    loadingRef.current = true;
    setLoading(true);
    setError(null);

    try {
      const offset = offsetRef.current;
      const page = await fetcher({ limit: pageSize, offset, signal: ctl.signal });
      if (ctl.signal.aborted) return;
      offsetRef.current = offset + page.items.length;
      setItems((prev) => {
        const key = getKeyRef.current;
        if (!key) return [...prev, ...page.items];
        const seen = new Set(prev.map(key));
        const novel = page.items.filter((it) => !seen.has(key(it)));
        return novel.length === page.items.length ? [...prev, ...page.items] : [...prev, ...novel];
      });
      setHasMore(page.hasMore);
    } catch (e) {
      if (ctl.signal.aborted) return;
      const name = (e as { name?: string } | null)?.name;
      if (name === "AbortError") return;
      setError(e instanceof Error ? e.message : "Failed to load");
    } finally {
      if (!ctl.signal.aborted) {
        loadingRef.current = false;
        setLoading(false);
      }
    }
  }, [fetcher, pageSize]);

  const reset = useCallback(() => {
    abortRef.current?.abort();
    abortRef.current = null;
    offsetRef.current = 0;
    loadingRef.current = false;
    hasMoreRef.current = true;
    setItems([]);
    setHasMore(true);
    setLoading(false);
    setError(null);
  }, []);

  useEffect(() => {
    reset();
    loadMore();
    return () => abortRef.current?.abort();
  }, [fetcher, reset, loadMore]);

  return { items, hasMore, loading, error, loadMore, reset, setItems };
}
