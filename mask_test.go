package tracer

import (
	"errors"
	"fmt"
	"testing"
)

func Test_Tracer_Mask(t *testing.T) {
	var (
		testErrorOne = fmt.Errorf("testErrorOne")
		testErrorTwo = &Error{Kind: "testErrorTwo"}
	)

	testCases := []struct {
		errFunc func() error
		cause   error
	}{
		// Case 000 ensures that without error there is no cause.
		{
			errFunc: func() error {
				return nil
			},
			cause: nil,
		},
		// Case 001 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
			cause: testErrorOne,
		},
		// Case 002 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
			cause: testErrorTwo,
		},
		// Case 003 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
			cause: testErrorOne,
		},
		// Case 004 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
			cause: testErrorTwo,
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
			cause: testErrorOne,
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
			cause: testErrorTwo,
		},
		// Case 007 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
			cause: testErrorTwo,
		},
		// Case 008 does error wrapping two times, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error type.
		{
			errFunc: func() error {
				var err error

				err = Maskf(testErrorTwo, "annotation")
				err = Mask(err)

				return err
			},
			cause: testErrorTwo,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			cause := tc.errFunc()

			if !errors.Is(cause, tc.cause) {
				t.Fatalf("expected %#v got %#v", tc.cause, cause)
			}
		})
	}
}
