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
	retries []int
}

func NewBackoff(retries []int) Backoff {
	return Backoff{
		retries: retries,
	}
}

type Retryable func() error

func (b Backoff) RetryableFunc(f Retryable) Retryable {
	return func() error {
		var err error
		for attempt := 0; ; attempt++ {
			err = f()
			if err == nil || attempt > len(b.retries)-1 {
				break
			}

			var retErr *RetryableError
			if !errors.As(err, &retErr) {
				break
			}

			time.Sleep(time.Duration(b.retries[attempt]) * time.Second)
		}

		return err
	}
}

func (b Backoff) Retry(f Retryable) error {
	return b.RetryableFunc(f)()
}
