BEGIN;
		-- Add Up migration here
		CREATE TABLE IF NOT EXISTS super_admin (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			password TEXT NOT NULL,
			last_login_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

	INSERT INTO IF NOT EXISTS super_admin (email, first_name, last_name, password, last_login_at) VALUES ('superadmin@example.com', 'Super', 'Admin', 'password', NOW());

	COMMIT;
	