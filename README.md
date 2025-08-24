## Database Sketch (PostgreSQL)
> Tooling: `sqlc` for typed queries (pgx), goose/migrate for migrations.

  ```sql
-- Tenancy & Identity
CREATE TABLE tenant (
  id UUID PK DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  billing_email TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE location (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL REFERENCES tenant(id),
  name TEXT, timezone TEXT, address JSONB, created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE team (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, location_id UUID,
  name TEXT NOT NULL, metadata JSONB DEFAULT '{}'::jsonb
);
CREATE TABLE app_user (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL,
  email CITEXT UNIQUE, full_name TEXT, status TEXT CHECK (status IN ('invited','active','suspended')),
  created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE role (
  id SERIAL PK, name TEXT UNIQUE
);
CREATE TABLE user_role (
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  role_id INT REFERENCES role(id), tenant_id UUID NOT NULL, PRIMARY KEY(user_id, role_id, tenant_id)
);

-- Resources & Assets
CREATE TABLE resource_type (id SERIAL PK, code TEXT UNIQUE, name TEXT);
CREATE TABLE resource (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, location_id UUID NOT NULL,
  type_id INT REFERENCES resource_type(id), name TEXT, capacity INT, features JSONB,
  open_hours JSONB, pricing JSONB, status TEXT CHECK (status IN ('active','maintenance','retired'))
);
CREATE TABLE asset (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, location_id UUID NOT NULL,
  tag TEXT UNIQUE, kind TEXT, status TEXT, warranty_until DATE, meta JSONB
);

-- Plans & Memberships
CREATE TABLE plan (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL,
  name TEXT, period TEXT CHECK (period IN ('monthly','annual')),
  price_cents INT, currency TEXT, credits JSONB, addons JSONB, tax_rate NUMERIC(5,2)
);
CREATE TABLE membership (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, team_id UUID, user_id UUID,
  plan_id UUID REFERENCES plan(id), status TEXT CHECK (status IN ('trial','active','paused','canceled')),
  starts_at DATE, ends_at DATE, next_renewal DATE
);

-- Bookings
CREATE TABLE booking (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, resource_id UUID NOT NULL,
  booker_id UUID NOT NULL, start_at TIMESTAMPTZ, end_at TIMESTAMPTZ,
  status TEXT CHECK (status IN ('tentative','confirmed','checked_in','checked_out','canceled')),
  recurrence JSONB, attendees JSONB
);

-- Billing
CREATE TABLE invoice (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, team_id UUID, user_id UUID,
  number TEXT UNIQUE, issued_at TIMESTAMPTZ, due_at TIMESTAMPTZ, status TEXT,
  subtotal_cents INT, tax_cents INT, total_cents INT, currency TEXT
);
CREATE TABLE invoice_line (
  id BIGSERIAL PK, invoice_id UUID REFERENCES invoice(id) ON DELETE CASCADE,
  description TEXT, quantity NUMERIC, unit_price_cents INT, total_cents INT, meta JSONB
);
CREATE TABLE payment (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, invoice_id UUID,
  provider TEXT, provider_ref TEXT, amount_cents INT, currency TEXT, status TEXT,
  received_at TIMESTAMPTZ
);

-- Support
CREATE TABLE ticket (
  id UUID PK DEFAULT gen_random_uuid(), tenant_id UUID NOT NULL, location_id UUID,
  category TEXT, title TEXT, description TEXT, priority TEXT, assignee_id UUID, status TEXT,
  sla_due TIMESTAMPTZ, created_by UUID, created_at TIMESTAMPTZ DEFAULT now()
);

-- Audit
CREATE TABLE audit_log (
  id BIGSERIAL PK, tenant_id UUID, actor_id UUID, action TEXT, entity TEXT, entity_id UUID,
  diff JSONB, created_at TIMESTAMPTZ DEFAULT now()
);
```
---

## API Design (External REST, versioned)


### Auth
- `POST /v1/auth/login` (if first‑party auth)
- `GET /v1/auth/me`

### Organizations & Users
- `GET/POST /v1/{t}/users`, `PATCH /v1/{t}/users/{id}`, `POST /v1/{t}/users/{id}:invite`
- `GET/POST /v1/{t}/teams`, add/remove members
- `GET/POST /v1/{t}/locations`
- `GET /v1/{t}/roles`, `PUT /v1/{t}/users/{id}/roles`

### Resources & Assets
- `GET/POST /v1/{t}/resources`, filters by type/location
- `PATCH /v1/{t}/resources/{id}`
- `GET/POST /v1/{t}/assets`, `PATCH /v1/{t}/assets/{id}`

### Plans & Memberships
- `GET/POST /v1/{t}/plans`
- `GET/POST /v1/{t}/memberships`, pause/cancel actions

### Bookings
- `POST /v1/{t}/bookings:availability` (check conflicts, pricing)
- `GET/POST /v1/{t}/bookings`, `PATCH /v1/{t}/bookings/{id}` (confirm/cancel/check‑in/out)

### Billing
- `GET/POST /v1/{t}/invoices`, `GET /v1/{t}/invoices/{id}` (PDF link)
- `POST /v1/{t}/payments/intents`, `POST /v1/{t}/payments/webhook` (provider adapter)

### Support & Ops
- `GET/POST /v1/{t}/tickets`, status changes, comments

### Search & Reports
- `GET /v1/{t}/search?q=…`
- `GET /v1/{t}/reports/occupancy?from=&to=&location=`

---

## Backend (Go) — Patterns & Skeleton
**Style**: Clean Architecture + DDD‑ish packages per context. Prefer **chi** for HTTP, **sqlc + pgx** for data, **otel** for tracing.
**Note**: This architecture may subject to change depending upon the requirements.

```
/erp
  /cmd
    /api (main)
  /internal
    /auth
    /org
    /rbac
    /resource
    /booking
    /plan
    /billing
    /support
    /audit
    /platform (db, cache, mq, config, logger)
  /migrations
  /api (openapi specs)
```

## RBAC Matrix (Excerpt)
| Resource      | SuperAdmin | OrgOwner | Manager | Member | Support | Finance |
|---------------|------------|----------|---------|--------|---------|---------|
| Orgs          | CRUD       | R/U      | R       | —      | R       | R       |
| Users/Teams   | R          | CRUD     | R       | R(self)| R       | R       |
| Spaces        | R          | CRUD     | CRUD    | R      | R       | R       |
| Bookings      | R          | CRUD     | CRUD    | CRUD(self)| R     | R       |
| Memberships   | R          | CRUD     | R       | R(self)| R       | R       |
| Invoices      | R          | R        | —       | R(own) | R       | CRUD    |
| Assets        | R          | R        | CRUD    | —      | R       | R       |
| Tickets       | R          | R        | R       | R(own) | CRUD    | R       |
| Audit Logs    | R          | R        | —       | —      | —       | —       |  

---

Store as `role_permissions(role_id, permission_code)` and cache in Redis.

