package main

import (
	"fmt"
	"testing"
)

func TestNfa(t *testing.T) {
	var data = []struct {
		email    string
		validity bool
	}{
		{email: "valid_email@example.com", validity: true},
		{email: "john.doe@email.com", validity: true},
		// ... rest of your data
	}

	t.Log("Starting to parse regex...")
	ctx := parse("[a-zA-Z][a-zA-Z0-9_.]+@[a-zA-Z0-9]+.[a-zA-Z]{2,}")
	t.Log("Regex parsed successfully")

	t.Log("Converting to NFA...")
	nfa := toNfa(ctx)
	t.Log("NFA created successfully")

	for i, instance := range data {
		t.Run(fmt.Sprintf("Test_%d: '%s'", i, instance.email), func(t *testing.T) {
			t.Logf("Testing email: %s (expected: %t)", instance.email, instance.validity)
			result := nfa.check(instance.email, -1)
			t.Logf("Result: %t", result)

			if result != instance.validity {
				t.Errorf("Expected: %t, got: %t", instance.validity, result)
			}
		})
	}
}
