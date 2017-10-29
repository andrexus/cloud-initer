package model

import (
	"time"

	"encoding/json"

	"github.com/boltdb/bolt"
)

var environmentBucket = []byte("environment")
var environmentKey = []byte("base-env")

type EnvironmentRepository interface {
	Get() (*Environment, error)
	Save(item *Environment) (*Environment, error)
}

type BoltEnvironmentRepository struct {
	db *bolt.DB
}

func NewEnvironmentRepository(db *bolt.DB) *BoltEnvironmentRepository {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(environmentBucket)
		return err
	})
	return &BoltEnvironmentRepository{db}
}

func (r *BoltEnvironmentRepository) Get() (*Environment, error) {
	var item *Environment
	err := r.db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket(environmentBucket)
		k := environmentKey
		itemData := b.Get(k)
		if len(itemData) == 0 {
			item = &Environment{UpdatedAt: time.Now()}
			return nil
		}
		item, err = decodeEnvironment(itemData)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *BoltEnvironmentRepository) Save(item *Environment) (*Environment, error) {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(environmentBucket)
		{
			item.UpdatedAt = time.Now()
			enc, err := item.encodeEnvironment()
			if err != nil {
				return err
			}
			return b.Put(environmentKey, enc)
		}
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (p *Environment) encodeEnvironment() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decodeEnvironment(data []byte) (*Environment, error) {
	var item *Environment
	err := json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
