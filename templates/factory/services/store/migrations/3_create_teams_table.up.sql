CREATE TABLE IF NOT EXISTS "teams" (
   id UUID PRIMARY KEY,
   name TEXT NOT NULL,
   organization_id UUID NOT NULL,
   FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS "team_users" (
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
