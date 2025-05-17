package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ledger/kafka"
	"ledger/mongo"
	"ledger/pg"
	"log"

	"github.com/google/uuid"
)

func Initialize(ctx context.Context) {
	// Create consumer handler for topics
	handler := kafka.NewConsumerHandler([]string{
		kafka.TopicCreateAccount,
		kafka.TopicAddBalance,
		kafka.TopicDeductBalance,
	})
	handler.StartConsuming(ctx)

	go func() {
		for msg := range handler.MessageChannel {
			log.Printf("Received message on topic %s: %s\n", *msg.TopicPartition.Topic, msg.Value)
			switch *msg.TopicPartition.Topic {
			case kafka.TopicAddBalance:
				var addBalanceMsg kafka.AddBalanceMessage
				if err := json.Unmarshal(msg.Value, &addBalanceMsg); err != nil {
					log.Printf("Failed to unmarshal add-balance message: %v\n", err)
					continue
				}
				err := HandleAddBalance(addBalanceMsg)
				if err != nil {
					log.Printf("Failed to handle add-balance message: %v\n", err)
					continue
				}
				log.Printf("User %s added balance: %f\n", addBalanceMsg.UserID, addBalanceMsg.Amount)
			case kafka.TopicDeductBalance:
				var deductBalanceMsg kafka.DeductBalanceMessage
				if err := json.Unmarshal(msg.Value, &deductBalanceMsg); err != nil {
					log.Printf("Failed to unmarshal deduct-balance message: %v\n", err)
					continue
				}
				err := HandleDeductBalance(deductBalanceMsg)
				if err != nil {
					log.Printf("Failed to handle deduct-balance message: %v\n", err)
					continue
				}
				log.Printf("User %s deducted balance: %f\n", deductBalanceMsg.UserID, deductBalanceMsg.Amount)
			case kafka.TopicCreateAccount:
				var createAccountMsg kafka.CreateAccountMessage
				if err := json.Unmarshal(msg.Value, &createAccountMsg); err != nil {

					log.Printf("Failed to unmarshal create-account message: %v\n", err)
					continue
				}
				err := HandleCreateAccount(createAccountMsg)
				if err != nil {
					log.Printf("Failed to handle deduct-balance message: %v\n", err)
					continue
				}
				log.Printf("User %s created account with initial balance: %f\n", createAccountMsg.UserID, createAccountMsg.InitialBalance)
			default:
				log.Printf("Unknown topic: %s\n", *msg.TopicPartition.Topic)
			}
		}
	}()
}

func HandleCreateAccount(msg kafka.CreateAccountMessage) error {
	ctx := context.Background()

	tx, err := pg.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return err
	}

	// Update balance synchronously inside transaction
	err = pg.CreateNewAccount(ctx, msg.UserID, float64(msg.InitialBalance), tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Commit the transaction before recording in Mongo
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	records := []mongo.LedgerRecord{
		{
			UserID:        msg.UserID,
			Operation:     "CreateAccount",
			Amount:        msg.InitialBalance,
			TransactionID: uuid.New().String(),
		},
	}

	err = mongo.RecordTransaction(ctx, records)
	if err != nil {
		return fmt.Errorf("failed to record ledger transaction: %w", err)
	}

	log.Printf("Handled account creation for user %s with initial balance %f\n", msg.UserID, msg.InitialBalance)
	return nil
}

func HandleAddBalance(msg kafka.AddBalanceMessage) error {
	ctx := context.Background()

	tx, err := pg.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = pg.UpdateBalance(ctx, msg.UserID, float64(msg.Amount), tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	records := []mongo.LedgerRecord{
		{
			UserID:        msg.UserID,
			Operation:     "AddBalance",
			Amount:        float64(msg.Amount),
			TransactionID: uuid.New().String(),
		},
	}

	if err := mongo.RecordTransaction(ctx, records); err != nil {
		return fmt.Errorf("failed to record ledger transaction: %w", err)
	}

	log.Printf("Handled balance addition for user %s with amount %f\n", msg.UserID, msg.Amount)
	return nil
}

func HandleDeductBalance(msg kafka.DeductBalanceMessage) error {
	ctx := context.Background()

	tx, err := pg.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = pg.UpdateBalance(ctx, msg.UserID, float64(-msg.Amount), tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	records := []mongo.LedgerRecord{
		{
			UserID:        msg.UserID,
			Operation:     "DeductBalance",
			Amount:        float64(-msg.Amount),
			TransactionID: uuid.New().String(),
		},
	}

	if err := mongo.RecordTransaction(ctx, records); err != nil {
		return fmt.Errorf("failed to record ledger transaction: %w", err)
	}

	log.Printf("Handled balance deduction for user %s with amount %f\n", msg.UserID, msg.Amount)
	return nil
}
