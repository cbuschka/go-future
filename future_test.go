package future

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPromiseIsPending(t *testing.T) {

	f := NewFuture()

	assert.True(t, f.IsPending())
	assert.False(t, f.IsRejected())
	assert.False(t, f.IsResolved())
}

func TestAfterResolvePromiseIsResolved(t *testing.T) {

	f := NewFuture()

	err := f.Resolve("yup")
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.False(t, f.IsPending())
	assert.False(t, f.IsRejected())
	assert.True(t, f.IsResolved())

	value, err := f.Await()
	assert.Equal(t, value, "yup")
	assert.Nil(t, err)
}

func TestAfterRejectPromiseIsResolved(t *testing.T) {

	f := NewFuture()

	err := f.Reject(fmt.Errorf("rejected"))
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.False(t, f.IsPending())
	assert.True(t, f.IsRejected())
	assert.False(t, f.IsResolved())

	value, err := f.Await()
	assert.Nil(t, value)
	assert.Equal(t, "rejected", err.Error())
}

func TestConcurrentAwaitResolving(t *testing.T) {

	f := NewFuture()

	go func() {
		time.Sleep(1 * time.Millisecond)
		f.Resolve("resolved")
	}()

	value, err := f.Await()

	assert.False(t, f.IsPending())
	assert.False(t, f.IsRejected())
	assert.True(t, f.IsResolved())

	assert.Equal(t, "resolved", value)
	assert.Nil(t, err)
}

func TestConcurrentAwaitRejecting(t *testing.T) {

	f := NewFuture()

	go func() {
		time.Sleep(1 * time.Millisecond)
		f.Reject(fmt.Errorf("rejected"))
	}()

	value, err := f.Await()

	assert.False(t, f.IsPending())
	assert.True(t, f.IsRejected())
	assert.False(t, f.IsResolved())

	assert.Nil(t, value)
	assert.Equal(t, "rejected", err.Error())
}

func TestDoubleAwait(t *testing.T) {

	f := NewFuture()

	go func() {
		time.Sleep(1 * time.Millisecond)
		f.Resolve("resolved")
	}()

	value, err := f.Await()
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "resolved", value)

	assert.False(t, f.IsPending())
	assert.False(t, f.IsRejected())
	assert.True(t, f.IsResolved())

	value, err = f.Await()
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "resolved", value)

	assert.False(t, f.IsPending())
	assert.False(t, f.IsRejected())
	assert.True(t, f.IsResolved())
}
