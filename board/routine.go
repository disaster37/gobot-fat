package board

import "context"

type Routine struct {
	chErr    chan error
	chRes    chan interface{}
	chCancel chan bool
}

// Permit to cancel current routine
func (r *Routine) Cancel() {
	r.chCancel <- true
}

func (r *Routine) Error() chan error {
	return r.chErr
}

func (r *Routine) Result() chan interface{} {
	return r.chRes
}

func NewRoutine(ctx context.Context, process func(ctx context.Context, chCancel chan bool) (res interface{}, err error)) *Routine {
	routine := &Routine{
		chErr:    make(chan error),
		chRes:    make(chan interface{}),
		chCancel: make(chan bool),
	}

	go func() {
		res, err := process(ctx, routine.chCancel)
		if err != nil {
			routine.chErr <- err
			return
		}

		routine.chRes <- res
	}()

	return routine
}
