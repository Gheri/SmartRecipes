package main

import "testing"

func TestDeliveryTimeQuery(t *testing.T) {
	t.Run("Test if 3AM - 10PM is within 2AM - 11PM", func(t *testing.T) {
		isInRange := isTimeRangeWithIn("3", "10", "2", "11")
		if !isInRange {
			t.Error("3AM - 10PM should be within 2AM - 11PM")
		}
	})
	t.Run("Test if 2AM - 11PM is within 2AM - 11PM", func(t *testing.T) {
		isInRange := isTimeRangeWithIn("2", "11", "2", "11")
		if !isInRange {
			t.Error("2AM - 11PM should be within 2AM - 11PM")
		}
	})
	t.Run("Test if 1AM - 11PM is not within 2AM - 11PM", func(t *testing.T) {
		isInRange := isTimeRangeWithIn("1", "11", "2", "11")
		if isInRange {
			t.Error("1AM - 11PM should not be within 2AM - 11PM")
		}
	})
}
