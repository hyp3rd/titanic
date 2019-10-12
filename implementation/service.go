package implementation

import (
	"context"
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
	ErrInconsistentUUIDs = errors.New("inconsistent UUIDs")
	ErrAlreadyExists     = errors.New("already exists")
	ErrNotFound          = errors.New("not found")
)

func (s *service) PostPeople(ctx context.Context, people titanic.People) (string, error) {
	logger := log.With(s.logger, "method", "PostPeople")
	uuid := uuid.New()

	people.UUID = uuid

	if err := s.repository.PostPeople(ctx, people); err != nil {
		level.Error(logger).Log("err", err)
		return "", titanic.ErrCmdRepository
	}
	return uuid.String(), nil
}

func (s *service) GetPeople(ctx context.Context, uuid uuid.UUID) (titanic.People, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[uuid.String()]
	if !ok {
		return titanic.People{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PutPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) error {
	if p.UUID.String() == "" {
		return ErrInconsistentUUIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	p.UUID = uuid
	s.m[uuid.String()] = p // PUT = create or update
	return nil
}

func (s *inmemService) PatchPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) error {
	if p.UUID.String() == "" {
		return ErrInconsistentUUIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[uuid.String()]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}

	// It should not possible to PATCH the UUID, and it should not be
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the People definition.

	if p.Survived != nil {
		existing.Survived = p.Survived
	}
	if p.Pclass != nil {
		existing.Pclass = p.Pclass
	}
	if p.Name != "" {
		existing.Name = p.Name
	}
	if p.Sex != "" {
		existing.Sex = p.Sex
	}
	if p.Age != nil {
		existing.Age = p.Age
	}
	if p.SiblingsSpousesAbroad != nil {
		existing.SiblingsSpousesAbroad = p.SiblingsSpousesAbroad
	}
	if p.ParentsChildrenAboard != nil {
		existing.ParentsChildrenAboard = p.ParentsChildrenAboard
	}
	if p.Fare != nil {
		existing.Fare = p.Fare
	}

	s.m[uuid.String()] = existing
	return nil
}

func (s *inmemService) DeletePeople(ctx context.Context, uuid uuid.UUID) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[uuid.String()]; !ok {
		return ErrNotFound
	}
	delete(s.m, uuid.String())
	return nil
}

func (s *inmemService) GetAllPeople(ctx context.Context) ([]titanic.People, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	p := []titanic.People{}
	for _, value := range s.m {
		p = append(p, value)
	}

	if p == nil {
		return []titanic.People{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) GetAPIStatus(ctx context.Context) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	status := "Healthy"

	return status, nil
}
