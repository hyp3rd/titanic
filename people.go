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
	Survived              bool      `json:"survived,omitempty"`
	Pclass                int       `json:"pclass,omitempty"`
	Name                  string    `json:"name,omitempty" valid:"length(2|48)"`
	Sex                   string    `json:"sex,omitempty" valid:"required"`
	Age                   int       `json:"age,omitempty" valid:"numeric"`
	SiblingsSpousesAbroad bool      `json:"siblings_spouses_abroad,omitempty"`
	ParentsChildrenAboard bool      `json:"parents_children_aboard,omitempty"`
	Fare                  float32   `json:"fare,omitempty" valid:"float"`
}

// Validate People struct. All the error can be catched with `db.GetErrors()`
func (people People) Validate(db *gorm.DB) {
	if people.Name == "" {
		db.AddError(errors.New("Name is required"))
	}
	if people.Sex == "" {
		db.AddError(errors.New("Sex value can't be blank"))
	}
	if people.Age <= 0 || people.Age >= 110 {
		db.AddError(errors.New("Age must be in a valid range"))
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

//sha256:a491b18f7b611590fc54be38f0af9a25dd5dc530dfcbca9a74823b4f8f3f091c
