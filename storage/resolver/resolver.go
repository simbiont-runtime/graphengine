// ---

package resolver

import (
	"context"
	"sync/atomic"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/internal/logutil"
)

type resolver struct {
	db     *pebble.DB
	ch     chan Task
	closed atomic.Bool
}

func newResolver(db *pebble.DB) *resolver {
	return &resolver{
		db: db,
		ch: make(chan Task, 512),
	}
}

func (r *resolver) run(ctx context.Context) {
	for {
		select {
		case task, ok := <-r.ch:
			if !ok {
				return
			}
			c := len(r.ch)
			tasks := make([]Task, 0, c+1)
			tasks = append(tasks, task)
			if c > 0 {
				for i := 0; i < c; i++ {
					tasks = append(tasks, <-r.ch)
				}
			}
			r.resolve(tasks)

		case <-ctx.Done():
			return
		}
	}
}

func (r *resolver) resolve(tasks []Task) {
	batch := r.db.NewBatch()
	for _, task := range tasks {
		var err error
		if task.CommitVer > 0 {
			err = Resolve(r.db, batch, task.Key, task.StartVer, task.CommitVer)
		} else {
			err = Rollback(r.db, batch, task.Key, task.StartVer)
		}
		if err != nil {
			logutil.Errorf("Resolve key failed, key:%v, startVer:%d, commitVer:%d, caused by:%+v",
				task.Key, task.StartVer, task.CommitVer, err)
		}
		if task.Notifier != nil {
			task.Notifier.Notify(err)
		}
	}
	err := batch.Commit(nil)
	if err != nil {
		logutil.Errorf("Commit batch failed: %+v", err)
	}
}

func (r *resolver) push(tasks ...Task) {
	for _, t := range tasks {
		r.ch <- t
	}
}

func (r *resolver) close() {
	if r.closed.Swap(true) {
		return
	}
	close(r.ch)
}
