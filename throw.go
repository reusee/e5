package e5

type throw struct {
	err error
}

func (t *throw) String() string {
	return t.err.Error()
}

func (t *throw) Error() string {
	return t.err.Error()
}

func (t *throw) Unwrap() error {
	return t.err
}

// Throw checks the error and if not nil, raise a panic which will be recovered by Handle
func Throw(err error, args ...error) error {
	if err == nil {
		return nil
	}
	if len(args) > 0 {
		err = Wrap.With(args...)(err)
	}
	if err == nil {
		return nil
	}
	panic(&throw{
		err: err,
	})
}
