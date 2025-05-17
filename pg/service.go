package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// UpdateBalance updates the user's balance if the user exists.
// It errors if the user does not already exist.
func UpdateBalance(ctx context.Context, userID string, amount float64, tx *sql.Tx) error {
	var err error
	internalTx := false

	if tx == nil {
		tx, err = DB.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		internalTx = true
	}

	query := `
		UPDATE user_balances
		SET balance = balance + $2
		WHERE user_id = $1
	`
	result, err := tx.ExecContext(ctx, query, userID, amount)
	if err != nil {
		if internalTx {
			_ = tx.Rollback()
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if internalTx {
			_ = tx.Rollback()
		}
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		if internalTx {
			_ = tx.Rollback()
		}
		return errors.New("no such user found to update balance")
	}

	if internalTx {
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}

// CreateNewAccount inserts a new user balance.
// It errors if the user already exists.
func CreateNewAccount(ctx context.Context, userID string, amount float64, tx *sql.Tx) error {
	var err error
	internalTx := false

	if tx == nil {
		tx, err = DB.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		internalTx = true
	}

	query := `
		INSERT INTO user_balances(user_id, balance)
		VALUES ($1, $2)
	`
	_, err = tx.ExecContext(ctx, query, userID, amount)
	if err != nil {
		if internalTx {
			_ = tx.Rollback()
		}
		return fmt.Errorf("failed to create new account: %w", err)
	}

	if internalTx {
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}

func GetBalance(ctx context.Context, userID string, tx *sql.Tx) (float64, error) {
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, `SELECT balance FROM user_balances WHERE user_id = $1`, userID)
	} else {
		row = DB.QueryRowContext(ctx, `SELECT balance FROM user_balances WHERE user_id = $1`, userID)
	}

	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return balance, nil
}
