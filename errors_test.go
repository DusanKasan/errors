package errors_test

import (
	"github.com/DusanKasan/errors"
	"fmt"
	"strings"
	"testing"
	errors2 "errors"
)

func ExampleNew() {
	err := errors.New("code", errors.Data{"key": "value"})
	fmt.Printf("Code: %q\n", err.Code())
	fmt.Printf("Error: %q\n", err.Error())
	fmt.Printf("Data: %v\n", err.Data())
	fmt.Printf("Has frames: %v\n", len(err.Frames()) > 0)
	fmt.Printf("Top frame:\n")
	fmt.Printf("\tLine: %v\n", err.Frames()[0].Line())
	fmt.Printf("\tFunction: %q\n", err.Frames()[0].Function())
	segments := strings.Split(err.Frames()[0].Path(), "/")
	fmt.Printf("\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("Cause: %v\n", err.Cause())

	// Output:
	// Code: "code"
	// Error: "code, data: map[key:value]"
	// Data: map[key:value]
	// Has frames: true
	// Top frame:
	// 	Line: 12
	// 	Function: "github.com/DusanKasan/errors_test.ExampleNew"
	// 	File: "errors_test.go"
	// Cause: <nil>
}

func ExampleWrap() {
	err1 := fmt.Errorf("random error")
	err2 := errors.Wrap(err1, "err2 code", errors.Data{"key": "val"})
	err3 := errors.Wrap(err2, "err3 code", errors.Data{"key2": "val2"})

	fmt.Printf("Code: %q\n", err3.Code())
	fmt.Printf("Error: %q\n", err3.Error())
	fmt.Printf("Data: %v\n", err3.Data())
	fmt.Printf("Frames count: %v\n", len(err3.Frames())) // this error doesn't have frames, only the lowest *errors.E in wrapping chain has frames
	fmt.Printf("Cause:\n")

	err := err3.Cause().(*errors.E)

	fmt.Printf("\tCode: %q\n", errors.Code(err))
	fmt.Printf("\tError: %q\n", err.Error())
	fmt.Printf("\tData: %v\n", err.Data())
	fmt.Printf("\tHas frames: %v\n", len(err.Frames()) > 0)
	fmt.Printf("\tTop frame:\n")
	fmt.Printf("\t\tLine: %v\n", err.Frames()[0].Line())
	fmt.Printf("\t\tFunction: %q\n", err.Frames()[0].Function())
	segments := strings.Split(err.Frames()[0].Path(), "/")
	fmt.Printf("\t\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("\tCause:\n")

	fmt.Printf("\t\tCode: %q\n", errors.Code(err.Cause())) // for non-implementors of `Code() interface{}` this returns empty sentinel value
	fmt.Printf("\t\tError: %q\n", err.Cause().Error())

	// Output:
	// Code: "err3 code"
	// Error: "err3 code, data: map[key2:val2]"
	// Data: map[key2:val2]
	// Frames count: 0
	// Cause:
	// 	Code: "err2 code"
	// 	Error: "err2 code, data: map[key:val]"
	// 	Data: map[key:val]
	// 	Has frames: true
	// 	Top frame:
	// 		Line: 38
	// 		Function: "github.com/DusanKasan/errors_test.ExampleWrap"
	// 		File: "errors_test.go"
	// 	Cause:
	// 		Code: ""
	// 		Error: "random error"


}

func ExampleCode() {
	fmt.Printf("implementor of `Code() interface{}` code: %q\n", errors.Code(errors.New("code"))) // for implementors of `Code() interface{}` this returns err.Code()
	fmt.Printf("non-implementor of `Code() interface{}` code: %q\n", errors.Code(fmt.Errorf("code"))) // for non-implementors of `Code() interface{}` this returns empty sentinel value

	// Output:
	// implementor of `Code() interface{}` code: "code"
	// non-implementor of `Code() interface{}` code: ""
}

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errors.New("code")
	}
}

func BenchmarkNativeNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errors2.New("code")
	}
}

func BenchmarkWrapWithoutStackCapturing(b *testing.B) {
	err := &errors.E{}
	for n := 0; n < b.N; n++ {
		errors.Wrap(err, "code")
	}
}

func BenchmarkWrapWithStackCapturing(b *testing.B) {
	err := errors2.New("err")
	for n := 0; n < b.N; n++ {
		errors.Wrap(err, "code")
	}
}