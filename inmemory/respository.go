package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"gitlab.com/hyperd/titanic"
)

// Response errors
var (
	ErrInconsistentID = errors.New("inconsistent ID")
	ErrAlreadyExists  = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
)

type repository struct {
	mtx    sync.RWMutex
	m      map[string]titanic.People
	logger log.Logger
}

// NewInmemService returns an in-memory storage
func NewInmemService(logger log.Logger) (titanic.Repository, error) {
	return &repository{
		m:      map[string]titanic.People{},
		logger: log.With(logger, "repository", "inmemory"),
	}, nil
}

func (r *repository) PostPeople(ctx context.Context, p titanic.People) (string, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	// id, err := uuid.Parse(vars["uuid"])
	id := uuid.New()

	p.ID = id

	if _, ok := r.m[p.ID.String()]; ok {
		return "", ErrAlreadyExists // POST = create, don't overwrite
	}
	r.m[p.ID.String()] = p
	return id.String(), nil
}

func (r *repository) GetPeopleByID(ctx context.Context, uuid uuid.UUID) (titanic.People, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	p, ok := r.m[uuid.String()]
	if !ok {
		return titanic.People{}, ErrNotFound
	}
	return p, nil
}

func (r *repository) PutPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) error {
	if p.ID.String() == "" {
		return ErrInconsistentID
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	existing, ok := r.m[uuid.String()]
	if !ok {
		p.ID = uuid // PUT can create
	}

	r.m[uuid.String()] = setPeople(p, existing)

	return nil
}

func (r *repository) PatchPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) error {
	if p.ID.String() == "" {
		return ErrInconsistentID
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	existing, ok := r.m[uuid.String()]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}

	r.m[uuid.String()] = setPeople(p, existing)
	return nil
}

func (r *repository) DeletePeople(ctx context.Context, uuid uuid.UUID) (string, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.m[uuid.String()]; !ok {
		return uuid.String(), ErrNotFound
	}
	delete(r.m, uuid.String())
	return uuid.String(), nil
}

func (r *repository) GetPeople(ctx context.Context) ([]titanic.People, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	p := []titanic.People{}
	for _, value := range r.m {
		p = append(p, value)
	}

	if p == nil {
		return []titanic.People{}, ErrNotFound
	}
	return p, nil
}

func setPeople(p titanic.People, existing titanic.People) titanic.People {

	// It should not possible to PATCH the ID, and it should not be
	// possible to PATCH or PUT any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Survived *bool in
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

	return existing
}
