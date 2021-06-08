package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Row struct {
	Key       string
	Value     interface{}
	CreatedAt int64
}

type Store struct {
	data map[string]Row
	lock sync.Mutex
	file string
}

func New(file string) (*Store, error) {
	store := Store{
		data: make(map[string]Row),
		lock: sync.Mutex{},
		file: file,
	}
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			os.Create(file)
		}
	} else {
		return nil, err
	}
	return &store, nil
}
func (s *Store) Set(key string, value interface{}) (Row, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	row := Row{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().Unix(),
	}
	s.data[key] = row
	return row, nil
}

func (s *Store) Get(key string) (Row, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.data[key] == (Row{}) {
		return Row{}, fmt.Errorf("key %s not found", key)
	}
	return s.data[key], nil
}

func (s *Store) Delete(key string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.data[key] == (Row{}) {
		return fmt.Errorf("key %s not found", key)
	}
	delete(s.data, key)
	return nil
}

func (s *Store) Keys() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	keys := []string{}
	for key, _ := range s.data {
		keys = append(keys, key)
	}
	return keys
}

func (s *Store) Save() error {
	f, err := os.OpenFile(s.file, os.O_RDWR, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	s.lock.Lock()
	defer s.lock.Unlock()

	byteValue, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	_, err = f.Write(byteValue)
	return err
}

func (s *Store) load() error {
	f, err := os.OpenFile(s.file, os.O_RDONLY, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	s.lock.Lock()
	defer s.lock.Unlock()

	var dataBytes []byte
	byteRead := 100
	for byteRead > 0 {
		byteRead, _ = f.Read(dataBytes)
	}
	if err := json.Unmarshal(dataBytes, &s.data); err != nil {
		return err
	}
	return nil
}
