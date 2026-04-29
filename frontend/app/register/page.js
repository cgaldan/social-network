"use client";

import { useState } from "react";
import Link from "next/link";
import { register } from "../../lib/api";

export default function RegisterPage() {
  const [form, setForm] = useState({
    email: "",
    password: "",
    firstName: "",
    lastName: "",
    dateOfBirth: "",
    nickname: "",
    gender: "",
    aboutMe: "",
    isPublic: true,
  });
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const updateField = (field, value) => {
    setForm((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setErrorMessage("");
    setSuccessMessage("");

    if (
      !form.email.trim() ||
      !form.password.trim() ||
      !form.firstName.trim() ||
      !form.lastName.trim() ||
      !form.dateOfBirth ||
      !form.nickname.trim() ||
      !form.gender
    ) {
      setErrorMessage("Please fill all required fields.");
      return;
    }

    setIsSubmitting(true);

    try {
      const payload = await register(form);
      localStorage.setItem("authToken", payload.token);
      localStorage.setItem("currentUser", JSON.stringify(payload.user));
      setSuccessMessage(
        `Registration successful. Welcome, ${payload.user?.nickname || "user"}!`,
      );
      setForm((prev) => ({
        ...prev,
        password: "",
      }));
    } catch (error) {
      setErrorMessage(error.message || "Unable to register.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className="login-page">
      <section className="login-card">
        <h1>Create account</h1>
        <p className="helper-text">Register to join the social network.</p>
        <p className="switch-auth-text">
          Already have an account? <Link href="/login">Login</Link>
        </p>

        <form className="login-form" onSubmit={handleSubmit}>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            value={form.email}
            onChange={(event) => updateField("email", event.target.value)}
            placeholder="you@example.com"
            autoComplete="email"
            disabled={isSubmitting}
          />

          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            value={form.password}
            onChange={(event) => updateField("password", event.target.value)}
            placeholder="minimum 6 characters"
            autoComplete="new-password"
            disabled={isSubmitting}
          />

          <label htmlFor="firstName">First name</label>
          <input
            id="firstName"
            type="text"
            value={form.firstName}
            onChange={(event) => updateField("firstName", event.target.value)}
            placeholder="first name"
            disabled={isSubmitting}
          />

          <label htmlFor="lastName">Last name</label>
          <input
            id="lastName"
            type="text"
            value={form.lastName}
            onChange={(event) => updateField("lastName", event.target.value)}
            placeholder="last name"
            disabled={isSubmitting}
          />

          <label htmlFor="dateOfBirth">Date of birth</label>
          <input
            id="dateOfBirth"
            type="date"
            value={form.dateOfBirth}
            onChange={(event) => updateField("dateOfBirth", event.target.value)}
            disabled={isSubmitting}
          />

          <label htmlFor="nickname">Nickname</label>
          <input
            id="nickname"
            type="text"
            value={form.nickname}
            onChange={(event) => updateField("nickname", event.target.value)}
            placeholder="nickname"
            disabled={isSubmitting}
          />

          <label htmlFor="gender">Gender</label>
          <select
            id="gender"
            value={form.gender}
            onChange={(event) => updateField("gender", event.target.value)}
            disabled={isSubmitting}
          >
            <option value="">Select gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>

          <label htmlFor="aboutMe">About me (optional)</label>
          <textarea
            id="aboutMe"
            rows={3}
            value={form.aboutMe}
            onChange={(event) => updateField("aboutMe", event.target.value)}
            placeholder="Tell us about yourself"
            disabled={isSubmitting}
          />

          <label className="checkbox-row" htmlFor="isPublic">
            <input
              id="isPublic"
              type="checkbox"
              checked={form.isPublic}
              onChange={(event) => updateField("isPublic", event.target.checked)}
              disabled={isSubmitting}
            />
            Public profile
          </label>

          {errorMessage ? <p className="error-message">{errorMessage}</p> : null}
          {successMessage ? <p className="success-message">{successMessage}</p> : null}

          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Creating account..." : "Register"}
          </button>
        </form>
      </section>
    </main>
  );
}
