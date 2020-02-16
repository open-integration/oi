.PHONY: test
test:
	@sh ./scripts/test.sh

test-fmt:
	@sh ./scripts/test-fmt.sh

fmt:
	go fmt ./...