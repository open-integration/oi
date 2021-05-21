
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



# Compile local oictl 
.PHONY: build-oictl
build-oictl::
	go build -o ./dist/oictl ./cmd/oictl


.PHONY: build-service-catalog
build-service-catalog:
	@bash ./scripts/build-service-catalog.sh

.PHONY: release-service-catalog
release-service-catalog:
	@bash ./scripts/release-service-catalog.sh

.PHONY: build-testing-image
build-testing-image:
	docker build . -f testing/Dockerfile -t openintegration/testing

.PHONY: generate-catalog-types
generate-catalog-types:
	# TODO: get all the data from service.yaml files
	mkdir -p catalog/services/github/endpoints/issuesearch
	quicktype -o catalog/services/github/endpoints/issuesearch/arguments.go -l go -s schema --src catalog/services/github/configs/endpoints/issuesearch/arguments.json --package issuesearch -t IssueSearchArguments
	quicktype -o catalog/services/github/endpoints/issuesearch/returns.go -l go -s schema --src catalog/services/github/configs/endpoints/issuesearch/returns.json --package issuesearch -t IssueSearchReturns
	make fmt