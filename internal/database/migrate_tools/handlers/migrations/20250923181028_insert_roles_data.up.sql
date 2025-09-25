BEGIN;
		-- Add Up migration here
		-- super admin that would have all permissions
		INSERT INTO roles (name) VALUES ('super_admin');
		-- admin that would have tenant permissions
		INSERT INTO roles (name) VALUES ('admin'); 
		-- member that would have member permissions
		INSERT INTO roles (name) VALUES ('member');
		-- user that would have user permissions
		INSERT INTO roles (name) VALUES ('user');
	COMMIT;
	