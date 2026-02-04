# pgparser - Pure Go PostgreSQL Parser

`pgparser` is a **pure Go** implementation of the PostgreSQL SQL parser. It translates the original PostgreSQL grammar (`gram.y`) into Go using `goyacc`, ensuring **100% compatibility** with the PostgreSQL parse tree structure.

> **Target PostgreSQL Version**: 17.7 (REL_17_STABLE)

## Features

- **100% Native Go**: No CGO, no external dependencies.
- **AST Compatibility**: The generated Abstract Syntax Tree (AST) strictly matches PostgreSQL's internal `Node` structures (`parsenodes.h`), enabling seamless integration with other PG-compatible tools.
- **Battle-Tested**: Verified against PostgreSQL's own massive regression test suite (~45,000 SQL statements) with a **99.6% pass rate**.
- **Thread Safe**: The parser is designed to be safe for concurrent use.

## Installation

```bash
go get github.com/bytebase/pgparser
```

## Quick Start

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bytebase/pgparser/parser"
)

func main() {
	sql := "SELECT * FROM users WHERE id > 100"

	// Parse the SQL string
	stmts, err := parser.Parse(sql)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	// Print the AST as JSON
	jsonBytes, _ := json.MarshalIndent(stmts, "", "  ")
	fmt.Println(string(jsonBytes))
}
```

## Architecture

This project is not a hand-written parser. It is a **port** of the official PostgreSQL source code:

1.  **Grammar**: `parser/gram.y` is a translated version of PostgreSQL's `src/backend/parser/gram.y`.
2.  **Lexer**: `parser/scan.l` logic is reimplemented in Go in `parser/lexer.go`.
3.  **Nodes**: `nodes/` contains Go struct definitions that map 1:1 to PostgreSQL's C structs.

## Development

### Requirements

- Go 1.18+
- `goyacc` (installed via `go install golang.org/x/tools/cmd/goyacc@latest`)

### Build & Test

```bash
# Regenerate parser (if gram.y changed)
make generate-parser

# Run tests
go test ./...

# Run PostgreSQL regression tests compatibility check
go test ./parser/pgregress -run TestPGRegressStats -v
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=bytebase/pgparser&type=date&legend=top-left)](https://www.star-history.com/#bytebase/pgparser&type=date&legend=top-left)

## License

This project is licensed under the MIT License. Portions derived from PostgreSQL are subject to the PostgreSQL License.
