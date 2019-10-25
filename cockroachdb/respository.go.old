package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"gitlab.com/hyperd/titanic"
)

// var (
// 	ErrRepository = errors.New("unable to handle request")
// )

type repository struct {
	db     *sql.DB
	logger log.Logger
}

// New returns a concrete repository backed by CockroachDB
func New(db *sql.DB, logger log.Logger) (titanic.Repository, error) {
	// return  repository
	return &repository{
		db:     db,
		logger: log.With(logger, "rep", "cockroachdb"),
	}, nil
}

func (repo *repository) PostPeople(ctx context.Context, people titanic.People) (string, error) {
	// Run a transaction to sync the query model.
	err := crdb.ExecuteTx(ctx, repo.db, nil, func(tx *sql.Tx) error {
		_, err := createPeople(tx, people)
		return err
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

func (repo *repository) GetPeopleByID(ctx context.Context, id uuid.UUID) (titanic.People, error) {
	var peopleRow = titanic.People{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT uuid, survived, pclass, name, sex, age, siblingsspousesabroad, parentschildrenaboard, fare FROM people WHERE uuid = $1",
		id).
		Scan(
			&peopleRow.UUID, &peopleRow.Survived, &peopleRow.Pclass, &peopleRow.Name, &peopleRow.Sex, &peopleRow.Age, &peopleRow.SiblingsSpousesAbroad, &peopleRow.ParentsChildrenAboard, &peopleRow.Fare,
		); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return peopleRow, err
	}

	return peopleRow, nil
}

func (repo *repository) PutPeople(ctx context.Context, id uuid.UUID, people titanic.People) error {
	return nil
}

func (repo *repository) PatchPeople(ctx context.Context, id uuid.UUID, people titanic.People) error {
	return nil
}

func (repo *repository) DeletePeople(ctx context.Context, id uuid.UUID) (string, error) {
	return "", nil
}

func (repo *repository) GetPeople(ctx context.Context) ([]titanic.People, error) {
	return []titanic.People{}, nil
}

// Private helpers
func createPeople(tx *sql.Tx, people titanic.People) (string, error) {
	// Insert a passenger into the "people" table.
	sql := `
			INSERT INTO people (uuid, survived, pclass, name, sex, age, siblingsspousesabroad, parentschildrenaboard, fare)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := tx.Exec(sql, people.UUID, people.Survived, people.Pclass, people.Name, people.Sex, people.Age, people.SiblingsSpousesAbroad, people.ParentsChildrenAboard, people.Fare)

	if err != nil {
		return people.UUID.String(), err
	}

	return people.UUID.String(), nil
}

// type People struct {
// 	UUID                  uuid.UUID `json:"uuid,omitempty"`
// 	Survived              *bool     `json:"survived,omitempty"`
// 	Pclass                *int      `json:"pclass,omitempty"`
// 	Name                  string    `json:"name,omitempty"`
// 	Sex                   string    `json:"sex,omitempty"`
// 	Age                   *int      `json:"age,omitempty"`
// 	SiblingsSpousesAbroad *bool     `json:"siblings_spouses_abroad,omitempty"`
// 	ParentsChildrenAboard *bool     `json:"parents_children_aboard,omitempty"`
// 	Fare                  *float32  `json:"fare,omitempty"`
// }

// func (repo *repository) ChangeOrderStatus(ctx context.Context, orderId string, status string) error {
// 	sql := `
// UPDATE orders
// SET status=$2
// WHERE id=$1`

// 	_, err := repo.db.ExecContext(ctx, sql, orderId, status)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
