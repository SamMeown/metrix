package backoff

import (
	"errors"
	"fmt"
	"time"
)

type RetryableError struct {
	Err error
}

func NewRetryableError(err error) error {
	return &RetryableError{
		Err: err,
	}
}

func (re *RetryableError) Error() string {
	return fmt.Sprintf("retryable error: %s", re.Err.Error())
}

func (re *RetryableError) Unwrap() error {
	return re.Err
}

type Backoff struct {
	retries        []int
	isRetryableErr func(error) bool
}

func NewBackoff(retries []int, retryableErrFunc func(error) bool) Backoff {
	return Backoff{
		retries:        retries,
		isRetryableErr: retryableErrFunc,
	}
}

type Retryable func() error

func (b Backoff) Retry(f Retryable) error {
	var err error
	for attempt := 0; ; attempt++ {
		err = f()
		if err == nil || attempt > len(b.retries)-1 {
			break
		}

		var retErr *RetryableError
		if !errors.As(err, &retErr) &&
			(b.isRetryableErr == nil || !b.isRetryableErr(err)) {
			break
		}

		time.Sleep(time.Duration(b.retries[attempt]) * time.Second)
	}

	return err
}
