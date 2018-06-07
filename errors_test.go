package errors_test

import (
	"fmt"
	"testing"

	"github.com/DusanKasan/errors"
	"strings"
)

func TestNew(t *testing.T) {
	err := errors.New("msg/code", errors.Data{"key": "val"})

	expect(t, err,
		code("msg/code"),
		data(errors.Data{"key": "val"}),
		origin("errors_test.go", 12, "github.com/DusanKasan/errors_test.TestNew"),
		noCause(),
	)
}

func TestWrap(t *testing.T) {
	err1 := errors.New("msg/code", errors.Data{"key": "val"})
	err2 := errors.Wrap(err1,  "msg/code2", errors.Data{"key2": "val2"})

	expect(t, err2,
		code("msg/code2"),
		data(errors.Data{"key2": "val2"}),
		framesCount(0),
		cause(
			code("msg/code"),
			data(errors.Data{"key": "val"}),
			origin("errors_test.go", 23, "github.com/DusanKasan/errors_test.TestWrap"),
			noCause(),
		),
	)
}

func TestCode(t *testing.T) {
	var err error

	err = errors.New("msg/code", errors.Data{"key": "val"})
	if errors.Code(err) != "msg/code" {
		t.Errorf("bad error code. expected: %v, got: %v", 1, errors.Code(err))
	}

	err = fmt.Errorf("err")
	if errors.Code(err) == nil {
		t.Errorf("no error code")
	}
}

type expectation func(*errors.E) []error

func expect(t *testing.T, e *errors.E, expectations ...expectation) {
	for _, exp := range expectations {
		if errs := exp(e); errs != nil {
			for _, err := range errs {
				t.Error(err)
			}
		}
	}
}

func data(fs errors.Data) expectation {
	return func(e *errors.E) []error {
		if len(fs) != len(e.Data()) {
			return []error{fmt.Errorf("bad data item count. expected: %v, got: %v", len(fs), len(e.Data()))}
		}

		var errs []error
		for k, v := range fs {
			if _, ok := e.Data()[k]; !ok {
				errs = append(errs, fmt.Errorf("data item with key %q not found", k))
				continue
			}

			if e.Data()[k] != v {
				errs = append(errs, fmt.Errorf("data item with key %q has bad value. expected: %v, got: %v", k, v, e.Data()[k]))
			}
		}

		return errs
	}
}

func code(c interface{}) expectation {
	return func(e *errors.E) []error {
		if e.Code() != c {
			return []error{fmt.Errorf("bad error code. expected: %v, got: %v", c, e.Code())}
		}

		return nil
	}
}

func framesCount(c int) expectation {
	return func(e *errors.E) []error {
		if len(e.Frames()) != c {
			return []error{fmt.Errorf("bad frame count. expected: %v, got: %v", c, len(e.Frames()))}
		}

		return nil
	}
}

func origin(file string, line int, function string) expectation {
	return func(e *errors.E) []error {
		if len(e.Frames()) == 0 {
			return []error{fmt.Errorf("unable to determine origin, no frames in error")}
		}

		frame := e.Frames()[0]
		var errs []error

		if !strings.HasSuffix(frame.Path(), file) {
			errs = append(errs, fmt.Errorf("bad origin file name. expected: %v, got: %v", file, frame.Path()))
		}

		if frame.Line() != line {
			errs = append(errs, fmt.Errorf("bad origin line number. expected: %v, got: %v", line, frame.Line()))
		}

		if frame.Function() != function {
			errs = append(errs, fmt.Errorf("bad origin function. expected: %v, got: %v", function, frame.Function()))
		}

		return errs
	}
}

func cause(expectations ...expectation) expectation {
	return func(e *errors.E) []error {
		err, ok := e.Cause().(*errors.E)
		if !ok {
			return []error{fmt.Errorf("cause is not errors.E, got: %#v", e)}
		}

		var errs []error

		for _, exp := range expectations {
			errs = append(errs, exp(err)...)
		}

		return errs
	}
}

func noCause() expectation {
	return func(e *errors.E) []error {
		if e.Cause() != nil {
			return []error{fmt.Errorf("cause is not empty, got: %#v", e.Cause())}
		}

		return nil
	}
}

