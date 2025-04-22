package main

// This program simulates a parking lot system
// Using Go's standard library to read commands from a file
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Car represents a parked car
type Car struct {
	RegistrationNo string
	Slot           int
}

// ParkingLot manages the parking slots
type ParkingLot struct {
	capacity int
	slots    map[int]*Car // Map of slot number to car
}

// NewParkingLot creates a new parking lot with given capacity
func NewParkingLot(capacity int) *ParkingLot {
	return &ParkingLot{
		capacity: capacity,
		slots:    make(map[int]*Car),
	}
}

// CreateParkingLot initializes the parking lot
func (p *ParkingLot) CreateParkingLot(capacity int) {
	p.capacity = capacity
	p.slots = make(map[int]*Car)
}

// Park allocates a slot to a car
func (p *ParkingLot) Park(regNo string) {
	if len(p.slots) >= p.capacity {
		fmt.Println("Sorry, parking lot is full")
		return
	}

	// Find the nearest available slot
	for i := 1; i <= p.capacity; i++ {
		if _, occupied := p.slots[i]; !occupied {
			p.slots[i] = &Car{RegistrationNo: regNo, Slot: i}
			fmt.Printf("Allocated slot number: %d\n", i)
			return
		}
	}
}

// Leave removes a car and calculates charges
func (p *ParkingLot) Leave(regNo string, hours int) {
	for slot, car := range p.slots {
		if car.RegistrationNo == regNo {
			charge := calculateCharge(hours)
			fmt.Printf("Registration number %s with Slot Number %d is free with Charge $%d\n",
				regNo, slot, charge)
			delete(p.slots, slot)
			return
		}
	}
	fmt.Printf("Registration number %s not found\n", regNo)
}

// Status prints the current parking lot status
func (p *ParkingLot) Status() {
	if len(p.slots) == 0 {
		return
	}

	fmt.Println("Slot No. Registration No.")
	for slot := 1; slot <= p.capacity; slot++ {
		if car, exists := p.slots[slot]; exists {
			fmt.Printf("%d %s\n", slot, car.RegistrationNo)
		}
	}
}

// calculateCharge computes parking fees
func calculateCharge(hours int) int {
	if hours <= 2 {
		return 10
	}
	return 10 + (hours-2)*10
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run parking_lot.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var parkingLot *ParkingLot
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "create_parking_lot":
			if len(parts) != 2 {
				fmt.Println("Invalid create_parking_lot command")
				continue
			}
			capacity, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid capacity")
				continue
			}
			parkingLot = NewParkingLot(capacity)
			parkingLot.CreateParkingLot(capacity)

		case "park":
			if len(parts) != 2 {
				fmt.Println("Invalid park command")
				continue
			}
			if parkingLot == nil {
				fmt.Println("Parking lot not initialized")
				continue
			}
			parkingLot.Park(parts[1])

		case "leave":
			if len(parts) != 3 {
				fmt.Println("Invalid leave command")
				continue
			}
			if parkingLot == nil {
				fmt.Println("Parking lot not initialized")
				continue
			}
			hours, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid hours")
				continue
			}
			parkingLot.Leave(parts[1], hours)

		case "status":
			if parkingLot == nil {
				fmt.Println("Parking lot not initialized")
				continue
			}
			parkingLot.Status()

		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
}