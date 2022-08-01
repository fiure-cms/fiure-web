package managers

import (
	"encoding/json"
	"errors"

	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/internal/models"
)

type PageManager struct {
	Sr      *fcore.RedisClientRing
	IsAdmin bool
}

func NewPageManager(sr *fcore.RedisClientRing, isAdmin bool) *PageManager {
	return &PageManager{
		Sr:      sr,
		IsAdmin: isAdmin,
	}
}

func (dm *PageManager) Count() (int, error) {
	return dm.Sr.Clients[fcore.LiveStore].Do(
		dm.Sr.Ctx,
		"stats",
		fcore.PagesBucket,
	).Int()
}

func (dm *PageManager) Has(uid string) (bool, error) {
	return dm.Sr.Clients[fcore.LiveStore].Do(
		dm.Sr.Ctx,
		"exists",
		fcore.PagesBucket,
		uid,
	).Bool()
}

func (dm *PageManager) Get(slug string, v interface{}) (interface{}, error) {
	result, err := dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "get", fcore.PagesBucket, slug).Text()
	if err != nil {
		return nil, err
	}

	if result != "" {
		err = json.Unmarshal([]byte(result), v)

		return v, err
	}

	return nil, err
}

func (dm *PageManager) Set(slug string, v interface{}) (string, error) {
	return dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "set", fcore.PagesBucket, slug, v).Text()
}

func (dm *PageManager) Del(slug string) error {
	count, err := dm.Sr.Clients[fcore.LiveStore].Del(dm.Sr.Ctx, fcore.PagesBucket, slug).Result()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("key not delete")
	}

	return nil
}

func (dm *PageManager) List(lastId string, limit int) (map[int]models.Page, error) {
	result, err := dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "list", fcore.PagesBucket, lastId, limit).StringSlice()
	if err != nil {
		return nil, err
	}

	items := map[int]models.Page{}
	if len(result) > 0 {
		for in, item := range result {
			page := models.Page{}
			_ = json.Unmarshal([]byte(item), &page)

			if !dm.IsAdmin && page.Status != fcore.StatusActive {
				continue
			}

			items[in] = page
		}
	}

	return items, err
}
