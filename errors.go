package errors

import (
	"fmt"
	"runtime"
)

type code string
const codeNotProvided = code("")

// New returns an error with the specified code, message and data.
func New(code interface{}, data ...Data) *E {
	if code == nil {
		code = codeNotProvided
	}

	err := &E{
		code:    code,
		data:    Data{},
		callers: callers(),
	}

	for _, f := range data {
		for k, v := range f {
			err.data[k] = v
		}
	}

	return err
}

// Wrap returns an error with the specified cause, code message and data.
func Wrap(cause error, code interface{}, data ...Data) *E {
	if code == nil {
		code = codeNotProvided
	}

	err := &E{
		code:  code,
		data:  Data{},
		cause: cause,
	}

	if _, ok := cause.(*E); !ok {
		err.callers = callers()
	}

	for _, f := range data {
		for k, v := range f {
			err.data[k] = v
		}
	}

	return err
}

// Code returns the error code from the passed error. If the passed error doesn't
// implement the `Code() interface{}` method, it will return nil.
func Code(err error) interface{} {
	if err == nil {
		return nil
	}

	if e, ok := err.(*E); ok {
		return e.code
	}

	return codeNotProvided
}

type Data map[string]interface{}

// E represents an error with code, message, data, stack trace and cause.
type E struct {
	code    interface{}
	data    Data
	callers []uintptr
	cause   error
}

// Error returns the message with data in it.
func (e *E) Error() string {
	msg := fmt.Sprintf("%v", e.code)
	if len(e.data) > 0 {
		msg = msg + fmt.Sprintf(", data: %v", e.data)
	}

	return msg
}

// Code returns the error code.
func (e *E) Code() interface{} {
	return e.code
}

// Cause returns the error that cause this one if this error was created using
// Wrap. Otherwise returns nil.
func (e *E) Cause() error {
	return e.cause
}

// Data returns the data.
func (e *E) Data() Data {
	return e.data
}

// Frames returns the stack frames of the stored error.
func (e *E) Frames() []Frame {
	var frames []Frame
	for _, pc := range e.callers {
		line, file, function := 0, "UNKNOWN", "UNKNOWN"
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			file, line = fn.FileLine(pc)
			function = fn.Name()
		}

		frames = append(frames, &frame{line, file, function})
	}

	return frames
}

type Frame interface {
	// Line returns the line number for this frame
	Line() int
	// Path returns the file path for this frame
	Path() string
	// Function returns the function name for this frame
	Function() string
}

type frame struct {
	line int
	path string
	function string
}

func (f *frame) Line() int {
	return f.line
}

func (f *frame) Path() string {
	return f.path
}

func (f *frame) Function() string {
	return f.function
}

// callers returns the slice of program counter without the top 4 (internal) frames
// which are always golang's internal calls so we don't need them.
func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

