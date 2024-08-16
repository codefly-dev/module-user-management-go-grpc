CREATE TABLE IF NOT EXISTS "users"
(
    id            UUID PRIMARY KEY,
    signed_up_at  DATE,
    last_login_at DATE,
    status        TEXT,
    email         TEXT,
    profile       JSONB
);


CREATE INDEX idx_users_id ON users (id);
CREATE INDEX idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS "users_auth"
(
    user_id UUID NOT NULL,
    auth_id TEXT NOT NULL,
    PRIMARY KEY (user_id, auth_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_auth_user_id ON users_auth (user_id);
CREATE INDEX idx_user_auth_auth_id ON users_auth (auth_id);


CREATE TABLE IF NOT EXISTS "permissions"
(
    id       UUID PRIMARY KEY,
    name     TEXT NOT NULL,
    resource TEXT NOT NULL,
    access   TEXT NOT NULL
);

CREATE INDEX idx_permissions_id ON permissions (id);

CREATE TABLE IF NOT EXISTS "roles"
(
    id   UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE INDEX idx_roles_id ON roles (id);


CREATE TABLE IF NOT EXISTS "role_permissions"
(
    role_id       UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions (permission_id);
