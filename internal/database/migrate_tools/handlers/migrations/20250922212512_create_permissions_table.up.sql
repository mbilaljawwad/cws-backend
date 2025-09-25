BEGIN;
		-- Add Up migration here
		CREATE TABLE IF NOT EXISTS permissions (
    id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- create tenant permissions ----------------------------------------------------
INSERT INTO permissions (name) VALUES ('tenant:create'); 
INSERT INTO permissions (name) VALUES ('tenant:read');
INSERT INTO permissions (name) VALUES ('tenant:update');
INSERT INTO permissions (name) VALUES ('tenant:delete');
-- tenant all permissions
INSERT INTO permissions (name) VALUES ('tenant:all');

--  create user permissions ----------------------------------------------------
INSERT INTO permissions (name) VALUES ('user:create');
INSERT INTO permissions (name) VALUES ('user:read');
INSERT INTO permissions (name) VALUES ('user:update');
INSERT INTO permissions (name) VALUES ('user:delete');
--  user all permissions
INSERT INTO permissions (name) VALUES ('user:all');

-- create location permissions ----------------------------------------------------
INSERT INTO permissions (name) VALUES ('location:create');
INSERT INTO permissions (name) VALUES ('location:read');
INSERT INTO permissions (name) VALUES ('location:update');
INSERT INTO permissions (name) VALUES ('location:delete');
--  location all permissions
INSERT INTO permissions (name) VALUES ('location:all');

--  create assets permissions ----------------------------------------------------
INSERT INTO permissions (name) VALUES ('assets:create');
INSERT INTO permissions (name) VALUES ('assets:read');
INSERT INTO permissions (name) VALUES ('assets:update');
INSERT INTO permissions (name) VALUES ('assets:delete');
--  assets all permissions
INSERT INTO permissions (name) VALUES ('assets:all');

	COMMIT;
	