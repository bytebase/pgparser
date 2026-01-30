package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single TransactionStmt
func parseTransactionStmt(t *testing.T, input string) *nodes.TransactionStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.TransactionStmt)
	if !ok {
		t.Fatalf("expected *nodes.TransactionStmt, got %T", result.Items[0])
	}
	return stmt
}

// TestTransactionBegin tests: BEGIN
func TestTransactionBegin(t *testing.T) {
	stmt := parseTransactionStmt(t, "BEGIN")

	if stmt.Kind != nodes.TRANS_STMT_BEGIN {
		t.Errorf("expected Kind TRANS_STMT_BEGIN (%d), got %d", nodes.TRANS_STMT_BEGIN, stmt.Kind)
	}
	if stmt.Chain {
		t.Error("expected Chain to be false")
	}
}

// TestTransactionBeginWork tests: BEGIN WORK
func TestTransactionBeginWork(t *testing.T) {
	stmt := parseTransactionStmt(t, "BEGIN WORK")

	if stmt.Kind != nodes.TRANS_STMT_BEGIN {
		t.Errorf("expected Kind TRANS_STMT_BEGIN (%d), got %d", nodes.TRANS_STMT_BEGIN, stmt.Kind)
	}
}

// TestTransactionBeginTransaction tests: BEGIN TRANSACTION
func TestTransactionBeginTransaction(t *testing.T) {
	stmt := parseTransactionStmt(t, "BEGIN TRANSACTION")

	if stmt.Kind != nodes.TRANS_STMT_BEGIN {
		t.Errorf("expected Kind TRANS_STMT_BEGIN (%d), got %d", nodes.TRANS_STMT_BEGIN, stmt.Kind)
	}
}

// TestTransactionStartTransaction tests: START TRANSACTION
func TestTransactionStartTransaction(t *testing.T) {
	stmt := parseTransactionStmt(t, "START TRANSACTION")

	if stmt.Kind != nodes.TRANS_STMT_START {
		t.Errorf("expected Kind TRANS_STMT_START (%d), got %d", nodes.TRANS_STMT_START, stmt.Kind)
	}
}

// TestTransactionCommit tests: COMMIT
func TestTransactionCommit(t *testing.T) {
	stmt := parseTransactionStmt(t, "COMMIT")

	if stmt.Kind != nodes.TRANS_STMT_COMMIT {
		t.Errorf("expected Kind TRANS_STMT_COMMIT (%d), got %d", nodes.TRANS_STMT_COMMIT, stmt.Kind)
	}
	if stmt.Chain {
		t.Error("expected Chain to be false")
	}
}

// TestTransactionCommitWork tests: COMMIT WORK
func TestTransactionCommitWork(t *testing.T) {
	stmt := parseTransactionStmt(t, "COMMIT WORK")

	if stmt.Kind != nodes.TRANS_STMT_COMMIT {
		t.Errorf("expected Kind TRANS_STMT_COMMIT (%d), got %d", nodes.TRANS_STMT_COMMIT, stmt.Kind)
	}
}

// TestTransactionEnd tests: END
func TestTransactionEnd(t *testing.T) {
	stmt := parseTransactionStmt(t, "END")

	if stmt.Kind != nodes.TRANS_STMT_COMMIT {
		t.Errorf("expected Kind TRANS_STMT_COMMIT (%d), got %d", nodes.TRANS_STMT_COMMIT, stmt.Kind)
	}
	if stmt.Chain {
		t.Error("expected Chain to be false")
	}
}

// TestTransactionRollback tests: ROLLBACK
func TestTransactionRollback(t *testing.T) {
	stmt := parseTransactionStmt(t, "ROLLBACK")

	if stmt.Kind != nodes.TRANS_STMT_ROLLBACK {
		t.Errorf("expected Kind TRANS_STMT_ROLLBACK (%d), got %d", nodes.TRANS_STMT_ROLLBACK, stmt.Kind)
	}
	if stmt.Chain {
		t.Error("expected Chain to be false")
	}
}

// TestTransactionRollbackWork tests: ROLLBACK WORK
func TestTransactionRollbackWork(t *testing.T) {
	stmt := parseTransactionStmt(t, "ROLLBACK WORK")

	if stmt.Kind != nodes.TRANS_STMT_ROLLBACK {
		t.Errorf("expected Kind TRANS_STMT_ROLLBACK (%d), got %d", nodes.TRANS_STMT_ROLLBACK, stmt.Kind)
	}
}

