create table jobs (
  id varchar primary key,
  queue varchar not null,
  status varchar not null,
  arguments jsonb not null,
  result jsonb,
  last_error varchar,
  scheduled_at timestamptz default now(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- index on queue
-- index on scheduled_at
-- index on status

-- insert into jobs (queue, status, arguments) values ('default', 'scheduled', '{"wow":123}');
-- insert into jobs (queue, status, arguments) values ('default', 'completed', '{"wow":152}');
-- insert into jobs (queue, status, arguments) values ('default', 'scheduled', '{"wow":365}');