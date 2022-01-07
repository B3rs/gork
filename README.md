# gork
golang pgsql backed worker library

## Introduction
Did you ever dream of executing code at a given time?
Did you ever dream about calling a third party but not blocking kafka consumers because of some broken constraint?
Did you have any concerns about "transactionality" in calling third party services?

gork is the answer

gork uses postgresql tables to create a queue of jobs that can be executed, re-executed, scheduled and retryed all within your beloved transactions.

see examples to see it in action.

### TODO
- [ ] proper documentation
- [ ] database indexes tuning
- [ ] improve code structure
- [ ] admin ui
- [ ] metrics and alerts
- [ ] job retry
- [ ] failure notifications (investigate sentry)
- [ ] performance benchmarks