// TestTransactionAbort tests: ABORT
func TestTransactionAbort(t *testing.T) {
	stmt := parseTransactionStmt(t, "ABORT")

	if stmt.Kind != nodes.TRANS_STMT_ROLLBACK {
		t.Errorf("expected Kind TRANS_STMT_ROLLBACK (%d), got %d", nodes.TRANS_STMT_ROLLBACK, stmt.Kind)
	}
	if stmt.Chain {
		t.Error("expected Chain to be false")
	}
}

// TestTransactionSavepoint tests: SAVEPOINT sp1
func TestTransactionSavepoint(t *testing.T) {
	stmt := parseTransactionStmt(t, "SAVEPOINT sp1")

	if stmt.Kind != nodes.TRANS_STMT_SAVEPOINT {
		t.Errorf("expected Kind TRANS_STMT_SAVEPOINT (%d), got %d", nodes.TRANS_STMT_SAVEPOINT, stmt.Kind)
	}
	if stmt.Savepoint != "sp1" {
		t.Errorf("expected Savepoint 'sp1', got %q", stmt.Savepoint)
	}
}

// TestTransactionReleaseSavepoint tests: RELEASE SAVEPOINT sp1
func TestTransactionReleaseSavepoint(t *testing.T) {
	stmt := parseTransactionStmt(t, "RELEASE SAVEPOINT sp1")

	if stmt.Kind != nodes.TRANS_STMT_RELEASE {
		t.Errorf("expected Kind TRANS_STMT_RELEASE (%d), got %d", nodes.TRANS_STMT_RELEASE, stmt.Kind)
	}
	if stmt.Savepoint != "sp1" {
		t.Errorf("expected Savepoint 'sp1', got %q", stmt.Savepoint)
	}
}

// TestTransactionRelease tests: RELEASE sp1
func TestTransactionRelease(t *testing.T) {
	stmt := parseTransactionStmt(t, "RELEASE sp1")

	if stmt.Kind != nodes.TRANS_STMT_RELEASE {
		t.Errorf("expected Kind TRANS_STMT_RELEASE (%d), got %d", nodes.TRANS_STMT_RELEASE, stmt.Kind)
	}
	if stmt.Savepoint != "sp1" {
		t.Errorf("expected Savepoint 'sp1', got %q", stmt.Savepoint)
	}
}

// TestTransactionRollbackToSavepoint tests: ROLLBACK TO SAVEPOINT sp1
func TestTransactionRollbackToSavepoint(t *testing.T) {
	stmt := parseTransactionStmt(t, "ROLLBACK TO SAVEPOINT sp1")

	if stmt.Kind != nodes.TRANS_STMT_ROLLBACK_TO {
		t.Errorf("expected Kind TRANS_STMT_ROLLBACK_TO (%d), got %d", nodes.TRANS_STMT_ROLLBACK_TO, stmt.Kind)
	}
	if stmt.Savepoint != "sp1" {
		t.Errorf("expected Savepoint 'sp1', got %q", stmt.Savepoint)
	}
}

// TestTransactionRollbackTo tests: ROLLBACK TO sp1
func TestTransactionRollbackTo(t *testing.T) {
	stmt := parseTransactionStmt(t, "ROLLBACK TO sp1")

	if stmt.Kind != nodes.TRANS_STMT_ROLLBACK_TO {
		t.Errorf("expected Kind TRANS_STMT_ROLLBACK_TO (%d), got %d", nodes.TRANS_STMT_ROLLBACK_TO, stmt.Kind)
	}
	if stmt.Savepoint != "sp1" {
		t.Errorf("expected Savepoint 'sp1', got %q", stmt.Savepoint)
	}
}

// TestTransactionCommitAndChain tests: COMMIT AND CHAIN
func TestTransactionCommitAndChain(t *testing.T) {
	stmt := parseTransactionStmt(t, "COMMIT AND CHAIN")

	if stmt.Kind != nodes.TRANS_STMT_COMMIT {
		t.Errorf("expected Kind TRANS_STMT_COMMIT (%d), got %d", nodes.TRANS_STMT_COMMIT, stmt.Kind)
	}
	if !stmt.Chain {
		t.Error("expected Chain to be true")
	}
}

// TestTransactionRollbackAndNoChain tests: ROLLBACK AND NO CHAIN
func TestTransactionRollbackAndNoChain(t *testing.T) {
	stmt := parseTransactionStmt(t, "ROLLBACK AND NO CHAIN")

	if stmt.Kind != nodes.TRANS_STMT_ROLLBACK {
		t.Errorf("expected Kind TRANS_STMT_ROLLBACK (%d), got %d", nodes.TRANS_STMT_ROLLBACK, stmt.Kind)
	}
	if stmt.Chain {
		t.Error("expected Chain to be false")
	}
}
