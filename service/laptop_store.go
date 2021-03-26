package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"grpc-course/pb"
	"sync"
)

var ErrAlreadyExists = errors.New("record already exists")

//LaptopStore is an interface to store laptop
type LaptopStore interface {
	//Save saves the laptop to the store
	Save (laptop *pb.Laptop) error
}

//InMemoryLaptopStore stores laptop in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data map[string]*pb.Laptop
}

//NewInMemoryLaptopStore returns a new InMemoryLaptopStore

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other, err := DeepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	return nil
}

//DeepCopy make a deep copy
func DeepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}