package openhab2awsiot

import (
	"fmt"
	"openhab2awsiot/transformer"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestOpenHab2AWSIoT_Transform(t *testing.T) {
	tr := OpenHab2AWSIoT{}
	from := &transformer.Message{
		Topic:   "/openhab/button1/state",
		Payload: []byte("ON"),
	}
	result := &transformer.Message{
		Topic:   "/awsiot/button1",
		Payload: []byte(`{"state":"ON"}`),
	}
	to, err := tr.Transform(from)
	ok(t, err)
	equals(t, result, to)

	from = &transformer.Message{
		Topic:   "/openhab/button1/command",
		Payload: []byte("ON"),
	}
	result = &transformer.Message{
		Topic:   "/awsiot/button1",
		Payload: []byte(`{"command":"ON"}`),
	}
	to, err = tr.Transform(from)
	ok(t, err)
	equals(t, result, to)

	from = &transformer.Message{
		Topic:   "/openhab/level1/state",
		Payload: []byte("42"),
	}
	result = &transformer.Message{
		Topic:   "/awsiot/level1",
		Payload: []byte(`{"state":42}`),
	}
	to, err = tr.Transform(from)
	ok(t, err)
	equals(t, result, to)

	from = &transformer.Message{
		Topic:   "/habopen/level1/state",
		Payload: []byte("42"),
	}
	to, err = tr.Transform(from)
	assert(t, err != nil, "First element error")

	from = &transformer.Message{
		Topic:   "/openhab/level1/unknown",
		Payload: []byte("42"),
	}
	to, err = tr.Transform(from)
	assert(t, err != nil, "Third element error")

	from = &transformer.Message{
		Topic:   "//openhab/level1/state",
		Payload: []byte("42"),
	}
	to, err = tr.Transform(from)
	assert(t, err != nil, "Extra /")
}
