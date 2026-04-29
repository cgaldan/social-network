"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function DashboardPage() {
  const router = useRouter();
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);
  const [currentUser, setCurrentUser] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("authToken");
    const storedUser = localStorage.getItem("currentUser");

    if (!token) {
      router.replace("/login");
      return;
    }

    if (storedUser) {
      try {
        setCurrentUser(JSON.parse(storedUser));
      } catch (error) {
        localStorage.removeItem("currentUser");
      }
    }

    setIsCheckingAuth(false);
  }, [router]);

  const handleLogout = () => {
    localStorage.removeItem("authToken");
    localStorage.removeItem("currentUser");
    router.replace("/login");
  };

  if (isCheckingAuth) {
    return (
      <main className="dashboard-page">
        <section className="dashboard-card">
          <p className="helper-text">Checking your session...</p>
        </section>
      </main>
    );
  }

  return (
    <main className="dashboard-page">
      <section className="dashboard-card">
        <h1>Welcome, {currentUser?.nickname || "user"}!</h1>
        <p className="helper-text">You are now logged in.</p>
        <div className="dashboard-details">
          <p>
            <strong>Email:</strong> {currentUser?.email || "Not available"}
          </p>
          <p>
            <strong>Name:</strong>{" "}
            {[currentUser?.first_name, currentUser?.last_name].filter(Boolean).join(" ") ||
              "Not available"}
          </p>
        </div>
        <button className="logout-button" type="button" onClick={handleLogout}>
          Logout
        </button>
      </section>
    </main>
  );
}
