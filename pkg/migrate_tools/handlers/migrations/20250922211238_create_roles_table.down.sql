BEGIN;
		-- Add Down migration here
		DROP TABLE IF EXISTS roles;
	COMMIT;
	