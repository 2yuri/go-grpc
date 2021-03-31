package service

import (
	"context"
	"errors"
	"fmt"
	"grpc-course/pb"
	"log"
	"sync"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrAlreadyExists = errors.New("record already exists")

//LaptopStore is an interface to store laptop
type LaptopStore interface {
	//Save saves the laptop to the store
	Save(laptop *pb.Laptop) error

	//Find find a laptop by ID
	Find(id string) (*pb.Laptop, error)

	//Search
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

//InMemoryLaptopStore stores laptop in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
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

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]

	if laptop == nil {
		return nil, nil
	}

	other, err := DeepCopy(laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}

func (store *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *pb.Filter,
	found func(laptop *pb.Laptop,
	) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		//time.Sleep(time.Second)
		if ctx.Err() == context.Canceled {
			log.Print("request is canceled")
			return status.Error(codes.Canceled, "request is canceled")
		}

		if ctx.Err() == context.DeadlineExceeded {
			log.Print("deadline is exceeded")
			return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
		}

		if isQualified(filter, laptop) {
			other, err := DeepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() > filter.GetMinCpuCores() {
		return false
	}
	if laptop.GetCpu().GetMinGhz() > filter.GetMinCpuGhz() {
		return false
	}
	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3
	case pb.Memory_KILOBYTE:
		return value << 13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGAYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}
