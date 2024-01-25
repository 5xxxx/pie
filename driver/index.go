package driver

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Indexes represents the interface for managing indexes in a database.
type Indexes interface {
	CreateIndexes(doc interface{}, ctx ...context.Context) ([]string, error)
	DropAll(doc interface{}, ctx ...context.Context) error
	DropOne(doc interface{}, name string, ctx ...context.Context) error
	AddIndex(keys interface{}, opt ...*options.IndexOptions) Indexes
	SetMaxTime(d time.Duration) Indexes
	SetCommitQuorumInt(quorum int32) Indexes
	SetCommitQuorumString(quorum string) Indexes
	SetCommitQuorumMajority() Indexes
	SetCommitQuorumVotingMembers() Indexes
	SetDatabase(db string) Indexes
	Collection(doc interface{}) Indexes
}
