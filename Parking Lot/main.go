package main

import (
	"errors"
	"fmt"
	"math"
	"time"

)

type VehicleType string

const (
	VehicleTypeBike VehicleType = "bike"
	VehicleTypeCar  VehicleType = "car"
)

type Vehicle struct {
	RegistrationNumber string
	Type               VehicleType
}

type ParkingSpot struct {
	ID                   string
	SupportedVehicleType VehicleType
	ParkedVehicle        *Vehicle
}

type ParkingFloor struct {
	FloorNumber int
	Spots       []*ParkingSpot
}

type ParkingLot struct {
	ID     int
	Floors []ParkingFloor
}

type Ticket struct {
	TicketNo       int
	VehicleDetails Vehicle
	ParkingSpotID  string
	EntryTime      time.Time
	ExitTime       *time.Time
	Fees           float64
}

type ParkingLotService struct {
	lot         *ParkingLot
	TicketRefNo int
	StoreTicket map[int]*Ticket
}

// Park a vehicle in a compatible, available spot.
func (ps *ParkingSpot) Park(vehicle *Vehicle) error {
	if vehicle == nil {
		return errors.New("vehicle cannot be nil")
	}

	if ps.ParkedVehicle != nil {
		return errors.New("spot is already occupied")
	}

	if ps.SupportedVehicleType != vehicle.Type {
		return errors.New("spot is not compatible with vehicle type")
	}

	ps.ParkedVehicle = vehicle
	return nil
}

// Find an available spot on one floor.
func (f *ParkingFloor) FindAvailableSpot(
	vehicleType VehicleType,
) (*ParkingSpot, error) {
	for _, spot := range f.Spots {
		if spot != nil &&
			spot.SupportedVehicleType == vehicleType &&
			spot.ParkedVehicle == nil {
			return spot, nil
		}
	}

	return nil, errors.New("no compatible spot available on this floor")
}

// Find an available spot across the entire parking lot.
func (pl *ParkingLot) FindSpot(
	vehicleType VehicleType,
) (*ParkingSpot, error) {
	for i := range pl.Floors {
		spot, err := pl.Floors[i].FindAvailableSpot(vehicleType)
		if err != nil {
			continue
		}

		if spot != nil {
			return spot, nil
		}
	}

	return nil, errors.New(
		"no compatible parking spot available in this parking lot",
	)
}

// Find a specific parking spot by its ID.
func (pl *ParkingLot) FindSpotByID(
	spotID string,
) (*ParkingSpot, error) {
	for i := range pl.Floors {
		for _, spot := range pl.Floors[i].Spots {
			if spot != nil && spot.ID == spotID {
				return spot, nil
			}
		}
	}

	return nil, errors.New("parking spot ID not found")
}

// Handle vehicle entry and create a ticket.
func (s *ParkingLotService) EnterVehicle(
	vehicle *Vehicle,
) (*Ticket, error) {
	if vehicle == nil {
		return nil, errors.New("vehicle cannot be nil")
	}

	if s.lot == nil {
		return nil, errors.New("parking lot is not configured")
	}

	// Find an available spot.
	spot, err := s.lot.FindSpot(vehicle.Type)
	if err != nil {
		return nil, err
	}

	// Park the vehicle.
	if err := spot.Park(vehicle); err != nil {
		return nil, err
	}

	// Ensure the ticket map is initialized.
	if s.StoreTicket == nil {
		s.StoreTicket = make(map[int]*Ticket)
	}

	// Generate the next ticket number.
	s.TicketRefNo++

	// Create the ticket.
	ticket := &Ticket{
		TicketNo:       s.TicketRefNo,
		VehicleDetails: *vehicle,
		ParkingSpotID:  spot.ID,
		EntryTime:      time.Now(),
	}

	// Store the ticket for later exit processing.
	s.StoreTicket[ticket.TicketNo] = ticket

	return ticket, nil
}

