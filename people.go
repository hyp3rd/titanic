package titanic

import (
	"context"

	"github.com/google/uuid"
)

// People represents a single passenger (People).
type People struct {
	// gorm.Model
	ID                    uuid.UUID `json:"uuid,omitempty" gorm:"primary_key"`
	Survived              bool      `json:"survived,omitempty"`
	Pclass                int       `json:"pclass,omitempty" valid:"numeric"`
	Name                  string    `json:"name,omitempty" valid:"alphanum,stringlength(2|70)"` // https://webarchive.nationalarchives.gov.uk/20100407173424/http://www.cabinetoffice.gov.uk/govtalk/schemasstandards/e-gif/datastandards.aspx
	Sex                   string    `json:"sex,omitempty" valid:"in(male|female|not declared)"`
	Age                   int       `json:"age,omitempty" valid:"numeric,range(0|116)"`
	SiblingsSpousesAbroad int       `json:"siblings_spouses_abroad,omitempty" valid:"numeric,range(0|20)"`
	ParentsChildrenAboard int       `json:"parents_children_aboard,omitempty" valid:"numeric,range(0|20)"`
	Fare                  float32   `json:"fare,omitempty" valid:"float"`
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
