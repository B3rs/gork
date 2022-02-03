# gork

golang, SQL based, simple enough, worker library

## Introduction

Gork is a worker library based on plain SQL instructions, it allows you to schedule and execute jobs in a simple and even transactional way (if needed).
In comparison to machinery or other frameworks gork focuses on simplicity and visibility instead of performances (while being performant enough).
Who needs redis when you already have an SQL database?

see examples to see it in action.

## How it works

gork under the hood uses a table like the following.

```
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

CREATE INDEX ON jobs (queue);
CREATE INDEX ON jobs (scheduled_at);
CREATE INDEX ON jobs (status);
CREATE INDEX ON jobs (started_at);
```

Every worker registered in the pool will poll the database with a

```
FOR UPDATE SKIP LOCKED
```

query to dequeue a job to be executed.

If the job gets stuck or the worker crashes a "reaper" will reschedule the job for you.

Gork jobs semantic is "AT LEAST ONCE"

## See it in action!

You can create jobs here
https://gork-client-example.herokuapp.com/

And check what happens here!
https://gork-worker-example.herokuapp.com/

### TODO

- [ ] proper documentation
- [ ] unit tests
- [ ] api tests
- [ ] metrics and alerts
- [ ] performance benchmarks
- [ ] workers statistics
