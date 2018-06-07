# errors

This is an opinionated error package for Go. It provides errors with stack traces, error codes, structured data and error re-wrapping. It does these things with minimal API and has extendability in mind.

For more information see godoc at http://godoc.org/github.com/DusanKasan/errors.

## Example usage

Here's a simple example of creating an error, re-wrapping it and then making a decision based on the error type.

```go
type code string

const CodeUserNotFound := code("user not found")

// return the error with specified code
func GetUserEmail(username string) (string, error) {
    return errors.New(CodeUserNotFound, errors.Data{"username": username})
}

func SendEmailToUser(email Email, username string) error {
    email, err := GetUserEmail(username)
    switch errors.Code(err) {
        case nil: // no error, we can continue
        case CodeUserNotFound:
            // process this specific error
        default:
            // wrap the error to add context and stop any
            return errors.Wrap(error, "unable to send email to user")
    }
    ...
}

```

## Error code instead of a message

The errors produced by this package no longer use the term "message". Instead they use "code" in its place. Code can be anything but most commonly it will be a string type variable. The reason for this change is that while it provides the same functionality as a message, it also allows you to create a exported constant of an unexported type and use it as a sentinel value. This means you can make decisions based on the error code.

The key part here is that the `errors.Code(error) interface{}` function returns the code of the error you pass to it. There are three distinct paths it can take:

 - if the error implements a `Code() interface{}` method, and its return value is non-nil, return that value
 - if the error doesn't implement `Code() interface{}` or its return value is nil, return a backup unexported sentinel value (empty string)
 - if the error is nil, return nil

## Structured data

Both `errors.New` and `errors.Wrap` accept `errors.Data` as their last parameter (which is also optional). `errors.Data` is simply and alias for `map[string]interface{}` and it allows you to add any contextual data to the error without printing them as is standard with the built-in error package. This is useful for example when outputting the error to a structured log.

## Error wrapping

A lot of times when working with error you just want to add more context to the error you already received from downstream. In the built-in errors package you do this by appending the error into a new error string. So you end up with stuff like `fmt.Errorf("New error description: %q", err)`.

This package provides the `errors.Wrap(cause error, code interface{}, data ...errors.Data) *errors.E` function that will create a new error from the specified inputs and sets the `cause` as a cause for the newly created error.

## Stack traces

You can get the stack trace for the current error by calling the `Frames() []errors.Frame` on `*errors.E`. It will return a slice of `Frame`s that you can read the file, line number and function name from. Note that only the first `errors.E` in the wrapping chain will have a stack trace (because capturing stack is costly).

## TODO

- better example