package tracer

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func Test_Tracer_Cause_Mask(t *testing.T) {
	var (
		testErrorOne = fmt.Errorf("testErrorOne")
		testErrorTwo = &Error{Kind: "testErrorTwo"}
	)

	testCases := []struct {
		errFunc func() error
		cause   error
	}{
		// Case 0 ensures that without error there is no cause.
		{
			errFunc: func() error {
				return nil
			},
			cause: nil,
		},
		// Case 1 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
			cause: testErrorOne,
		},
		// Case 2 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
			cause: testErrorTwo,
		},
		// Case 3 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
			cause: testErrorOne,
		},
		// Case 4 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
			cause: testErrorTwo,
		},
		// Case 5 does error wrapping two times. The error is simply made up by
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
		// Case 6 does error wrapping two times. The error is simply made up by
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
		// Case 7 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
			cause: testErrorTwo,
		},
		// Case 8 does error wrapping two times, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cause := Cause(tc.errFunc())

			if !errors.Is(cause, tc.cause) {
				t.Fatalf("expected %#v got %#v", tc.cause, cause)
			}
		})
	}
}
