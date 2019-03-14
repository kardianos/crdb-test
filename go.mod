module github.com/kardianos/crdb-test

go 1.12

require (
	bitbucket.org/kardianos/table v0.0.4
	github.com/lib/pq v1.0.0
)

replace github.com/lib/pq => ./pq
