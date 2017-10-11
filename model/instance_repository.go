package model

import (
	"time"

	"encoding/json"

	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
)

var instanceBucket = []byte("instances")

type InstanceRepository interface {
	FindAll() ([]Instance, error)
	FindOne(id string) (*Instance, error)
	FindByIPAddress(IPAddress string) (*Instance, error)
	FindByMACAddress(MACAddress string) (*Instance, error)
	Save(item *Instance) (*Instance, error)
	Delete(id string) error
}

type BoltInstanceRepository struct {
	db *bolt.DB
}

func NewInstanceRepository(db *bolt.DB) *BoltInstanceRepository {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(instanceBucket)
		return err
	})
	return &BoltInstanceRepository{db}
}

func (r *BoltInstanceRepository) FindAll() ([]Instance, error) {
	items := []Instance{}

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(instanceBucket)
		b.ForEach(func(k, v []byte) error {
			item, err := decode(v)
			if err != nil {
				return err
			}
			items = append(items, *item)
			return nil
		})

		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *BoltInstanceRepository) FindOne(id string) (*Instance, error) {
	var item *Instance
	err := r.db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket(instanceBucket)
		k := []byte(id)
		itemData := b.Get(k)
		if len(itemData) == 0 {
			return nil
		}
		item, err = decode(itemData)
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

func (r *BoltInstanceRepository) FindByIPAddress(IPAddress string) (*Instance, error) {
	items, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.IPAddress == IPAddress {
			return &item, nil
		}
	}
	return nil, nil
}

func (r *BoltInstanceRepository) FindByMACAddress(MACAddress string) (*Instance, error) {
	items, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.MACAddress == MACAddress {
			return &item, nil
		}
	}
	return nil, nil
}

func (r *BoltInstanceRepository) Save(item *Instance) (*Instance, error) {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(instanceBucket)
		if item.ID == "" {
			item.ID = bson.NewObjectId()
			item.CreatedAt = time.Now()
			item.UpdatedAt = time.Now()
			enc, err := item.encode()
			if err != nil {
				return err
			}
			return b.Put([]byte(item.ID.Hex()), enc)
		} else {
			item.UpdatedAt = time.Now()
			enc, err := item.encode()
			if err != nil {
				return err
			}
			return b.Put([]byte(item.ID.Hex()), enc)
		}
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *BoltInstanceRepository) Delete(id string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(instanceBucket)
		k := []byte(id)
		return b.Delete(k)
	})
}

func (p *Instance) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Instance, error) {
	var item *Instance
	err := json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
