BEGIN;
		-- Add Down migration here
		TRUNCATE TABLE roles CASCADE;
		
	COMMIT;
	