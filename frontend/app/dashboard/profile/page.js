"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  deleteCurrentUser,
  getCurrentUser,
  getStoredToken,
  updateCurrentUser,
} from "../../../lib/api";
import { toDateInputValue } from "../../../lib/dates";

export default function ProfilePage() {
  const router = useRouter();
  const token = getStoredToken();
  const [form, setForm] = useState({
    email: "",
    firstName: "",
    lastName: "",
    dateOfBirth: "",
    nickname: "",
    gender: "",
    aboutMe: "",
    isPublic: true,
    avatarPath: "",
  });
  const [error, setError] = useState("");
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const res = await getCurrentUser(token);
        const u = res.user;
        if (!u) return;
        setForm({
          email: u.email || "",
          firstName: u.first_name || "",
          lastName: u.last_name || "",
          dateOfBirth: toDateInputValue(u.date_of_birth),
          nickname: u.nickname || "",
          gender: u.gender || "",
          aboutMe: u.about_me || "",
          isPublic: !!u.is_public,
          avatarPath: u.avatar_path || "",
        });
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    })();
  }, [token]);

  const updateField = (field, value) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setError("");
    setMsg("");
    try {
      const res = await updateCurrentUser(token, {
        email: form.email,
        first_name: form.firstName,
        last_name: form.lastName,
        date_of_birth: `${form.dateOfBirth}T00:00:00Z`,
        nickname: form.nickname,
        gender: form.gender,
        about_me: form.aboutMe,
        is_public: form.isPublic,
        avatar_path: form.avatarPath,
      });
      if (res.user) {
        localStorage.setItem("currentUser", JSON.stringify(res.user));
      }
      setMsg("Profile saved.");
    } catch (e) {
      setError(e.message);
    }
  };

  const handleDelete = async () => {
    if (
      !confirm(
        "Delete your account permanently? This cannot be undone.",
      )
    ) {
      return;
    }
    setError("");
    try {
      await deleteCurrentUser(token);
      localStorage.removeItem("authToken");
      localStorage.removeItem("currentUser");
      router.replace("/login");
    } catch (e) {
      setError(e.message);
    }
  };

  if (loading) {
    return (
      <div className="page-stack">
        <p className="helper-text">Loading profile…</p>
      </div>
    );
  }

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Profile</h1>
        <form className="stack-form login-form" onSubmit={handleSave}>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            value={form.email}
            onChange={(e) => updateField("email", e.target.value)}
            required
          />
          <label htmlFor="firstName">First name</label>
          <input
            id="firstName"
            value={form.firstName}
            onChange={(e) => updateField("firstName", e.target.value)}
            required
          />
          <label htmlFor="lastName">Last name</label>
          <input
            id="lastName"
            value={form.lastName}
            onChange={(e) => updateField("lastName", e.target.value)}
            required
          />
          <label htmlFor="dob">Date of birth</label>
          <input
            id="dob"
            type="date"
            value={form.dateOfBirth}
            onChange={(e) => updateField("dateOfBirth", e.target.value)}
            required
          />
          <label htmlFor="nickname">Nickname</label>
          <input
            id="nickname"
            value={form.nickname}
            onChange={(e) => updateField("nickname", e.target.value)}
            required
          />
          <label htmlFor="gender">Gender</label>
          <select
            id="gender"
            value={form.gender}
            onChange={(e) => updateField("gender", e.target.value)}
            required
          >
            <option value="">Select</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
          <label htmlFor="about">About me</label>
          <textarea
            id="about"
            rows={3}
            value={form.aboutMe}
            onChange={(e) => updateField("aboutMe", e.target.value)}
          />
          <label htmlFor="avatar">Avatar path</label>
          <input
            id="avatar"
            value={form.avatarPath}
            onChange={(e) => updateField("avatarPath", e.target.value)}
            placeholder="optional URL or path"
          />
          <label className="checkbox-row" htmlFor="isPublic">
            <input
              id="isPublic"
              type="checkbox"
              checked={form.isPublic}
              onChange={(e) => updateField("isPublic", e.target.checked)}
            />
            Public profile
          </label>
          {error ? <p className="error-message">{error}</p> : null}
          {msg ? <p className="success-message">{msg}</p> : null}
          <button type="submit">Save changes</button>
        </form>

        <div className="danger-zone">
          <h2>Danger zone</h2>
          <button type="button" className="logout-button" onClick={handleDelete}>
            Delete account
          </button>
        </div>
      </section>
    </div>
  );
}
