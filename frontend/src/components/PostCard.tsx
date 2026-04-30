import Link from "next/link";
import { formatDate } from "@/lib/format";
import type { Post } from "@/types/api";

export function PostCard({ post }: { post: Post }) {
  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex flex-wrap items-center gap-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
        <span>{post.category || "General"}</span>
        <span>/</span>
        <span>{post.privacy_level}</span>
        {post.group_id ? (
          <>
            <span>/</span>
            <span>Group #{post.group_id}</span>
          </>
        ) : null}
      </div>
      <Link href={`/posts/${post.id}`}>
        <h2 className="mt-3 text-2xl font-bold text-slate-950 transition hover:text-sky-700">
          {post.title}
        </h2>
      </Link>
      <p className="mt-3 whitespace-pre-wrap text-slate-700">{post.content}</p>
      {post.media_url ? (
        <a
          className="mt-3 inline-flex text-sm font-semibold text-sky-700"
          href={post.media_url}
          rel="noreferrer"
          target="_blank"
        >
          View media
        </a>
      ) : null}
      <footer className="mt-5 flex flex-wrap gap-4 text-sm text-slate-500">
        <span>By {post.author || `User #${post.user_id}`}</span>
        <span>{formatDate(post.created_at)}</span>
        <span>{post.comment_count} comments</span>
        <span>{post.like_count} likes</span>
      </footer>
    </article>
  );
}
