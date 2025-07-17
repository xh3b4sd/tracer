package tracer

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Tracer_Json_Interface ensures that this error implementation complies
// with the native JSON Marschaler interface. This test already fails at compile
// time if *Error does not implement json.Marshaler.
func Test_Tracer_Json_Interface(t *testing.T) {
	var _ json.Marshaler = &Error{}
}

// Test_Tracer_Json_String ensures that the error wrapping, error matching and
// respective JSON encoding works properly for the *Error type, including its
// context and tracing annotations.
//
//	go test ./... -run Test_Tracer_Json_String -update
func Test_Tracer_Json_String(t *testing.T) {
	var (
		testErrorOne       = fmt.Errorf("test error one message")
		testErrorTwo       = &Error{Description: "test error two description"}
		testErrorThree     = fmt.Errorf("executing \".github/dependabot.yaml\"")
		testErrorFour      = &Error{}
		alreadyExistsError = &Error{Description: "alreadyExistsError", Context: []Context{{Key: "code", Val: "invalidArgument"}}}
	)

	testCases := []struct {
		err error
		cau error
		neg error
	}{
		// Case 000 ensures that without error there is no cause.
		{
			err: nil,
			cau: nil,
			neg: testErrorTwo,
		},
		// Case 001 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			err: func() error {
				return testErrorOne
			}(),
			cau: testErrorOne,
			neg: testErrorTwo,
		},
		// Case 002 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			err: func() error {
				return testErrorTwo
			}(),
			cau: testErrorTwo,
			neg: testErrorThree,
		},
		// Case 003 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			err: func() error {
				return Mask(testErrorOne)
			}(),
			cau: testErrorOne,
			neg: nil,
		},
		// Case 004 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			err: func() error {
				return Mask(testErrorTwo)
			}(),
			cau: testErrorTwo,
			neg: testErrorThree,
		},
		// Case 005 does error wrapping two times. The error is simply made up by
		// fmt.Errorf.
		{
			err: func() error {
				var err error

				err = Mask(testErrorOne)
				err = Mask(err)

				return err
			}(),
			cau: testErrorOne,
			neg: testErrorThree,
		},
		// Case 006 does error wrapping two times. The error is simply made up by
		// using the tracer error type.
		{
			err: func() error {
				var err error

				err = Mask(testErrorTwo)
				err = Mask(err)

				return err
			}(),
			cau: testErrorTwo,
			neg: testErrorOne,
		},
		// Case 007 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			err: func() error {
				var err error // nolint:gosimple

				err = Mask(testErrorTwo, Context{Key: "annotation", Val: "foo bar"})

				return err
			}(),
			cau: testErrorTwo,
			neg: testErrorOne,
		},
		// Case 008 does error wrapping two times, while the first wrapping is
		// annotated. The annotations use the same key twice.
		{
			err: func() error {
				var err error

				err = Mask(testErrorTwo,
					Context{Key: "annotation", Val: "foo bar"},
					Context{Key: "something", Val: "one two"},
					Context{Key: "annotation", Val: "567 xyz"},
				)
				err = Mask(err)

				return err
			}(),
			cau: testErrorTwo,
			neg: testErrorOne,
		},
		// Case 009 does not use error wrapping. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			err: func() error {
				var err error // nolint:gosimple

				err = testErrorThree

				return err
			}(),
			cau: testErrorThree,
			neg: testErrorOne,
		},
		// Case 010 does error wrapping one time. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			err: func() error {
				var err error // nolint:gosimple

				err = Mask(testErrorThree)

				return err
			}(),
			cau: testErrorThree,
			neg: nil,
		},
		// Case 011 does error wrapping two times. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			err: func() error {
				var err error

				err = Mask(testErrorThree)
				err = Mask(err)

				return err
			}(),
			cau: testErrorThree,
		},
		// Case 012 does error wrapping twice. The error type is *Error. The
		// original instance contains context annotations already.
		{
			err: func() error {
				var err error

				err = Mask(alreadyExistsError)
				// keep this line between the two Mask calls
				err = Mask(err)

				return err
			}(),
			cau: alreadyExistsError,
			neg: testErrorOne,
		},
		// Case 013 wraps nil errors at first. The first wrapped error wraps with
		// context annotations.
		{
			err: func() error {
				var err error

				err = Mask(err)
				err = Mask(err) // nolint:staticcheck,ineffassign
				err = Mask(testErrorTwo, Context{Key: "re-source", Val: "id"})
				err = Mask(err)
				err = Mask(err, Context{Key: "more", Val: "info"})
				err = Mask(err)

				return err
			}(),
			cau: testErrorTwo,
			neg: testErrorOne,
		},
		// Case 014 does error wrapping one time. The error is of type *Error
		// without any description.
		{
			err: func() error {
				var err error // nolint:gosimple

				err = Mask(testErrorFour)

				return err
			}(),
			cau: testErrorFour,
			neg: alreadyExistsError,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			{
				if !errors.Is(tc.err, tc.cau) {
					t.Fatalf("expected %#v got %#v", tc.cau, tc.err)
				}
				if errors.Is(tc.err, tc.neg) {
					t.Fatalf("expected %#v got %#v", false, true)
				}
			}

			{
				act, exp := actExp(filepath.Join("testdata", fmt.Sprintf("case.%03d.golden", i)), Json(tc.err))
				if dif := cmp.Diff(exp, act); dif != "" {
					t.Fatalf("-expected +actual:\n%s", dif)
				}
			}
		})
	}
}

func actExp(pat string, inp string) (string, string) {
	var err error

	// Change the given file paths in our Error type's JSON output to
	// avoid prefixes like "/Users/username/go/src/" so this test can be
	// executed on different machines. See the golden files in the
	// testdata folder for specific examples.
	var act []byte
	{
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		{
			inp = strings.ReplaceAll(inp, cwd, "--REPLACED--")
		}

		b := &bytes.Buffer{}
		err = json.Indent(b, []byte(inp), "", "\t")
		if err != nil {
			panic(err)
		}

		act = []byte(b.String() + "\n")
	}

	if *update {
		err := os.WriteFile(pat, act, 0644) // nolint:gosec
		if err != nil {
			panic(err)
		}
	}

	var exp []byte
	{
		exp, err = os.ReadFile(pat)
		if err != nil {
			panic(err)
		}
	}

	return string(act), string(exp)
}
