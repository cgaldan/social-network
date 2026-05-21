import { resolveMediaUrl } from "@/lib/api";
import { initials } from "@/lib/format";

interface Props {
  src?: string;
  firstName?: string;
  lastName?: string;
  nickname?: string;
  size?: number;
  className?: string;
}

export function Avatar({ src, firstName = "", lastName = "", nickname, size = 40, className = "" }: Props) {
  const label = initials(firstName, lastName) || (nickname?.[0]?.toUpperCase() ?? "?");
  const resolved = resolveMediaUrl(src);
  if (resolved) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={resolved}
        alt={nickname ?? `${firstName} ${lastName}`}
        width={size}
        height={size}
        style={{ width: size, height: size }}
        className={`rounded-full object-cover ring-1 ring-slate-200 ${className}`}
      />
    );
  }
  return (
    <div
      style={{ width: size, height: size, fontSize: size * 0.4 }}
      className={`flex items-center justify-center rounded-full bg-indigo-100 font-medium text-indigo-700 ${className}`}
    >
      {label}
    </div>
  );
}
