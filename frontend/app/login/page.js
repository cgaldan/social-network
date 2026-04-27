"use client";

import { useState } from "react";
import Link from "next/link";
import { login } from "../../lib/api";

export default function LoginPage() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setErrorMessage("");
    setSuccessMessage("");

    if (!identifier.trim() || !password.trim()) {
      setErrorMessage("Please enter both identifier and password.");
      return;
    }

    setIsSubmitting(true);

    try {
      const payload = await login({ identifier, password });

      localStorage.setItem("authToken", payload.token);
      localStorage.setItem("currentUser", JSON.stringify(payload.user));

      setSuccessMessage(`Welcome back, ${payload.user?.nickname || "user"}!`);
      setPassword("");
    } catch (error) {
      setErrorMessage(error.message || "Unable to login.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className="login-page">
      <section className="login-card">
        <h1>Login</h1>
        <p className="helper-text">Use your nickname or email to sign in.</p>
        <p className="switch-auth-text">
          New here? <Link href="/register">Create an account</Link>
        </p>

        <form className="login-form" onSubmit={handleSubmit}>
          <label htmlFor="identifier">Identifier</label>
          <input
            id="identifier"
            type="text"
            name="identifier"
            value={identifier}
            onChange={(event) => setIdentifier(event.target.value)}
            placeholder="email or nickname"
            autoComplete="username"
            disabled={isSubmitting}
          />

          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            name="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="your password"
            autoComplete="current-password"
            disabled={isSubmitting}
          />

          {errorMessage ? <p className="error-message">{errorMessage}</p> : null}
          {successMessage ? <p className="success-message">{successMessage}</p> : null}

          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Logging in..." : "Login"}
          </button>
        </form>
      </section>
    </main>
  );
}
