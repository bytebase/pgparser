.PHONY: all build generate test clean

# PostgreSQL source location
PG_SRC := $(HOME)/Github/postgres

all: generate build test

build:
	go build ./...

generate: generate-parser generate-nodes

generate-parser:
	cd parser && go generate

generate-nodes:
	go run ./cmd/pgsema-gen/... \
		-parsenodes $(PG_SRC)/src/include/nodes/parsenodes.h \
		-primnodes $(PG_SRC)/src/include/nodes/primnodes.h \
		-outdir nodes

test:
	go test ./...

clean:
	rm -f parser/y.go
	rm -f nodes/*_generated.go

# Install goyacc if not present
install-tools:
	go install golang.org/x/tools/cmd/goyacc@latest

# Copy PostgreSQL regression test SQL files
copy-regress-sql:
	mkdir -p parser/pgregress/testdata/sql
	cp $(PG_SRC)/src/test/regress/sql/*.sql parser/pgregress/testdata/sql/

# Update known_failures.json baseline
update-known-failures:
	cd parser/pgregress && go test -run TestPGRegress -update -v -count=1

# Analyze PostgreSQL grammar (for development)
analyze-gram:
	@echo "Grammar lines: $$(wc -l < $(PG_SRC)/src/backend/parser/gram.y)"
	@echo "makeNode calls: $$(grep -c 'makeNode(' $(PG_SRC)/src/backend/parser/gram.y)"
	@echo "list_make calls: $$(grep -c 'list_make' $(PG_SRC)/src/backend/parser/gram.y)"
	@echo "Keywords: $$(grep -c 'PG_KEYWORD' $(PG_SRC)/src/include/parser/kwlist.h)"
