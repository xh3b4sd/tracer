package tracer

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Tracer_JSON tests the masking behaviour based on our Error type's JSON
// output. The tests use golden file references. In case the golden files change
// something is broken. In case intentional changes get introduced the golden
// files have to be updated. In case the golden files have to be adjusted,
// simply provide the -update flag when running the tests.
//
//     go test . -run Test_Tracer_JSON -update
//
func Test_Tracer_JSON(t *testing.T) {
	var (
		testErrorOne   = fmt.Errorf("test error one")
		testErrorTwo   = &Error{Kind: "testErrorTwo"}
		testErrorThree = fmt.Errorf("executing \".github/dependabot.yaml\"")
	)

	testCases := []struct {
		errFunc func() error
	}{
		// Case 0 ensures that without error there is no cause.
		{
			errFunc: func() error {
				return nil
			},
		},
		// Case 1 does not use error wrapping. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return testErrorOne
			},
		},
		// Case 2 does not use error wrapping. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return testErrorTwo
			},
		},
		// Case 3 does error wrapping one time. The error is simply made up by
		// fmt.Errorf.
		{
			errFunc: func() error {
				return Mask(testErrorOne)
			},
		},
		// Case 4 does error wrapping one time. The error is simply made up by
		// using the tracer error type.
		{
			errFunc: func() error {
				return Mask(testErrorTwo)
			},
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
		},
		// Case 7 does error wrapping one time, while the first wrapping is
		// annotated. The error is simply made up by using the tracer error
		// type.
		{
			errFunc: func() error {
				return Maskf(testErrorTwo, "annotation")
			},
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
		},
		// Case 9 does not use error wrapping. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return testErrorThree
			},
		},
		// Case 10 does error wrapping one time. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				return Mask(testErrorThree)
			},
		},
		// Case 11 does error wrapping two times. The error is simply made up by
		// fmt.Errorf. The error message contains escaped double quotes.
		{
			errFunc: func() error {
				var err error

				err = Mask(testErrorThree)
				err = Mask(err)

				return err
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s := JSON(tc.errFunc())

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
				err := ioutil.WriteFile(p, actual, 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := ioutil.ReadFile(p)
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
	return "case-" + strconv.Itoa(i) + ".golden"
}
