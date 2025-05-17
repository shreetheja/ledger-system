package kafka

import "time"

// Topic constants
const (
	TopicCreateAccount = "create-account"
	TopicAddBalance    = "add-balance"
	TopicDeductBalance = "deduct-balance"
)

// Base struct for all Kafka messages
type BaseMessage struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateAccountMessage represents a new user account creation event
type CreateAccountMessage struct {
	BaseMessage
	InitialBalance float64 `json:"initial_balance"`
}

// AddBalanceMessage represents an event where funds are added to a user's account
type AddBalanceMessage struct {
	BaseMessage
	Amount float64 `json:"amount"`
}

// DeductBalanceMessage represents an event where funds are deducted from a user's account
type DeductBalanceMessage struct {
	BaseMessage
	Amount float64 `json:"amount"`
}

var Topics = []string{
	TopicCreateAccount,
	TopicAddBalance,
	TopicDeductBalance,
}
