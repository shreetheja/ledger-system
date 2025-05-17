package service

import (
	"context"
	"ledger/kafka"
	"ledger/pg"
	"log"
	"time"
)

func GetUserBalance(userID string) (float64, error) {
	balance, err := pg.GetBalance(context.Background(), userID, nil)
	if err != nil {
		log.Printf("Error getting balance for user %s: %v", userID, err)
		return 0.0, err
	}
	return balance, nil
}

func AddAmount(userID string, amount float64) error {
	err := kafka.SendAddBalanceMessage(kafka.AddBalanceMessage{
		Amount: amount,
		BaseMessage: kafka.BaseMessage{
			UserID:    userID,
			Timestamp: time.Now(),
		},
	})
	if err != nil {
		log.Printf("Error creating account for user %s: %v", userID, err)
		return err
	}
	return nil
}

func DeductAmount(userID string, amount float64) error {
	err := kafka.SendDeductBalanceMessage(kafka.DeductBalanceMessage{
		Amount: amount,
		BaseMessage: kafka.BaseMessage{
			UserID:    userID,
			Timestamp: time.Now(),
		},
	})
	if err != nil {
		log.Printf("Error creating account for user %s: %v", userID, err)
		return err
	}
	return nil
}

func CreateAccount(userID string, initialBalance float64) error {
	err := kafka.SendCreateAccountMessage(kafka.CreateAccountMessage{
		InitialBalance: initialBalance,
		BaseMessage: kafka.BaseMessage{
			UserID:    userID,
			Timestamp: time.Now(),
		},
	})
	if err != nil {
		log.Printf("Error creating account for user %s: %v", userID, err)
		return err
	}
	return nil
}
