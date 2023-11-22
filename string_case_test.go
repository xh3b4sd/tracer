package tracer

import (
	"fmt"
	"testing"
)

func Test_Tracer_toStringCase(t *testing.T) {
	testCases := []struct {
		InputString    string
		ExpectedString string
	}{
		// Case 000 camel case to string case with lower start
		{
			InputString:    "fooBar",
			ExpectedString: "foo bar",
		},
		// Case 001 camel case to string case with lower start and longer input
		{
			InputString:    "fooBarBazupKick",
			ExpectedString: "foo bar bazup kick",
		},
		// Case 002 camel case to string case with upper start
		{
			InputString:    "FooBar",
			ExpectedString: "foo bar",
		},
		// Case 003 camel case to string case with upper start and longer input
		{
			InputString:    "FooBarBazupKick",
			ExpectedString: "foo bar bazup kick",
		},
		// Case 004 real private error kind
		{
			InputString:    "authenticationError",
			ExpectedString: "authentication error",
		},
		// Case 005 real public error kind
		{
			InputString:    "AuthenticationError",
			ExpectedString: "authentication error",
		},
		// Case 006 camel case with abbreviation at the start
		{
			InputString:    "APINotAvailableError",
			ExpectedString: "api not available error",
		},
		// Case 007 camel case with abbreviation in the middle
		{
			InputString:    "invalidHTTPStatusError",
			ExpectedString: "invalid http status error",
		},
		// Case 008 camel case with abbreviation at the end
		{
			InputString:    "fooBarBAZ",
			ExpectedString: "foo bar baz",
		},
		// Case 009 with version numbers at the start
		{
			InputString:    "v2RouteNotReachable",
			ExpectedString: "v2 route not reachable",
		},
		// Case 010 with version numbers in the middle
		{
			InputString:    "oldV2RouteNotReachable",
			ExpectedString: "old v2 route not reachable",
		},
		// Case 011 with version numbers in the middle
		{
			InputString:    "oldV2RouteNotReachable",
			ExpectedString: "old v2 route not reachable",
		},
		// Case 012 with version numbers at the end does not work
		{
			InputString:    "statusCode200",
			ExpectedString: "status code200",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			output := toStringCase(tc.InputString)
			if output != tc.ExpectedString {
				t.Fatalf("expected %#v got %#v", tc.ExpectedString, output)
			}
		})
	}
}
