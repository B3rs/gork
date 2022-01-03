create table jobs (
  id bigserial primary key,
  queue text not null,
  status text not null,
  arguments jsonb not null,
  result jsonb,
  last_error text,
  scheduled_at timestamptz default now(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

insert into jobs (queue, status, arguments) values ('default', 'scheduled', '{"wow":123}');
insert into jobs (queue, status, arguments) values ('default', 'completed', '{"wow":152}');
insert into jobs (queue, status, arguments) values ('default', 'scheduled', '{"wow":365}');