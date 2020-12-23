package driver

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Indexes interface {
	CreateIndexes(ctx context.Context, doc interface{}) ([]string, error)

	DropAll(ctx context.Context, doc interface{}) error

	DropOne(ctx context.Context, doc interface{}, name string) error

	AddIndex(keys interface{}, opt ...*options.IndexOptions) Indexes

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) Indexes

	// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
	SetCommitQuorumInt(quorum int32) Indexes

	// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
	SetCommitQuorumString(quorum string) Indexes

	// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
	SetCommitQuorumMajority() Indexes

	// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
	SetCommitQuorumVotingMembers() Indexes

	SetDatabase(db string) Indexes

	Collection(doc interface{}) Indexes
}
