package service_test

import (
	"context"
	"database/sql"
	"ledger/kafka"
	"ledger/mongo"
	"ledger/pg"
	"ledger/service"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────────────────────────────────────

// noOpDB patches the parts of database/sql that the handlers touch.
func noOpDB(p *gomonkey.Patches) {
	// BeginTx returns a dummy *sql.Tx
	p.ApplyMethod(reflect.TypeOf(&sql.DB{}), "BeginTx",
		func(_ *sql.DB, _ context.Context, _ *sql.TxOptions) (*sql.Tx, error) {
			return &sql.Tx{}, nil
		})

	// Commit / Rollback become no-ops
	p.ApplyMethod(reflect.TypeOf(&sql.Tx{}), "Commit", func(_ *sql.Tx) error { return nil })
	p.ApplyMethod(reflect.TypeOf(&sql.Tx{}), "Rollback", func(_ *sql.Tx) error { return nil })
}

// ─────────────────────────────────────────────────────────────────────────────
// tests
// ─────────────────────────────────────────────────────────────────────────────

func TestHandleCreateAccount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	noOpDB(patches)

	// Stub pg.CreateNewAccount
	patches.ApplyFunc(pg.CreateNewAccount,
		func(_ context.Context, user string, amt float64, _ *sql.Tx) error {
			assert.Equal(t, "user-1", user)
			assert.Equal(t, 100.0, amt)
			return nil
		})

	// Stub mongo.RecordTransaction
	patches.ApplyFunc(mongo.RecordTransaction,
		func(_ context.Context, rec []mongo.LedgerRecord) error {
			assert.Len(t, rec, 1)
			assert.Equal(t, "CreateAccount", rec[0].Operation)
			return nil
		})

	msg := kafka.CreateAccountMessage{InitialBalance: 100, BaseMessage: kafka.BaseMessage{UserID: "user-1"}}
	err := service.HandleCreateAccount(msg)

	assert.NoError(t, err)
}

func TestHandleAddBalance(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	noOpDB(patches)

	// Stub pg.UpdateBalance
	patches.ApplyFunc(pg.UpdateBalance,
		func(_ context.Context, user string, delta float64, _ *sql.Tx) error {
			assert.Equal(t, "user-2", user)
			assert.Equal(t, 50.0, delta)
			return nil
		})

	// Stub mongo.RecordTransaction
	patches.ApplyFunc(mongo.RecordTransaction,
		func(_ context.Context, rec []mongo.LedgerRecord) error {
			assert.Equal(t, "AddBalance", rec[0].Operation)
			return nil
		})

	msg := kafka.AddBalanceMessage{Amount: 50, BaseMessage: kafka.BaseMessage{UserID: "user-2"}}
	err := service.HandleAddBalance(msg)

	assert.NoError(t, err)
}

func TestHandleDeductBalance(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	noOpDB(patches)

	// Stub pg.UpdateBalance (negative amount is passed in handler)
	patches.ApplyFunc(pg.UpdateBalance,
		func(_ context.Context, user string, delta float64, _ *sql.Tx) error {
			assert.Equal(t, "user-3", user)
			assert.Equal(t, -25.0, delta)
			return nil
		})

	// Stub mongo.RecordTransaction
	patches.ApplyFunc(mongo.RecordTransaction,
		func(_ context.Context, rec []mongo.LedgerRecord) error {
			assert.Equal(t, "DeductBalance", rec[0].Operation)
			return nil
		})

	msg := kafka.DeductBalanceMessage{Amount: 25, BaseMessage: kafka.BaseMessage{UserID: "user-3"}}
	err := service.HandleDeductBalance(msg)

	assert.NoError(t, err)
}
