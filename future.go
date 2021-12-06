package future

import (
	"fmt"
	"sync"
)

const (
	pending  = 0
	resolved = 1
	rejected = 2
)

type Future interface {
	IsRejected() bool
	IsResolved() bool
	IsPending() bool
	Resolve(value interface{}) error
	Reject(err error) error
	MustResolve(value interface{})
	MustReject(err error)
	Await() (interface{}, error)
}

type defaultFuture struct {
	cond   *sync.Cond
	status int8
	err    error
	value  interface{}
}

func NewFuture() Future {
	lock := &sync.Mutex{}
	cond := sync.NewCond(lock)

	return Future(&defaultFuture{cond: cond, status: pending, err: nil, value: nil})
}

func ResolvedFuture(value interface{}) Future {
	future := NewFuture()
	future.MustResolve(value)
	return future
}

func RejectedFuture(err error) Future {
	future := NewFuture()
	future.MustReject(err)
	return future
}

func (future *defaultFuture) MustResolve(value interface{}) {
	err := future.Resolve(value)
	if err != nil {
		panic(err)
	}
}

func (future *defaultFuture) Resolve(value interface{}) error {
	future.cond.L.Lock()
	defer future.cond.L.Unlock()
	switch future.status {
	case resolved:
		return fmt.Errorf("already resolved")
	case rejected:
		return fmt.Errorf("already rejected")
	case pending:
		future.status = resolved
		future.value = value
		future.err = nil
		future.cond.Broadcast()
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}

func (future *defaultFuture) MustReject(err error) {
	rejectionRerr := future.Reject(err)
	if rejectionRerr != nil {
		panic(rejectionRerr)
	}
}

func (future *defaultFuture) Reject(err error) error {
	future.cond.L.Lock()
	defer future.cond.L.Unlock()
	switch future.status {
	case resolved:
		return fmt.Errorf("already resolved")
	case rejected:
		return fmt.Errorf("already rejected")
	case pending:
		future.status = rejected
		future.value = nil
		future.err = err
		future.cond.Broadcast()
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}

func (future *defaultFuture) IsResolved() bool {
	future.cond.L.Lock()
	defer future.cond.L.Unlock()
	switch future.status {
	case resolved:
		return true
	case rejected:
		return false
	case pending:
		return false
	default:
		panic(fmt.Errorf("invalid status"))
	}
}

func (future *defaultFuture) IsPending() bool {
	future.cond.L.Lock()
	defer future.cond.L.Unlock()
	switch future.status {
	case resolved:
		return false
	case rejected:
		return false
	case pending:
		return true
	default:
		panic(fmt.Errorf("invalid status"))
	}
}

func (future *defaultFuture) IsRejected() bool {
	future.cond.L.Lock()
	defer future.cond.L.Unlock()
	switch future.status {
	case resolved:
		return false
	case rejected:
		return true
	case pending:
		return false
	default:
		panic(fmt.Errorf("invalid status"))
	}
}

func (future *defaultFuture) Await() (interface{}, error) {

	future.cond.L.Lock()
	defer future.cond.L.Unlock()
	if future.status == pending {
		future.cond.Wait()
	}
	switch future.status {
	case resolved:
		return future.value, nil
	case rejected:
		return nil, future.err
	case pending:
		return nil, fmt.Errorf("still pending")
	default:
		return nil, fmt.Errorf("invalid status")
	}
}
