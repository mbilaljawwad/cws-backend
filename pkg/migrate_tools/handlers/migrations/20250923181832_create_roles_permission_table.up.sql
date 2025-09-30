BEGIN;

    alter table roles if not exists add constraint cws_roles_name_key unique (name);
    alter table permissions if not exists add constraint cws_permissions_name_key unique (name);

		-- Add Up migration here
		create table if not exists roles_permissions (
			role_name text not null,
			permission_name text not null,
			primary key (role_name, permission_name),
      foreign key (role_name) references roles(name) on delete cascade,
      foreign key (permission_name) references permissions(name) on delete cascade
		);
 
  create index if not exists idx_roles_permissions_role_name on roles_permissions(role_name);
  create index if not exists idx_roles_permissions_permission_name on roles_permissions(permission_name);

		-- Seed role→permission links (safe re-run)
insert into if not exists roles_permissions (role_name, permission_name) values
  -- super_admin (all)
  ('super_admin', 'tenant:all'),
  ('super_admin', 'user:all'),
  ('super_admin', 'location:all'),
  ('super_admin', 'assets:all'),

  -- admin
  ('admin', 'tenant:read'),
  ('admin', 'user:all'),
  ('admin', 'location:all'),
  ('admin', 'assets:all'),

  -- member
  ('member', 'tenant:read'),
  ('member', 'user:create'),
  ('member', 'user:read'),
  ('member', 'user:delete'),
  ('member', 'location:read'),
  ('member', 'assets:read'),

  -- user
  ('user', 'user:read'),
  ('user', 'user:update')
	on conflict do nothing;
	COMMIT;
	