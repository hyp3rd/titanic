package implementation

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"gitlab.com/hyperd/titanic"
)

// service implements the Titanic Service
type service struct {
	repository titanic.Repository
	logger     log.Logger
}

// NewService creates and returns a new Titanic service instance
func NewService(rep titanic.Repository, logger log.Logger) titanic.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

// Response errors
var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

func (s *service) PostPeople(ctx context.Context, people titanic.People) (string, error) {
	logger := log.With(s.logger, "method", "PostPeople")
	uuid := uuid.New()

	people.ID = uuid

	id, err := s.repository.PostPeople(ctx, people)

	if err != nil {
		level.Error(logger).Log("err", err)
		return id, titanic.ErrCmdRepository
	}
	return id, err
}

func (s *service) GetPeopleByID(ctx context.Context, uuid uuid.UUID) (titanic.People, error) {
	logger := log.With(s.logger, "method", "GetPeopleByID")
	people, err := s.repository.GetPeopleByID(ctx, uuid)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return people, titanic.ErrNotFound
		}
		return people, titanic.ErrQueryRepository
	}
	return people, err
}

func (s *service) PutPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) error {
	logger := log.With(s.logger, "method", "PutPeople")
	if err := s.repository.PutPeople(ctx, uuid, p); err != nil {
		level.Error(logger).Log("err", err)
		return titanic.ErrCmdRepository
	}
	return nil
}

func (s *service) PatchPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) error {
	logger := log.With(s.logger, "method", "PatchPeople")
	if err := s.repository.PatchPeople(ctx, uuid, p); err != nil {
		level.Error(logger).Log("err", err)
		return titanic.ErrCmdRepository
	}
	return nil
}

func (s *service) DeletePeople(ctx context.Context, uuid uuid.UUID) (string, error) {
	logger := log.With(s.logger, "method", "DeletePeople")
	id, err := s.repository.DeletePeople(ctx, uuid)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return uuid.String(), titanic.ErrNotFound
		}
		return uuid.String(), titanic.ErrQueryRepository
	}
	return id, err
}

func (s *service) GetPeople(ctx context.Context) ([]titanic.People, error) {
	logger := log.With(s.logger, "method", "GetPeople")
	people, err := s.repository.GetPeople(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return people, titanic.ErrNotFound
		}
		return people, titanic.ErrQueryRepository
	}
	return people, err
}
