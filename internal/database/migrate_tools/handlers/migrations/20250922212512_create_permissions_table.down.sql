BEGIN;
		-- Add Down migration here
		DROP TABLE IF EXISTS permissions;
	COMMIT;
	