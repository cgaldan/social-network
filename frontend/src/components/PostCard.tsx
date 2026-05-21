import Link from "next/link";
import { resolveMediaUrl } from "@/lib/api";
import { formatRelative } from "@/lib/format";
import type { Post } from "@/types/api";

const PRIVACY_LABEL: Record<string, string> = {
  public: "Public",
  almost_private: "Followers",
  private: "Private",
};

export function PostCard({ post }: { post: Post }) {
  return (
    <article className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex items-center justify-between text-sm text-slate-500">
        <div className="flex items-center gap-2">
          <Link
            href={`/users/${post.user_id}`}
            className="font-medium text-slate-700 hover:text-indigo-600 hover:underline"
          >
            @{post.author}
          </Link>
          <span aria-hidden>•</span>
          <span>{formatRelative(post.created_at)}</span>
          {post.category ? (
            <>
              <span aria-hidden>•</span>
              <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs">{post.category}</span>
            </>
          ) : null}
        </div>
        <span className="rounded-full bg-slate-50 px-2 py-0.5 text-xs text-slate-500">
          {PRIVACY_LABEL[post.privacy_level] ?? post.privacy_level}
        </span>
      </div>

      <h2 className="mt-3 text-lg font-semibold text-slate-900">
        <Link href={`/posts/${post.id}`} className="hover:underline">
          {post.title}
        </Link>
      </h2>
      <p className="mt-1 whitespace-pre-wrap text-sm text-slate-700">{post.content}</p>

      {post.media_url ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={resolveMediaUrl(post.media_url)}
          alt=""
          className="mt-3 max-h-96 w-full rounded-lg object-cover"
        />
      ) : null}

      <div className="mt-4 flex items-center gap-4 text-sm text-slate-500">
        <Link href={`/posts/${post.id}`} className="hover:text-indigo-600">
          💬 {post.comment_count} comments
        </Link>
      </div>
    </article>
  );
}
