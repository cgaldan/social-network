import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Social Network",
  description: "A Next.js frontend for the Go social network backend.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
