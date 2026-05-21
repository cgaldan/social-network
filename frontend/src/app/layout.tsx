import type { Metadata } from "next";
import { AuthProvider } from "@/components/AuthProvider";
import { ThemeProvider, themeInitScript } from "@/components/ThemeProvider";
import "./globals.css";

export const metadata: Metadata = {
  title: "Social Network",
  description: "A small social network",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: themeInitScript }} />
      </head>
      <body>
        <ThemeProvider>
          <AuthProvider>{children}</AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
