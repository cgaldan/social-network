"use client";

import { useEffect, useRef } from "react";

interface Props {
  onIntersect: () => void;
  enabled?: boolean;
  rootMargin?: string;
  className?: string;
}


export function InfiniteScrollSentinel({
  onIntersect,
  enabled = true,
  rootMargin = "300px",
  className,
}: Props) {
  const ref = useRef<HTMLDivElement | null>(null);
  const cbRef = useRef(onIntersect);
  cbRef.current = onIntersect;

  useEffect(() => {
    if (!enabled || !ref.current) return;
    const obs = new IntersectionObserver(
      (entries) => {
        if (entries.some((e) => e.isIntersecting)) cbRef.current();
      },
      { rootMargin },
    );
    obs.observe(ref.current);
    return () => obs.disconnect();
  }, [enabled, rootMargin]);

  return <div ref={ref} aria-hidden className={className} style={{ height: 1 }} />;
}
