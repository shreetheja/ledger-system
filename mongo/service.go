package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Pass a slice of LedgerRecord, so you can record multiple ops atomically.
func RecordTransaction(ctx context.Context, records []LedgerRecord) error {
	session, err := MongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongo session: %w", err)
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, rec := range records {
			rec.Timestamp = time.Now().UTC() // Set timestamp if not set

			_, err := LedgerCollection.InsertOne(sessCtx, rec)
			if err != nil {
				return nil, fmt.Errorf("failed to insert ledger record: %w", err)
			}
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return fmt.Errorf("mongo transaction failed: %w", err)
	}

	return nil
}

func GetUserLogs(userID string) ([]LedgerRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := LedgerCollection.Find(ctx, map[string]string{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find user logs: %w", err)
	}
	defer cursor.Close(ctx)

	var records []LedgerRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode user logs: %w", err)
	}

	return records, nil
}