// Calculate the fee using started-hour billing.
func (s *ParkingLotService) CalculateFees(
	ticket *Ticket,
	exitTime time.Time,
) (float64, error) {
	if ticket == nil {
		return 0, errors.New("ticket cannot be nil")
	}

	duration := exitTime.Sub(ticket.EntryTime)
	if duration < 0 {
		return 0, errors.New("exit time cannot be before entry time")
	}

	// Round up partial hours.
	billableHours := math.Ceil(duration.Hours())

	var hourlyRate float64

	switch ticket.VehicleDetails.Type {
	case VehicleTypeCar:
		hourlyRate = 50
	case VehicleTypeBike:
		hourlyRate = 20
	default:
		return 0, errors.New("unsupported vehicle type")
	}

	return billableHours * hourlyRate, nil
}

// Handle vehicle exit, calculate the fee, and release the spot.
func (s *ParkingLotService) ExitVehicle(
	ticketNo int,
) (float64, error) {
	// 1. Find the ticket.
	ticket, ok := s.StoreTicket[ticketNo]
	if !ok || ticket == nil {
		return 0, errors.New("ticket not found")
	}

	// 2. Reject tickets that have already been completed.
	if ticket.ExitTime != nil {
		return 0, errors.New("vehicle has already exited")
	}

	if s.lot == nil {
		return 0, errors.New("parking lot is not configured")
	}

	// 3. Find the assigned spot.
	spot, err := s.lot.FindSpotByID(ticket.ParkingSpotID)
	if err != nil {
		return 0, err
	}

	// 4. Validate that the assigned spot contains this vehicle.
	if spot.ParkedVehicle == nil {
		return 0, errors.New("assigned parking spot is already empty")
	}

	if spot.ParkedVehicle.RegistrationNumber !=
		ticket.VehicleDetails.RegistrationNumber {
		return 0, errors.New(
			"assigned spot contains a different vehicle",
		)
	}

	// 5. Capture one exit timestamp.
	exitTime := time.Now()

	// 6. Calculate the fee before changing the parking state.
	fee, err := s.CalculateFees(ticket, exitTime)
	if err != nil {
		return 0, err
	}

	// 7. Complete the ticket.
	ticket.ExitTime = &exitTime
	ticket.Fees = fee

	// 8. Release the spot.
	spot.ParkedVehicle = nil

	return fee, nil
}

func main() {
	bikeSpot := &ParkingSpot{
		ID:                   "B1",
		SupportedVehicleType: VehicleTypeBike,
	}

	carSpot := &ParkingSpot{
		ID:                   "C1",
		SupportedVehicleType: VehicleTypeCar,
	}

	floor := ParkingFloor{
		FloorNumber: 1,
		Spots:       []*ParkingSpot{bikeSpot, carSpot},
	}

	lot := ParkingLot{
		ID:     1,
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	// Create a car.
	car := Vehicle{
		RegistrationNumber: "MH01AB1234",
		Type:               VehicleTypeCar,
	}

	// Enter the parking lot.
	ticket, err := service.EnterVehicle(&car)
	if err != nil {
		fmt.Println("Failed to park vehicle:", err)
		return
	}

	fmt.Println("Vehicle entered successfully")
	fmt.Println("Ticket number:", ticket.TicketNo)
	fmt.Println("Assigned spot:", ticket.ParkingSpotID)
	fmt.Println("Entry time:", ticket.EntryTime.Format(time.RFC3339))
	fmt.Println("Spot C1 occupied:", carSpot.ParkedVehicle != nil)

	// Exit the parking lot.
	fee, err := service.ExitVehicle(ticket.TicketNo)
	if err != nil {
		fmt.Println("Failed to exit vehicle:", err)
		return
	}

	fmt.Println("\nVehicle exited successfully")
	fmt.Printf("Parking fee: ₹%.2f\n", fee)
	fmt.Println("Exit time:", ticket.ExitTime.Format(time.RFC3339))
	fmt.Println("Spot C1 available:", carSpot.ParkedVehicle == nil)
}