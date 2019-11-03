package titanic

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// People represents a single passenger (People).
// UUID should be globally unique.
type People struct {
	// gorm.Model
	ID                    uuid.UUID `json:"uuid,omitempty" gorm:"primary_key"`
	Survived              *bool     `json:"survived,omitempty"`
	Pclass                *int      `json:"pclass,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Sex                   string    `json:"sex,omitempty"`
	Age                   int       `json:"age,omitempty"`
	SiblingsSpousesAbroad *bool     `json:"siblings_spouses_abroad,omitempty"`
	ParentsChildrenAboard *bool     `json:"parents_children_aboard,omitempty"`
	Fare                  *float32  `json:"fare,omitempty"`
}

// Validate People struct. All the error can be catched with `db.GetErrors()`
func (people People) Validate(db *gorm.DB) {
	if people.Age >= 18 {
		db.AddError(errors.New("Age need to be 18+"))
	}
	if people.Name == "" {
		db.AddError(errors.New("Name can't be blank"))
	}
}

// Repository describes the persistence on people model
type Repository interface {
	PostPeople(ctx context.Context, p People) (string, error)
	GetPeopleByID(ctx context.Context, ID uuid.UUID) (People, error)
	PutPeople(ctx context.Context, ID uuid.UUID, p People) error
	PatchPeople(ctx context.Context, ID uuid.UUID, p People) error
	DeletePeople(ctx context.Context, ID uuid.UUID) (string, error)
	GetPeople(ctx context.Context) ([]People, error)
}
