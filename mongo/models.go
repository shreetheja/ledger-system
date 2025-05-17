package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LedgerRecord struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	Operation     string             `bson:"operation"`      // e.g., "CreateAccount", "AddBalance", "DeductBalance"
	Amount        float64            `bson:"amount"`         // positive or negative
	Timestamp     time.Time          `bson:"timestamp"`      // when transaction happened
	TransactionID string             `bson:"transaction_id"` // optional to correlate multiple ops in one transaction
}
