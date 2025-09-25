BEGIN;
		-- Add Down migration here
		alter table roles drop constraint if exists cws_roles_name_key;
		alter table permissions drop constraint if exists cws_permissions_name_key;
		drop table if exists roles_permissions;
		drop index if exists idx_roles_permissions_role_name;
		drop index if exists idx_roles_permissions_permission_name;
	COMMIT;
	