package main

import (
	"bytes"
	"os"
	"testing"
)

// captureOutput captures stdout output for testing
func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout
	buf.ReadFrom(r)
	return buf.String()
}

func TestCreateParkingLot(t *testing.T) {
	parkingLot := NewParkingLot(0)
	parkingLot.CreateParkingLot(6)

	if parkingLot.capacity != 6 {
		t.Errorf("Expected capacity 6, got %d", parkingLot.capacity)
	}
	if len(parkingLot.slots) != 0 {
		t.Errorf("Expected empty slots map, got %d slots", len(parkingLot.slots))
	}
}

func TestPark(t *testing.T) {
	parkingLot := NewParkingLot(3)
	parkingLot.CreateParkingLot(3)

	tests := []struct {
		regNo        string
		expectedSlot int
		expectedOut  string
	}{
		{"KA-01-HH-1234", 1, "Allocated slot number: 1\n"},
		{"KA-01-HH-9999", 2, "Allocated slot number: 2\n"},
		{"KA-01-BB-0001", 3, "Allocated slot number: 3\n"},
		{"KA-01-HH-7777", 0, "Sorry, parking lot is full\n"},
	}

	for _, test := range tests {
		output := captureOutput(func() {
			parkingLot.Park(test.regNo)
		})

		if output != test.expectedOut {
			t.Errorf("For regNo %s, expected output %q, got %q", test.regNo, test.expectedOut, output)
		}

		if test.expectedSlot != 0 {
			if car, exists := parkingLot.slots[test.expectedSlot]; !exists || car.RegistrationNo != test.regNo {
				t.Errorf("For regNo %s, expected slot %d with car, got none or different car", test.regNo, test.expectedSlot)
			}
		}
	}
}

func TestLeave(t *testing.T) {
	tests := []struct {
		name        string
		regNo       string
		hours       int
		expectedOut string
		shouldExist map[int]string // Expected slots and regNos after Leave
	}{
		{
			name:        "Leave existing car",
			regNo:       "KA-01-HH-1234",
			hours:       4,
			expectedOut: "Registration number KA-01-HH-1234 with Slot Number 1 is free with Charge $30\n",
			shouldExist: map[int]string{2: "KA-01-HH-9999"},
		},
		{
			name:        "Leave another existing car",
			regNo:       "KA-01-HH-9999",
			hours:       2,
			expectedOut: "Registration number KA-01-HH-9999 with Slot Number 2 is free with Charge $10\n",
			shouldExist: map[int]string{1: "KA-01-HH-1234"},
		},
		{
			name:        "Leave non-existent car",
			regNo:       "KA-01-XX-9999",
			hours:       3,
			expectedOut: "Registration number KA-01-XX-9999 not found\n",
			shouldExist: map[int]string{1: "KA-01-HH-1234", 2: "KA-01-HH-9999"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parkingLot := NewParkingLot(3)
			parkingLot.CreateParkingLot(3)
			// Suppress Park output
			captureOutput(func() { parkingLot.Park("KA-01-HH-1234") }) // Slot 1
			captureOutput(func() { parkingLot.Park("KA-01-HH-9999") }) // Slot 2

			output := captureOutput(func() {
				parkingLot.Leave(test.regNo, test.hours)
			})

			if output != test.expectedOut {
				t.Errorf("For regNo %s, hours %d, expected output %q, got %q", test.regNo, test.hours, test.expectedOut, output)
			}

			// Check remaining slots
			for slot, expectedRegNo := range test.shouldExist {
				if car, exists := parkingLot.slots[slot]; !exists || car.RegistrationNo != expectedRegNo {
					t.Errorf("Expected slot %d to contain %s, got %v", slot, expectedRegNo, car)
				}
			}
			// Ensure no unexpected slots are occupied
			for slot, car := range parkingLot.slots {
				if _, expected := test.shouldExist[slot]; !expected {
					t.Errorf("Unexpected car %s in slot %d", car.RegistrationNo, slot)
				}
			}
		})
	}
}

func TestCalculateCharge(t *testing.T) {
	tests := []struct {
		hours    int
		expected int
	}{
		{1, 10},
		{2, 10},
		{3, 20},
		{4, 30},
		{5, 40},
	}

	for _, test := range tests {
		result := calculateCharge(test.hours)
		if result != test.expected {
			t.Errorf("For %d hours, expected charge $%d, got $%d", test.hours, test.expected, result)
		}
	}
}

func TestStatus(t *testing.T) {
	parkingLot := NewParkingLot(3)
	parkingLot.CreateParkingLot(3)
	// Suppress Park output
	captureOutput(func() { parkingLot.Park("KA-01-HH-1234") }) // Slot 1
	captureOutput(func() { parkingLot.Park("KA-01-HH-9999") }) // Slot 2

	output := captureOutput(func() {
		parkingLot.Status()
	})

	expected := "Slot No. Registration No.\n1 KA-01-HH-1234\n2 KA-01-HH-9999\n"
	if output != expected {
		t.Errorf("Expected status output:\n%s\nGot:\n%s", expected, output)
	}
}

func TestStatusEmpty(t *testing.T) {
	parkingLot := NewParkingLot(3)
	parkingLot.CreateParkingLot(3)

	output := captureOutput(func() {
		parkingLot.Status()
	})

	if output != "" {
		t.Errorf("Expected empty status output, got %q", output)
	}
}