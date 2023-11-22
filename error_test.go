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

func Test_Tracer_Error_Error(t *testing.T) {
	var (
		testErrorOne       = fmt.Errorf("test error one")
		testErrorTwo       = &Error{Kind: "testErrorTwo"}
		testErrorThree     = fmt.Errorf("executing \".github/dependabot.yaml\"")
		alreadyExistsError = &Error{Kind: "alreadyExistsError"}
	)

	testCases := []struct {
		errFunc func() error
		message string
	}{
		// Case 000 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
			message: "test error one",
		},
		// Case 001 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
			message: "test error two",
		},
		// Case 002 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
			message: "test error one",
		},
		// Case 003 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
			message: "test error two",
		},
		// Case 004 does error wrapping two times. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorOne)
				err = Mask(err)

				return err
			},
			message: "test error one",
		},
		// Case 005 does error wrapping two times. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorTwo)
				err = Mask(err)

				return err
			},
			message: "test error two",
		},
		// Case 006 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
			message: "test error two: annotation",
		},
		// Case 007 does error wrapping two times, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				var err error

				err = Maskf(testErrorTwo, "annotation")
				err = Mask(err)

				return err
			},
			message: "test error two: annotation",
		},
		// Case 008 does not use error wrapping. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return testErrorThree
			},
			message: "executing \".github/dependabot.yaml\"",
		},
		// Case 009 does error wrapping one time. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return Mask(testErrorThree)
			},
			message: "executing \".github/dependabot.yaml\"",
		},
		// Case 010 does error wrapping two times. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorThree)
				err = Mask(err)

				return err
			},
			message: "executing \".github/dependabot.yaml\"",
		},
		// Case 011 ensures that masking errors without annotation generates the
		// error message representation based on its error kind.
		{
			errFunc: func() error {
				return Mask(alreadyExistsError)
			},
			message: "already exists error",
		},
		// Case 012 is similar to the above, but additionally combines the error kind
		// and error annotation in a way that the " error" suffix of the error kind
		// based error message is removed in order to produce a more natural human
		// readable error message.
		{
			errFunc: func() error {
				return Maskf(alreadyExistsError, "some thing")
			},
			message: "already exists: some thing",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			message := tc.errFunc().Error()

			if message != tc.message {
				t.Fatalf("expected %#v got %#v", tc.message, message)
			}
		})
	}
}

func Test_Tracer_Error_Copy(t *testing.T) {
	a := &Error{Kind: "testError"}
	b := a.Copy()

	if !errors.Is(a, a) {
		t.Fatalf("%#v and %#v must be equal", a, a)
	}
	if errors.Is(a, b) {
		t.Fatalf("%#v and %#v must not be equal", a, b)
	}
}

func Test_Tracer_Error_Is(t *testing.T) {
	var (
		testErrorOne         = fmt.Errorf("testErrorOne")
		testErrorOneSameKind = fmt.Errorf("testErrorOne")
		testErrorTwo         = &Error{Kind: "testErrorTwo"}
		testErrorTwoSameKind = &Error{Kind: "testErrorTwo"}
		testErrorThree       = &Error{Kind: "testErrorThree"}
	)

	testCases := []struct {
		one   error
		two   error
		equal bool
	}{
		// Case 000 ensures that two equal errors are detected to be equal. Both
		// errors are of the same arbitrary error type.
		{
			one:   testErrorOne,
			two:   testErrorOne,
			equal: true,
		},
		// Case 001 ensures that two equal errors are detected to be equal.
		{
			one:   testErrorTwo,
			two:   testErrorTwo,
			equal: true,
		},
		// Case 002 ensures that two different errors are not detected to be
		// equal. One error is our Error type. The other other error is some
		// arbitrary error type.
		{
			one:   testErrorTwo,
			two:   testErrorOne,
			equal: false,
		},
		// Case 003 ensures that two different errors are not detected to be
		// equal. Both errors are our Error types, but each is a different
		// instance.
		{
			one:   testErrorTwo,
			two:   testErrorThree,
			equal: false,
		},
		// Case 004 ensures that two different errors are not detected to be
		// equal. Both errors are our Error types. Both errors have the same
		// kind but each is a different instance.
		{
			one:   testErrorTwo,
			two:   testErrorTwoSameKind,
			equal: false,
		},
		// Case 005 ensures that two different errors are not detected to be
		// equal. Both errors are of the same arbitrary error type. Both errors
		// have the same kind but each is a different instance.
		{
			one:   testErrorOne,
			two:   testErrorOneSameKind,
			equal: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			equal := errors.Is(tc.one, tc.two)

			if equal != tc.equal {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.equal, equal))
			}
		})
	}
}

