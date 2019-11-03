package cockroachdb

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/hyperd/titanic"
)

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

// New returns a concrete repository backed by CockroachDB
func New(db *gorm.DB, logger log.Logger) (titanic.Repository, error) {
	// return  repository
	return &repository{
		db:     db,
		logger: log.With(logger, "rep", "cockroachdb"),
	}, nil
}

// Functions of type `txnFunc` are passed as arguments to our
// `runTransaction` wrapper that handles transaction retries for us
// (see implementation below).
type txnFunc func(*gorm.DB) error

// This function is used for testing the transaction retry loop.  It
// can be deleted from production code.
var forceRetryLoop txnFunc = func(db *gorm.DB) error {

	// The first statement in a transaction can be retried transparently
	// on the server, so we need to add a dummy statement so that our
	// force_retry statement isn't the first one.
	if err := db.Exec("SELECT now()").Error; err != nil {
		return err
	}
	// Used to force a transaction retry.  Can only be run as the
	// 'root' user.
	if err := db.Exec("SELECT crdb_internal.force_retry('1s'::INTERVAL)").Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) PostPeople(ctx context.Context, people titanic.People) (string, error) {
	// Run a transaction to sync the query model.
	id := uuid.New()
	people.ID = id

	repo.db.Create(&titanic.People{
		ID:                    people.ID,
		Survived:              people.Survived,
		Pclass:                people.Pclass,
		Name:                  people.Name,
		Sex:                   people.Sex,
		Age:                   people.Age,
		SiblingsSpousesAbroad: people.SiblingsSpousesAbroad,
		ParentsChildrenAboard: people.ParentsChildrenAboard,
		Fare:                  people.Fare})

	return id.String(), nil
}

func (repo *repository) GetPeopleByID(ctx context.Context, id uuid.UUID) (titanic.People, error) {
	var people = titanic.People{}

	repo.db.Where("id = ?", id).First(&people)

	return people, nil
}

func (repo *repository) PutPeople(ctx context.Context, id uuid.UUID, people titanic.People) error {
	// Update multiple attributes with `struct`, will only update those changed & non blank fields
	repo.db.Model(&people).Updates(titanic.People{
		ID:                    id,
		Survived:              people.Survived,
		Pclass:                people.Pclass,
		Name:                  people.Name,
		Sex:                   people.Sex,
		Age:                   people.Age,
		SiblingsSpousesAbroad: people.SiblingsSpousesAbroad,
		ParentsChildrenAboard: people.ParentsChildrenAboard,
		Fare:                  people.Fare,
	})

	return nil
}

func (repo *repository) PatchPeople(ctx context.Context, id uuid.UUID, people titanic.People) error {
	// Update multiple attributes with `map`, will only update those changed fields
	repo.db.Model(&people).Updates(map[string]interface{}{
		"Survived":              people.Survived,
		"Pclass":                people.Pclass,
		"Name":                  people.Name,
		"Sex":                   people.Sex,
		"Age":                   people.Age,
		"SiblingsSpousesAbroad": people.SiblingsSpousesAbroad,
		"ParentsChildrenAboard": people.ParentsChildrenAboard,
		"Fare":                  people.Fare,
	})
	return nil
}

func (repo *repository) DeletePeople(ctx context.Context, id uuid.UUID) (string, error) {

	repo.db.Where("id = ?", id).Delete(&titanic.People{})

	return id.String(), nil
}

func (repo *repository) GetPeople(ctx context.Context) ([]titanic.People, error) {
	var people []titanic.People

	repo.db.Find(&people)

	return people, nil
}
