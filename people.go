package titanic

import (
	"context"

	"github.com/google/uuid"
)

// People represents a single passenger (People).
// UUID should be globally unique.
type People struct {
	UUID                  uuid.UUID `json:"uuid,omitempty"`
	Survived              *bool     `json:"survived,omitempty"`
	Pclass                *int      `json:"pclass,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Sex                   string    `json:"sex,omitempty"`
	Age                   *int      `json:"age,omitempty"`
	SiblingsSpousesAbroad *bool     `json:"siblings_spouses_abroad,omitempty"`
	ParentsChildrenAboard *bool     `json:"parents_children_aboard,omitempty"`
	Fare                  *float32  `json:"fare,omitempty"`
}

// Repository describes the persistence on people model
type Repository interface {
	PostPeople(ctx context.Context, p People) error
	GetPeople(ctx context.Context, uuid uuid.UUID) (People, error)
	PutPeople(ctx context.Context, uuid uuid.UUID, p People) error
	PatchPeople(ctx context.Context, uuid uuid.UUID, p People) error
	DeletePeople(ctx context.Context, uuid uuid.UUID) error
	GetAllPeople(ctx context.Context) ([]People, error)
}
