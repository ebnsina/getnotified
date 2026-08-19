create extension if not exists pgcrypto;

create table if not exists organizations (
  id         uuid primary key default gen_random_uuid(),
  name       text not null,
  slug       text not null unique,
  created_at timestamptz not null default now()
);

create table if not exists monitors (
  id                   uuid primary key default gen_random_uuid(),
  org_id               uuid not null references organizations(id) on delete cascade,
  name                 text not null,
  type                 text not null default 'http',   -- http | tcp | ssl_expiry
  target               text not null,
  interval_seconds     int  not null default 60,
  timeout_seconds      int  not null default 10,
  expected_status      int[] not null default '{200}',
  ssl_warn_days        int  not null default 14,
  failure_threshold    int  not null default 2,
  tags                 text[] not null default '{}',
  paused               bool not null default false,
  status               text not null default 'pending', -- pending | up | down
  consecutive_failures int  not null default 0,
  next_check_at        timestamptz not null default now(),
  created_at           timestamptz not null default now()
);

create table if not exists checks (
  id          bigserial primary key,
  monitor_id  uuid not null references monitors(id) on delete cascade,
  checked_at  timestamptz not null default now(),
  ok          bool not null,
  latency_ms  int  not null,
  status_code int,
  error       text,
  region      text not null default 'local'
);
create index if not exists checks_monitor_time on checks (monitor_id, checked_at desc);

create table if not exists incidents (
  id          uuid primary key default gen_random_uuid(),
  monitor_id  uuid not null references monitors(id) on delete cascade,
  started_at  timestamptz not null default now(),
  resolved_at timestamptz,
  cause       text
);
create unique index if not exists incidents_one_open on incidents (monitor_id) where resolved_at is null;

create table if not exists notification_channels (
  id         uuid primary key default gen_random_uuid(),
  org_id     uuid not null references organizations(id) on delete cascade,
  name       text not null,
  type       text not null,  -- slack | email | webhook | sms | whatsapp | imessage
  config     jsonb not null default '{}',
  created_at timestamptz not null default now()
);

create table if not exists monitor_channels (
  monitor_id uuid not null references monitors(id) on delete cascade,
  channel_id uuid not null references notification_channels(id) on delete cascade,
  primary key (monitor_id, channel_id)
);

-- Row-level security. Policies bind to getnotified_app, so the owner role that
-- runs migrations is unaffected and the app role is confined to app.org_id.
do $$ begin
  if not exists (select from pg_roles where rolname = 'getnotified_app') then
    create role getnotified_app nologin;
  end if;
end $$;

-- Returns null when app.org_id is unset, so policies fail closed.
create or replace function current_org() returns uuid language sql stable as $$
  select nullif(current_setting('app.org_id', true), '')::uuid
$$;

do $$
declare t text;
begin
  foreach t in array array['organizations', 'monitors', 'checks', 'incidents',
                           'notification_channels', 'monitor_channels'] loop
    execute format('alter table %I enable row level security', t);
    execute format('drop policy if exists org_isolation on %I', t);
  end loop;
end $$;

create policy org_isolation on organizations
  for all to getnotified_app using (id = current_org()) with check (id = current_org());

create policy org_isolation on monitors
  for all to getnotified_app using (org_id = current_org()) with check (org_id = current_org());

create policy org_isolation on notification_channels
  for all to getnotified_app using (org_id = current_org()) with check (org_id = current_org());

-- checks, incidents and monitor_channels inherit isolation through monitors,
-- which is itself filtered by the policy above.
create policy org_isolation on checks
  for all to getnotified_app using (monitor_id in (select id from monitors))
  with check (monitor_id in (select id from monitors));

create policy org_isolation on incidents
  for all to getnotified_app using (monitor_id in (select id from monitors))
  with check (monitor_id in (select id from monitors));

create policy org_isolation on monitor_channels
  for all to getnotified_app using (monitor_id in (select id from monitors))
  with check (monitor_id in (select id from monitors));