func Test_Tracer_Error_JSON(t *testing.T) {
	var (
		testErrorOne       = fmt.Errorf("test error one")
		testErrorTwo       = &Error{Kind: "testErrorTwo"}
		testErrorThree     = fmt.Errorf("executing \".github/dependabot.yaml\"")
		alreadyExistsError = &Error{Kind: "alreadyExistsError", Code: "invalidArgument"}
	)

	testCases := []struct {
		errFunc func() error
	}{
		// Case 000 ensures that without error there is no cause.
		{
			errFunc: func() error {
				return nil
			},
		},
		// Case 001 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
		},
		// Case 002 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
		},
		// Case 003 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
		},
		// Case 004 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
		},
		// Case 005 does error wrapping two times. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorOne)
				err = Mask(err)

				return err
			},
		},
		// Case 006 does error wrapping two times. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorTwo)
				err = Mask(err)

				return err
			},
		},
		// Case 007 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
		},
		// Case 008 does error wrapping two times, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				var err error

				err = Maskf(testErrorTwo, "annotation")
				err = Mask(err)

				return err
			},
		},
		// Case 009 does not use error wrapping. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return testErrorThree
			},
		},
		// Case 010 does error wrapping one time. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return Mask(testErrorThree)
			},
		},
		// Case 011 does error wrapping two times. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorThree)
				err = Mask(err)

				return err
			},
		},
		// Case 012
		{
			errFunc: func() error {
				var err error

				err = Mask(alreadyExistsError)
				err = Mask(err)

				return err
			},
		},
		// Case 013
		{
			errFunc: func() error {
				var err error

				err = Maskf(alreadyExistsError, "re-source-id")
				err = Mask(err)

				return err
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			s := mustJson(tc.errFunc())

			// Change the given file paths in our Error type's JSON output to
			// avoid prefixes like "/Users/username/go/src/" so this test can be
			// executed on different machines. See the golden files in the
			// testdata folder for specific examples.
			var actual []byte
			{
				p, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				s = strings.ReplaceAll(s, p, "--REPLACED--")

				b := &bytes.Buffer{}
				err = json.Indent(b, []byte(s), "", "\t")
				if err != nil {
					t.Fatal(err)
				}

				actual = []byte(b.String() + "\n")
			}

			p := filepath.Join("testdata/json", fileName(i))
			if *update {
				err := os.WriteFile(p, actual, 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(actual, expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(expected), string(actual)))
			}
		})
	}
}

func Test_Tracer_Error_Stack(t *testing.T) {
	var (
		testErrorOne       = fmt.Errorf("test error one")
		testErrorTwo       = &Error{Kind: "testErrorTwo"}
		testErrorThree     = fmt.Errorf("executing \".github/dependabot.yaml\"")
		alreadyExistsError = &Error{Kind: "alreadyExistsError", Code: "invalidArgument"}
	)

	testCases := []struct {
		errFunc func() error
	}{
		// Case 000 ensures that without error there is no cause.
		{
			errFunc: func() error {
				return nil
			},
		},
		// Case 001 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
		},
		// Case 002 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
		},
		// Case 003 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
		},
		// Case 004 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
		},
		// Case 005 does error wrapping two times. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorOne)
				err = Mask(err)

				return err
			},
		},
		// Case 006 does error wrapping two times. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorTwo)
				err = Mask(err)

				return err
			},
		},
		// Case 007 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
		},
		// Case 008 does error wrapping two times, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				var err error

				err = Maskf(testErrorTwo, "annotation")
				err = Mask(err)

				return err
			},
		},
		// Case 009 does not use error wrapping. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return testErrorThree
			},
		},
		// Case 010 does error wrapping one time. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return Mask(testErrorThree)
			},
		},
		// Case 011 does error wrapping two times. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorThree)
				err = Mask(err)

				return err
			},
		},
		// Case 012
		{
			errFunc: func() error {
				var err error

				err = Mask(alreadyExistsError)
				err = Mask(err)

				return err
			},
		},
		// Case 013
		{
			errFunc: func() error {
				var err error

				err = Maskf(alreadyExistsError, "re-source-id")
				err = Mask(err)

				return err
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			s := Stack(tc.errFunc())

			// Change the given file paths in our Error type's JSON output to
			// avoid prefixes like "/Users/username/go/src/" so this test can be
			// executed on different machines. See the golden files in the
			// testdata folder for specific examples.
			var actual []byte
			{
				p, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				s = strings.ReplaceAll(s, p, "--REPLACED--")

				b := &bytes.Buffer{}
				err = json.Indent(b, []byte(s), "", "\t")
				if err != nil {
					t.Fatal(err)
				}

				actual = []byte(b.String() + "\n")
			}

			p := filepath.Join("testdata/stack", fileName(i))
			if *update {
				err := os.WriteFile(p, actual, 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(actual, expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(expected), string(actual)))
			}
		})
	}
}

func fileName(i int) string {
	return "case-" + fmt.Sprintf("%03d", i) + ".golden"
}
