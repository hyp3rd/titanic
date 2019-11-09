package titanic

import (
	"context"

	"github.com/google/uuid"
)

// People represents a single passenger (People).
// UUID should be globally unique.
type People struct {
	// gorm.Model
	ID                    uuid.UUID `json:"uuid,omitempty" gorm:"primary_key" valid:"uuidv4"`
	Survived              bool      `json:"survived,omitempty"`
	Pclass                int       `json:"pclass,omitempty" valid:"numeric"`
	Name                  string    `json:"name,omitempty" valid:"alphanum,stringlength(2|70)"` // https://webarchive.nationalarchives.gov.uk/20100407173424/http://www.cabinetoffice.gov.uk/govtalk/schemasstandards/e-gif/datastandards.aspx
	Sex                   string    `json:"sex,omitempty" valid:"in(male|female|not declared)"`
	Age                   int       `json:"age,omitempty" valid:"numeric,range(0|116)"`
	SiblingsSpousesAbroad int       `json:"siblings_spouses_abroad,omitempty" valid:"numeric,range(0|20)"`
	ParentsChildrenAboard int       `json:"parents_children_aboard,omitempty" valid:"numeric,range(0|20)"`
	Fare                  float32   `json:"fare,omitempty" valid:"float"`
}

// sha256:9f40d23ad125b9f6aa5a706b3f6d3e6b78f5937813812f2ca7e62769055104f3
// Validate People struct. All the error can be catched with `db.GetErrors()`
// func (people People) Validate(db *gorm.DB) {
// 	nameIsCorrect, _ := regexp.MatchString(`^[A-Za-z0-9 _]*[A-Za-z0-9][A-Za-z0-9 _]*$`, people.Name)
// 	if !nameIsCorrect {
// 		db.AddError(errors.New("Name: invalid format"))
// 	}
// 	if !isValidSex(people.Sex) {
// 		db.AddError(errors.New("Sex: invalid value"))
// 	}
// 	if people.Age < 0 || people.Age > 116 {
// 		db.AddError(errors.New("Age: value in invalid range"))
// 	}
// 	if people.SiblingsSpousesAbroad < 0 || people.SiblingsSpousesAbroad > 20 {
// 		db.AddError(errors.New("SiblingsSpousesAbroad: value in invalid range"))
// 	}
// 	if people.ParentsChildrenAboard < 0 || people.ParentsChildrenAboard > 20 {
// 		db.AddError(errors.New("ParentsChildrenAboard: value in invalid range"))
// 	}
// }

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
// func isValidSex(a string) bool {
// 	list := [...]string{"male", "female", "not declared", ""}

// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }
