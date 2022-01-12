CREATE TABLE jobs (
  id varchar primary key,
  queue varchar not null,
  status varchar not null,
  arguments jsonb not null default '{}'::jsonb,
  result jsonb not null default '{}'::jsonb,
  last_error varchar,
  retry_count integer not null default 0,
  options jsonb not null default '{}'::jsonb,
  scheduled_at timestamptz default now(),
  started_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- index on queue
CREATE INDEX ON jobs (queue);
-- index on scheduled_at
CREATE INDEX ON jobs (scheduled_at);
-- index on status
CREATE INDEX ON jobs (status);
-- index on started_at
CREATE INDEX ON jobs (started_at);