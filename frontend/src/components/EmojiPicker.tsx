"use client";

import { useEffect, useRef, useState } from "react";

const CATEGORIES: Record<string, string[]> = {
  Smileys: [
    "😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🙂", "🙃",
    "😉", "😊", "😇", "🥰", "😍", "🤩", "😘", "😗", "☺️", "😚",
    "😙", "🥲", "😋", "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭",
    "🤫", "🤔", "🤐", "🤨", "😐", "😑", "😶", "😏", "😒", "🙄",
    "😬", "😮‍💨", "🤥", "😌", "😔", "😪", "🤤", "😴", "😷", "🤒",
    "🤕", "🤢", "🤮", "🤧", "🥵", "🥶", "🥴", "😵", "🤯", "🤠",
  ],
  Hearts: ["❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍", "🤎", "💔", "❣️", "💕", "💞", "💓", "💗", "💖", "💘", "💝"],
  Gestures: ["👍", "👎", "👌", "🤌", "🤏", "✌️", "🤞", "🫰", "🤟", "🤘", "🤙", "👈", "👉", "👆", "🖕", "👇", "☝️", "👋", "🤚", "🖐️", "✋", "🖖", "👏", "🙌", "🙏"],
  Animals: ["🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼", "🐨", "🐯", "🦁", "🐮", "🐷", "🐸", "🐵", "🐔", "🐧", "🐦", "🐤", "🦄", "🦋", "🐢", "🐙", "🐠"],
  Food: ["🍎", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇", "🍓", "🫐", "🍒", "🥭", "🍍", "🥥", "🥝", "🍅", "🍆", "🥑", "🥦", "🌽", "🌶️", "🥕", "🥐", "🥖", "🧀", "🍔", "🍕", "🌮", "🍣", "🍩", "🍪"],
  Activities: ["⚽", "🏀", "🏈", "⚾", "🎾", "🏐", "🏉", "🎱", "🏓", "🏸", "🎮", "🎲", "🎯", "🎤", "🎧", "🎵", "🎶", "🎸", "🎹", "🎺"],
  Objects: ["💻", "📱", "⌨️", "🖥️", "🖱️", "💾", "💿", "📷", "📹", "📺", "📞", "☎️", "🔋", "🔌", "💡", "🔦", "🕯️", "💸", "💰", "🎁", "📦", "📚", "📝", "✏️", "🔑", "🔒"],
  Symbols: ["✅", "❌", "⭐", "🌟", "✨", "⚡", "🔥", "💯", "❓", "❗", "‼️", "⁉️", "🆗", "🆕", "🆒", "🆓", "🔔", "🔕", "🎉", "🎊", "🎈", "🎂", "💬", "💭"],
};

const CATEGORY_NAMES = Object.keys(CATEGORIES);

interface Props {
  onPick: (emoji: string) => void;
  align?: "left" | "right";
}

export function EmojiPicker({ onPick, align = "right" }: Props) {
  const [open, setOpen] = useState(false);
  const [category, setCategory] = useState(CATEGORY_NAMES[0]);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const handler = (e: MouseEvent) => {
      if (!containerRef.current?.contains(e.target as Node)) setOpen(false);
    };
    const onEsc = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", handler);
    document.addEventListener("keydown", onEsc);
    return () => {
      document.removeEventListener("mousedown", handler);
      document.removeEventListener("keydown", onEsc);
    };
  }, [open]);

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label="Insert emoji"
        className="rounded-lg border border-slate-300 bg-white px-2 py-2 text-base hover:bg-slate-50"
      >
        😊
      </button>
      {open ? (
        <div
          className={`absolute bottom-full z-20 mb-2 w-72 rounded-xl border border-slate-200 bg-white shadow-xl ${
            align === "right" ? "right-0" : "left-0"
          }`}
        >
          <div className="flex gap-1 overflow-x-auto border-b border-slate-200 p-2">
            {CATEGORY_NAMES.map((name) => (
              <button
                key={name}
                type="button"
                onClick={() => setCategory(name)}
                className={`shrink-0 rounded-md px-2 py-1 text-xs font-medium transition ${
                  category === name
                    ? "bg-indigo-50 text-indigo-700"
                    : "text-slate-600 hover:bg-slate-100"
                }`}
              >
                {name}
              </button>
            ))}
          </div>
          <div className="grid max-h-56 grid-cols-8 gap-1 overflow-y-auto p-2 text-xl">
            {CATEGORIES[category].map((emoji, i) => (
              <button
                key={`${category}-${i}`}
                type="button"
                onClick={() => {
                  onPick(emoji);
                }}
                className="rounded-md p-1 hover:bg-slate-100"
                aria-label={`Insert ${emoji}`}
              >
                {emoji}
              </button>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}
