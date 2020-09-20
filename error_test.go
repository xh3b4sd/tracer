package tracer

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Tracer_Error_Error(t *testing.T) {
	var (
		testErrorOne   = fmt.Errorf("test error one")
		testErrorTwo   = &Error{Kind: "testErrorTwo"}
		testErrorThree = fmt.Errorf("executing \".github/dependabot.yaml\"")
	)

	testCases := []struct {
		errFunc func() error
		message string
	}{
		// Case 0 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
			message: "test error one",
		},
		// Case 1 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
			message: "test error two",
		},
		// Case 2 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
			message: "test error one",
		},
		// Case 3 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
			message: "test error two",
		},
		// Case 4 does error wrapping two times. The error is simply made up by
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
		// Case 5 does error wrapping two times. The error is simply made up by
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
		// Case 6 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
			message: "annotation",
		},
		// Case 7 does error wrapping two times, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				var err error

				err = Maskf(testErrorTwo, "annotation")
				err = Mask(err)

				return err
			},
			message: "annotation",
		},
		// Case 8 does not use error wrapping. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return testErrorThree
			},
			message: "executing \".github/dependabot.yaml\"",
		},
		// Case 9 does error wrapping one time. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return Mask(testErrorThree)
			},
			message: "executing \".github/dependabot.yaml\"",
		},
		// Case 10 does error wrapping two times. The error is simply made up by
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
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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
		// Case 0 ensures that two equal errors are detected to be equal. Both
		// errors are of the same arbitrary error type.
		{
			one:   testErrorOne,
			two:   testErrorOne,
			equal: true,
		},
		// Case 1 ensures that two equal errors are detected to be equal.
		{
			one:   testErrorTwo,
			two:   testErrorTwo,
			equal: true,
		},
		// Case 2 ensures that two different errors are not detected to be
		// equal. One error is our Error type. The other other error is some
		// arbitrary error type.
		{
			one:   testErrorTwo,
			two:   testErrorOne,
			equal: false,
		},
		// Case 3 ensures that two different errors are not detected to be
		// equal. Both errors are our Error types, but each is a different
		// instance.
		{
			one:   testErrorTwo,
			two:   testErrorThree,
			equal: false,
		},
		// Case 4 ensures that two different errors are not detected to be
		// equal. Both errors are our Error types. Both errors have the same
		// kind but each is a different instance.
		{
			one:   testErrorTwo,
			two:   testErrorTwoSameKind,
			equal: false,
		},
		// Case 5 ensures that two different errors are not detected to be
		// equal. Both errors are of the same arbitrary error type. Both errors
		// have the same kind but each is a different instance.
		{
			one:   testErrorOne,
			two:   testErrorOneSameKind,
			equal: false,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			equal := errors.Is(tc.one, tc.two)

			if equal != tc.equal {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.equal, equal))
			}
		})
	}
}
