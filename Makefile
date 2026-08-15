TEST_THREAD_COUNT = 1

.PHONY: test

test:
	go test -parallel=$(TEST_THREAD_COUNT) -v ./...

