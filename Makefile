.PHONY: test-all
test-all:
	make test
	make test-race

.PHONY: test
test:
	go clean -testcache
	go test -v ./... --cover

.PHONY: test-race
test-race:
	go clean -testcache
	go test -v ./... --cover --race
