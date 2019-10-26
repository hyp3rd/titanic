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
	people.UUID = id

	repo.db.Create(&titanic.People{
		UUID:                  people.UUID,
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

	repo.db.Where("uuid = ?", id).First(&people)

	return people, nil
}

func (repo *repository) PutPeople(ctx context.Context, id uuid.UUID, people titanic.People) error {
	return nil
}

func (repo *repository) PatchPeople(ctx context.Context, id uuid.UUID, people titanic.People) error {
	return nil
}

func (repo *repository) DeletePeople(ctx context.Context, id uuid.UUID) (string, error) {

	repo.db.Where("uuid = ?", id).Delete(&titanic.People{})

	return id.String(), nil
}

func (repo *repository) GetPeople(ctx context.Context) ([]titanic.People, error) {
	var people []titanic.People

	repo.db.Find(&people)

	return people, nil
}

// func transferFunds(db *gorm.DB, fromID int, toID int, amount int) error {
//     var fromAccount Account
//     var toAccount Account

//     db.First(&fromAccount, fromID)
//     db.First(&toAccount, toID)

//     if fromAccount.Balance < amount {
//         return fmt.Errorf("account %d balance %d is lower than transfer amount %d", fromAccount.ID, fromAccount.Balance, amount)
//     }

//     fromAccount.Balance -= amount
//     toAccount.Balance += amount

//     if err := db.Save(&fromAccount).Error; err != nil {
//         return err
//     }
//     if err := db.Save(&toAccount).Error; err != nil {
//         return err
//     }
//     return nil
// }