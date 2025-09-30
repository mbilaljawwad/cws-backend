BEGIN;
		-- Add Down migration here
		DROP TABLE IF EXISTS super_admin;
	COMMIT;
	