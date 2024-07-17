CREATE TABLE IF NOT EXISTS "organizations" (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    owner UUID NOT NULL
);
