const FALLBACK_API_BASE_URL = "http://localhost:8080";

export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || FALLBACK_API_BASE_URL;

export async function login({ identifier, password }) {
  const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      identifier,
      password,
    }),
  });

  let payload = null;

  try {
    payload = await response.json();
  } catch (error) {
    throw new Error("Unable to parse server response.");
  }

  if (!response.ok || !payload.success) {
    throw new Error(payload?.message || "Login failed.");
  }

  return payload;
}

export async function register({
  email,
  password,
  firstName,
  lastName,
  dateOfBirth,
  nickname,
  gender,
  aboutMe = "",
  isPublic = true,
  avatarPath = "",
}) {
  const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      password,
      first_name: firstName,
      last_name: lastName,
      date_of_birth: `${dateOfBirth}T00:00:00Z`,
      nickname,
      gender,
      about_me: aboutMe,
      is_public: isPublic,
      avatar_path: avatarPath,
    }),
  });

  let payload = null;

  try {
    payload = await response.json();
  } catch (error) {
    throw new Error("Unable to parse server response.");
  }

  if (!response.ok || !payload.success) {
    throw new Error(payload?.message || "Registration failed.");
  }

  return payload;
}
