import "./globals.css";

export const metadata = {
  title: "Social Network",
  description: "Login to your social network account",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
