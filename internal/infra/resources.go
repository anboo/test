package infra

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Resources struct {
	Env    Env
	DB     *gorm.DB
	Logger *slog.Logger
}

func Init(ctx context.Context) (*Resources, error) {
	r := &Resources{}

	errGrp, ctx := errgroup.WithContext(ctx)

	err := r.initEnv()
	if err != nil {
		return nil, err
	}

	r.initLogger()

	errGrp.Go(func() error {
		r.Logger.Info("starting db connection")
		defer r.Logger.Info("done db connection")
		return r.initDb(ctx)
	})

	return r, errGrp.Wait()
}
