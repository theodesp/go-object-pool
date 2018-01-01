package go_object_pool

import (
	"errors"
	"sync"
)

type PooledObject interface {
	Reset()
}

// All pooled object factories must satisfy this interface
type PooledObjectFactory interface {
	Create() (PooledObject, error)
}

// All pool implementations must satisfy this interface
type Pool interface {
	Get() (PooledObject, error)
	Return(obj PooledObject) error
}

/**
 * Fixed size Object pool
 */
type FixedPool struct {
	// List of available objects to share
	available []PooledObject
	// List of in-use objects that are currently reserved
	inUse []PooledObject
	// Maximum size permitted
	capacity int
	// For protecting updates
	mu *sync.Mutex
	// For creating or closing the Objects
	factory PooledObjectFactory
}

// Creates a new fixed pool with capacity
func NewFixedPool(capacity int, factory PooledObjectFactory) *FixedPool {
	return &FixedPool{
		available: make([]PooledObject, 0),
		inUse:     make([]PooledObject, 0),
		capacity:  capacity,
		mu:        new(sync.Mutex),
		factory:   factory,
	}
}

func (p *FixedPool) Get() (PooledObject, error) {
	p.mu.Lock()

	var obj PooledObject
	var err error

	if len(p.available) == 0 {
		// Make sure we don't exceed capacity
		if len(p.inUse) == p.capacity {
			err = errors.New("fixed Pool reached maximum capacity")
		} else {
			obj, err = p.factory.Create()
			p.inUse = append(p.inUse, obj)
		}
	} else {
		// pop
		obj, p.available = p.available[0], p.available[1:]
		err = nil
		p.inUse = append(p.inUse, obj)
	}

	p.mu.Unlock()

	return obj, err
}

func (p *FixedPool) Return(obj PooledObject) error {
	obj.Reset()

	var err error

	p.mu.Lock()
	if idx := findIndex(obj, p.inUse); idx != -1 {
		// Delete at index
		p.inUse = append(p.inUse[:idx], p.inUse[idx+1:]...)
		p.available = append(p.available, obj)
		err = nil
	} else {
		err = errors.New("unrecognized pooled object returned")
	}
	p.mu.Unlock()

	return err
}

func findIndex(target PooledObject, slice []PooledObject) int {
	for idx, obj := range slice {
		if target == obj {
			return idx
		}
	}

	return -1
}
