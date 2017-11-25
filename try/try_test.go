package try_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bruston/try"
)

var someErr = errors.New("some error")

func TestDo(t *testing.T) {
	buf := &bytes.Buffer{}
	err := try.Do(5, func() error {
		buf.WriteByte('.')
		return someErr
	}, nil)
	if err != someErr {
		t.Errorf("expecting error :%s, got %s", someErr, err)
	}
	const expected = "....."
	if buf.String() != expected {
		t.Errorf("expecting %x, got %x", expected, buf.String())
	}
}

func TestDoSucceeds(t *testing.T) {
	c := 0
	err := try.Do(5, func() error {
		if c == 3 {
			return nil
		}
		c++
		return someErr
	}, nil)
	if err != nil {
		t.Errorf("errs should be nil")
	}
}

func TestDelay(t *testing.T) {
	duration := time.Second
	start := time.Now()
	delay := try.Delay(duration)
	delay.Throttle()
	elapsed := time.Since(start)
	if elapsed < duration {
		t.Errorf("delay was: %s, expecting no less than 1 second", elapsed)
	}
}

func TestBackoff(t *testing.T) {
	const (
		multiplier = 1
		duration   = time.Second
	)
	backoff := &try.Backoff{
		Duration:   duration,
		Multiplier: multiplier,
	}
	start := time.Now()
	backoff.Throttle()
	backoff.Throttle()
	elapsed := time.Since(start)
	if elapsed < duration*multiplier {
		t.Errorf("delay was %s, expecting no less than 3 seconds", elapsed)
	}
}

type testThrottler struct {
	throttled bool
}

func (th *testThrottler) Throttle() {
	th.throttled = true
}

func TestDoThrottled(t *testing.T) {
	throttler := &testThrottler{}
	try.Do(2, func() error {
		return someErr
	}, throttler)
	if !throttler.throttled {
		t.Errorf("throttler was not called")
	}
}

func Example() {
	buf := &bytes.Buffer{}
	try.Do(5, func() error {
		buf.WriteString(".")
		return errors.New("this function can fail")
	}, nil)
	fmt.Println(buf.String())
	// Output:
	// .....
}

func ExampleDelay() {
	try.Do(5, func() error {
		return errors.New("this function can fail")
	}, try.Delay(time.Second))
}

func ExampleBackoff() {
	backoff := &try.Backoff{
		Duration:   time.Microsecond * 100,
		Multiplier: 2,
	}
	try.Do(
		5, func() error {
			return errors.New("this function can fail")
		},
		backoff,
	)
}
