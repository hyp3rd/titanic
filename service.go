package titanic

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Response errors
var (
	ErrInconsistentUUIDs = errors.New("inconsistent UUIDs")
	ErrAlreadyExists     = errors.New("already exists")
	ErrNotFound          = errors.New("not found")
	ErrCmdRepository     = errors.New("unable to command repository")
	ErrQueryRepository   = errors.New("unable to query repository")
)

// Service is a CRUD interface for People in the Titanic collection.
type Service interface {
	PostPeople(ctx context.Context, p People) (string, error)
	GetPeopleByID(ctx context.Context, uuid uuid.UUID) (People, error)
	PutPeople(ctx context.Context, uuid uuid.UUID, p People) error
	PatchPeople(ctx context.Context, uuid uuid.UUID, p People) error
	DeletePeople(ctx context.Context, uuid uuid.UUID) (string, error)
	GetPeople(ctx context.Context) ([]People, error)
}
