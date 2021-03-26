package workflow

import (
	"context"
	"math"
	"sync"
)

const abortIndex int8 = math.MaxInt8 / 2

type Event = func(ctx context.Context, data *WorkData) error

type funcs struct {
	mux   sync.Mutex
	fs    []Event
	index int8
}

func NewFuncs() *funcs {
	return &funcs{
		fs:    []Event{},
		index: -1,
	}
}

func (fs *funcs) Add(f Event) {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	fs.fs = append(fs.fs, f)
}

func (fs *funcs) next(ctx context.Context, data *WorkData) error {
	fs.index++
	for fs.index < int8(len(fs.fs)) {
		err := fs.fs[fs.index](ctx, data)
		if err != nil {
			fs.abort()
			return err
		}
		fs.index++
	}
	return nil
}

func (fs *funcs) isAborted() bool {
	return fs.index >= abortIndex
}

func (fs *funcs) abort() {
	fs.index = abortIndex
}

func (fs *funcs) reset() {
	fs.index = -1
}
