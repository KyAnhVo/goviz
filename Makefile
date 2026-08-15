TEST_THREAD_COUNT = 1

.PHONY: test lint

test:
	@ go test -parallel=$(TEST_THREAD_COUNT) -v ./... | grep -v 'no test files'

lint:
	@ echo 'LINTING'
	@ - golangci-lint run
	@ echo -e '\n----------------------------------------------------------------------------\n'
	@ echo 'VULNERABILITY CHECK'
	@ - govulncheck ./...

