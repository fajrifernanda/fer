package repository

import (
	"context"

	"github.com/go-redsync/redsync"
	redcachekeeper "github.com/kumparan/cacher"
	"github.com/kumparan/tapao"

	"github.com/jinzhu/gorm"
	"github.com/kumparan/go-lib/utils"
	log "github.com/sirupsen/logrus"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"
)

// HelloRepository :nodoc:
type HelloRepository interface {
	FindByID(c context.Context, id int64) (*model.Greeting, error)
	Create(c context.Context, g *model.Greeting) error
	Update(c context.Context, id int64, g *model.Greeting) error
}

type helloRepo struct {
	db          *gorm.DB
	cacheKeeper redcachekeeper.Keeper
}

// NewHelloRepository create new repository
func NewHelloRepository(d *gorm.DB, k redcachekeeper.Keeper) HelloRepository {
	return &helloRepo{
		db:          d,
		cacheKeeper: k,
	}
}

// FindByID find object with specific id
func (r *helloRepo) FindByID(ctx context.Context, id int64) (*model.Greeting, error) {
	var greeting *model.Greeting

	cacheKey := model.NewGreetingCacheKeyByID(id)
	greeting, mu, err := r.findFromCacheByID(cacheKey)
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.Dump(ctx),
			"id":      id}).
			Error(err)
		return nil, err
	}

	if greeting != nil {
		return greeting, nil
	}

	var g model.Greeting
	err = r.db.First(&g, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.Dump(ctx),
			"id":      id}).
			Error(err)
	}

	greeting = &g
	newGreetingMsgPack, err := tapao.Marshal(greeting)
	if err != nil {
		log.WithField("msg", utils.Dump(greeting)).Error(err)
		return nil, err
	}

	r.cacheKeeper.Store(mu, redcachekeeper.NewItem(cacheKey, newGreetingMsgPack))
	return greeting, err
}

// Create Greeting
func (r *helloRepo) Create(ctx context.Context, g *model.Greeting) error {
	tx := r.db.Begin()
	err := tx.Create(g).Error
	if err != nil {
		log.WithFields(log.Fields{
			"context":  utils.DumpIncomingContext(ctx),
			"greeting": utils.Dump(g)}).
			Error(err)

		tx.Rollback()
		return err
	}

	r.cacheKeeper.StoreWithoutBlocking(redcachekeeper.NewItem(model.NewGreetingCacheKeyByID(g.ID), utils.ToByte(g)))
	return tx.Commit().Error
}

// Update Greeting by ID
func (r *helloRepo) Update(ctx context.Context, id int64, g *model.Greeting) error {
	tx := r.db.Begin()

	err := tx.Model(g).Omit(g.ImmutableColumns()...).Save(g).Error
	if err != nil {
		log.WithFields(log.Fields{
			"context":  utils.DumpIncomingContext(ctx),
			"greeting": utils.Dump(g)}).
			Error(err)
		tx.Rollback()
		return err
	}

	err = tx.First(&g, g.ID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"id":      g.ID}).
			Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.WithField("tx", utils.Dump(tx)).Error(err)
		return err
	}

	_ = r.cacheKeeper.Purge(model.NewGreetingCacheKeyByID(g.ID))
	return nil
}

func (r *helloRepo) findFromCacheByID(key string) (g *model.Greeting, mu *redsync.Mutex, err error) {
	reply, mu, err := r.cacheKeeper.GetOrLock(key)
	if err != nil {
		return
	}

	if reply == nil {
		return
	}

	g, err = model.NewGreetingFromInterface(reply)

	return
}
