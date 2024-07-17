CREATE TABLE IF NOT EXISTS "users" (
    id UUID PRIMARY KEY,
    auth_signup_id TEXT NOT NULL,
    signed_up_at DATE,
    last_login_at DATE,
    status TEXT,
    email TEXT NOT NULL,
    profile JSONB
);
