module github.com/B3rs/gork/example/worker

go 1.17

replace github.com/B3rs/gork => ./../..

require (
	github.com/B3rs/gork v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.4
)

require github.com/golang/mock v1.6.0 // indirect
