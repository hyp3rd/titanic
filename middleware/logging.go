package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"

	"gitlab.com/hyperd/titanic"
)

// LoggingMiddleware provides basic logging Middleware
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next titanic.Service) titanic.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   titanic.Service
	logger log.Logger
}

func (mw loggingMiddleware) PostPeople(ctx context.Context, p titanic.People) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostPeople", "uuid", p.UUID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostPeople(ctx, p)
}

func (mw loggingMiddleware) GetPeopleByID(ctx context.Context, uuid uuid.UUID) (p titanic.People, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetPeople", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetPeopleByID(ctx, uuid)
}

func (mw loggingMiddleware) PutPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutPeople", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutPeople(ctx, uuid, p)
}

func (mw loggingMiddleware) PatchPeople(ctx context.Context, uuid uuid.UUID, p titanic.People) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchPeople", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchPeople(ctx, uuid, p)
}

func (mw loggingMiddleware) DeletePeople(ctx context.Context, uuid uuid.UUID) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeletePeople", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeletePeople(ctx, uuid)
}

func (mw loggingMiddleware) GetPeople(ctx context.Context) (allPeople []titanic.People, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAllPeople", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetPeople(ctx)
}

func (mw loggingMiddleware) GetAPIStatus(ctx context.Context) (status string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAPIStatus", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAPIStatus(ctx)
}
