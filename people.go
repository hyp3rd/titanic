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
	Survived              bool      `json:"survived,omitempty" valid:"bool"`
	Pclass                int       `json:"pclass,omitempty" valid:"numeric"`
	Name                  string    `json:"name,omitempty" valid:"length(2|48)"`
	Sex                   string    `json:"sex,omitempty" valid:"required"`
	Age                   int       `json:"age,omitempty" valid:"numeric"`
	SiblingsSpousesAbroad int       `json:"siblings_spouses_abroad,omitempty" valid:"numeric"`
	ParentsChildrenAboard int       `json:"parents_children_aboard,omitempty" valid:"numeric"`
	Fare                  float32   `json:"fare,omitempty" valid:"float"`
}

// Validate People struct. All the error can be catched with `db.GetErrors()`
func (people People) Validate(db *gorm.DB) {
	if people.Name == "" {
		db.AddError(errors.New("Name: is required"))
	}
	if !isValidSex(people.Sex) {
		db.AddError(errors.New("Sex: invalid value"))
	}
	if people.Age < 0 || people.Age > 116 {
		db.AddError(errors.New("Age: value in invalid range"))
	}
	if people.SiblingsSpousesAbroad < 0 || people.SiblingsSpousesAbroad > 20 {
		db.AddError(errors.New("SiblingsSpousesAbroad: value in invalid range"))
	}
	if people.ParentsChildrenAboard < 0 || people.ParentsChildrenAboard > 20 {
		db.AddError(errors.New("ParentsChildrenAboard: value in invalid range"))
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

// Validation helpers
func isValidSex(a string) bool {
	list := [...]string{"male", "female", "not declared", ""}

	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
