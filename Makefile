
.PHONY: test-all
test-all: test-dir test-fmt gocyclo spellcheck lint security-check test

.PHONY: test
test:
	@sh ./scripts/test.sh

test-fmt:
	@sh ./scripts/test-fmt.sh

test-dir:
	@bash ./scripts/test-clean-dir.sh

# Fix fmt errors in file
fmt:
	go fmt ./...

release:
	@sh ./scripts/release.sh

# spellcheck Finds commonly misspelled English words
.PHONY: spellcheck
spellcheck:
	@misspell -error .

# Gocyclo calculates cyclomatic complexities of functions in Go source code.
# The cyclomatic complexity of a function is calculated according to the following rules: 
# 1 is the base complexity of a function +1 for each 'if', 'for', 'case', '&&' or '||'
# Go Report Card warns on functions with cyclomatic complexity > 15.
.PHONY: gocyclo
gocyclo:
	@gocyclo -over 15 .

.PHONY: lint
lint:
	@staticcheck ./...

.PHONY: security-check
security-check:
	@gosec ./... -nosec


# Generate mock struct from interface
# example: make mock PKG=./pkg/runtime NAME=Runtime
.PHONY: mock
mock:
	@sh ./scripts/mock.sh $(PKG) $(NAME)

export-vars:
	@sh ./scripts/export-vars.sh



.PHONY: build
build:
	@bash ./scripts/build.sh


.PHONY: release-binaries
release-binaries:
	@bash ./scripts/release-binaries.sh

.PHONY: build-testing-image
build-testing-image:
	docker build . -f testing/Dockerfile -t openintegration/testing

.PHONY: generate-catalog-types
generate-catalog-types:
	@bash ./scripts/code-gen.sh
	make fmt