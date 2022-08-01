package managers

import (
	"encoding/json"
	"errors"

	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/internal/models"
)

type PostManager struct {
	Sr      *fcore.RedisClientRing
	IsAdmin bool
}

func NewPostManager(sr *fcore.RedisClientRing, isAdmin bool) *PostManager {
	return &PostManager{
		Sr:      sr,
		IsAdmin: isAdmin,
	}
}

func (dm *PostManager) GenerateUID() string {
	var uid string

	for {
		// Generate New One
		guid := fcore.GenerateRandUID()

		// Check guid is really unique
		found, _ := dm.Has(guid)
		if !found {
			uid = guid
			break
		}
	}

	return uid
}

func (dm *PostManager) Count() (int, error) {
	return dm.Sr.Clients[fcore.LiveStore].Do(
		dm.Sr.Ctx,
		"stats",
		fcore.ItemsBucket,
	).Int()
}

func (dm *PostManager) Has(uid string) (bool, error) {
	return dm.Sr.Clients[fcore.LiveStore].Do(
		dm.Sr.Ctx,
		"exists",
		fcore.ItemsBucket,
		uid,
	).Bool()
}

func (dm *PostManager) Get(uid string, v interface{}) (interface{}, error) {
	result, err := dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "get", fcore.ItemsBucket, uid).Text()
	if err != nil {
		return nil, err
	}

	if result != "" {
		err = json.Unmarshal([]byte(result), v)
		return v, err
	}

	return nil, err
}

func (dm *PostManager) Set(uid string, v interface{}) (string, error) {
	return dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "set", fcore.ItemsBucket, uid, v).Text()
}

func (dm *PostManager) Del(uid string) error {
	count, err := dm.Sr.Clients[fcore.LiveStore].Del(dm.Sr.Ctx, fcore.ItemsBucket, uid).Result()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("key not delete")
	}

	return nil
}

func (dm *PostManager) List(lastId string, limit int) (map[int]models.Post, error) {
	result, err := dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "list", fcore.ItemsBucket, lastId, limit).StringSlice()
	if err != nil {
		return nil, err
	}

	return dm.GetResultFromJson(result), err
}

func (dm *PostManager) PrevList(lastId string, limit int) (map[int]models.Post, error) {
	result, err := dm.Sr.Clients[fcore.LiveStore].Do(dm.Sr.Ctx, "prevlist", fcore.ItemsBucket, lastId, limit).StringSlice()
	if err != nil {
		return nil, err
	}

	return dm.GetResultFromJson(result), err
}

func (dm *PostManager) GetResultFromJson(result []string) map[int]models.Post {
	items := map[int]models.Post{}
	if len(result) > 0 {
		for in, item := range result {
			postModel := models.Post{}
			_ = json.Unmarshal([]byte(item), &postModel)

			if !dm.IsAdmin && postModel.Status != fcore.StatusActive {
				continue
			}

			items[in] = postModel
		}
	}

	return items
}

func (dm *PostManager) GetResultFromUIDS(uids []string) map[int]models.Post {
	items := map[int]models.Post{}
	if len(uids) > 0 {
		for in, uid := range uids {
			postModel := models.Post{}
			_, _ = dm.Get(uid, &postModel)

			if !dm.IsAdmin && postModel.Status != fcore.StatusActive {
				continue
			}

			items[in] = postModel
		}
	}

	return items
}
