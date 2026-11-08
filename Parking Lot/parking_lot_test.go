package main

import (
	"testing"
	"time"

)

// func TestEnterVehicle(t *testing.T) {

// 	tests := []struct {
// 		name       string
// 		vehicle    *Vehicle
// 		spot       *ParkingSpot
// 		wantErr    bool
// 		wantSpotID string
// 	}{
// 		{
// 			name: "successful car entry",
// 			vehicle: &Vehicle{
// 				RegistrationNumber: "MH01AB1234",
// 				Type:               VehicleTypeCar,
// 			},
// 			spot: &ParkingSpot{
// 				ID:                   "C1",
// 				SupportedVehicleType: VehicleTypeCar,
// 			},
// 			wantErr:    false,
// 			wantSpotID: "C1",
// 		},
// 		{
// 			name:    "nil vehicle",
// 			vehicle: nil,
// 			spot: &ParkingSpot{
// 				ID:                   "C1",
// 				SupportedVehicleType: VehicleTypeCar,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "no compatible parking spot",
// 			vehicle: &Vehicle{
// 				RegistrationNumber: "MH01AB1234",
// 				Type:               VehicleTypeCar,
// 			},
// 			spot: &ParkingSpot{
// 				ID:                   "B1",
// 				SupportedVehicleType: VehicleTypeBike,
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {

// 		t.Run(tt.name, func(t *testing.T) {

// 			// ARRANGE: Create a fresh parking lot for each test.

// 			floor := ParkingFloor{
// 				FloorNumber: 1,
// 				Spots:       []*ParkingSpot{tt.spot},
// 			}

// 			lot := ParkingLot{
// 				Floors: []ParkingFloor{floor},
// 			}

// 			service := ParkingLotService{
// 				lot:         &lot,
// 				TicketRefNo: 0,
// 				StoreTicket: make(map[int]*Ticket),
// 			}

// 			// ACT: Execute EnterVehicle.

// 			ticket, err := service.EnterVehicle(tt.vehicle)

// 			// ASSERT: Verify the result.

// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("expected error: %v, got: %v",
// 					tt.wantErr, err)
// 			}

// 			// If an error was expected, no further success checks are needed.
// 			if tt.wantErr {
// 				if ticket != nil {
// 					t.Error("expected nil ticket when entry fails")
// 				}
// 				return
// 			}

// 			if ticket == nil {
// 				t.Fatal("expected ticket, got nil")
// 			}

// 			if ticket.ParkingSpotID != tt.wantSpotID {
// 				t.Errorf("expected spot %s, got %s",
// 					tt.wantSpotID, ticket.ParkingSpotID)
// 			}

// 			if tt.spot.ParkedVehicle == nil {
// 				t.Error("expected parking spot to be occupied")
// 			}

// 			if ticket.TicketNo != 1 {
// 				t.Errorf("expected ticket number 1, got %d",
// 					ticket.TicketNo)
// 			}

// 			// Verify ticket storage.
// 			storedTicket, ok := service.StoreTicket[ticket.TicketNo]

// 			if !ok {
// 				t.Error("ticket was not stored in StoreTicket")
// 			} else if storedTicket != ticket {
// 				t.Error("stored ticket does not match returned ticket")
// 			}
// 		})
// 	}
// }

func TestEnterVehicle_Success(t *testing.T) {

	// ARRANGE: Prepare test data

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
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	car := Vehicle{
		RegistrationNumber: "MH01AB1234",
		Type:               VehicleTypeCar,
	}

	// ACT: Execute the function being tested

	ticket, err := service.EnterVehicle(&car)

	// ASSERT: Check the results

	// 1. Entry should succeed without an error.
	if err != nil {
		t.Fatalf("expected successful entry, got error: %v", err)
	}

	// 2. A ticket should be generated.
	if ticket == nil {
		t.Fatal("expected a ticket, got nil")
	}

	// 3. Correct parking spot should be assigned.
	if ticket.ParkingSpotID != "C1" {
		t.Errorf("expected parking spot C1, got %s", ticket.ParkingSpotID)
	}

	// 4. Vehicle should be parked in the assigned spot.
	if carSpot.ParkedVehicle == nil {
		t.Fatal("expected C1 to be occupied, but it is empty")
	}

	if carSpot.ParkedVehicle.RegistrationNumber != car.RegistrationNumber {
		t.Errorf("expected vehicle %s, got %s",
			car.RegistrationNumber,
			carSpot.ParkedVehicle.RegistrationNumber)
	}

	// 5. Generated ticket should be stored.
	if _, ok := service.StoreTicket[ticket.TicketNo]; !ok {
		t.Error("expected ticket to be stored in StoreTicket")
	}
}

func TestEnterVehicle_NilVehicle(t *testing.T) {

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
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	// ACT
	ticket, err := service.EnterVehicle(nil)

	// ASSERT
	if err == nil {
		t.Fatal("expected error for nil vehicle, got nil")
	}

	if ticket != nil {
		t.Fatal("expected nil ticket for nil vehicle, got non-nil")
	}

	if carSpot.ParkedVehicle != nil || bikeSpot.ParkedVehicle != nil {
		t.Fatal("expected all parking spots to remain empty")
	}
}

func TestEnterVehicle_NoAvailableSpot(t *testing.T) {

	bikeSpot := &ParkingSpot{
		ID:                   "B1",
		SupportedVehicleType: VehicleTypeBike,
	}

	carSpot := &ParkingSpot{
		ID:                   "C1",
		SupportedVehicleType: VehicleTypeCar,
	}

	// Arrange: Occupy the car spot.
	parkedCar := &Vehicle{
		RegistrationNumber: "MH02CD5678",
		Type:               VehicleTypeCar,
	}

	if err := carSpot.Park(parkedCar); err != nil {
		t.Fatalf("failed to arrange occupied spot: %v", err)
	}

	floor := ParkingFloor{
		FloorNumber: 1,
		Spots:       []*ParkingSpot{bikeSpot, carSpot},
	}

	lot := ParkingLot{
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	// ACT
	newCar := Vehicle{
		RegistrationNumber: "MH01AB1234",
		Type:               VehicleTypeCar,
	}

	ticket, err := service.EnterVehicle(&newCar)

	// ASSERT
	if err == nil {
		t.Fatal("expected error for no available spot, got nil")
	}

	if ticket != nil {
		t.Fatal("expected nil ticket, got non-nil")
	}

	if carSpot.ParkedVehicle.RegistrationNumber != parkedCar.RegistrationNumber {
		t.Error("existing parked vehicle should remain unchanged")
	}

	if bikeSpot.ParkedVehicle != nil {
		t.Error("bike spot should remain empty")
	}

	if len(service.StoreTicket) != 0 {
		t.Error("no ticket should be stored when entry fails")
	}
}

func TestExitVehicle_Success(t *testing.T) {

	// ARRANGE: Set up parking lot and service.

	carSpot := &ParkingSpot{
		ID:                   "C1",
		SupportedVehicleType: VehicleTypeCar,
	}

	floor := ParkingFloor{
		FloorNumber: 1,
		Spots:       []*ParkingSpot{carSpot},
	}

	lot := ParkingLot{
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	// Park a car and generate an active ticket.
	car := &Vehicle{
		RegistrationNumber: "MH01AB1234",
		Type:               VehicleTypeCar,
	}

	ticket, err := service.EnterVehicle(car)
	if err != nil {
		t.Fatalf("failed to park vehicle: %v", err)
	}
	// Set entry time one minute in the past.
	// This ensures one billable hour at ₹50/hour.
	ticket.EntryTime = time.Now().Add(-time.Minute)

	// ACT: Exit the vehicle.
	fee, err := service.ExitVehicle(ticket.TicketNo)

	// ASSERT: Verify the results.

	// 1. Exit should succeed without an error.
	if err != nil {
		t.Fatalf("expected successful exit, got error: %v", err)
	}

	// 2. Verify the parking fee.
	if fee != 50 {
		t.Errorf("expected fee 50, got %f", fee)
	}

	// 3. Verify that the parking spot is released.
	if carSpot.ParkedVehicle != nil {
		t.Error("expected parking spot to be empty after exit")
	}

	// 4. Verify that exit time is recorded.
	if ticket.ExitTime == nil {
		t.Error("expected ExitTime to be recorded")
	}

	// 5. Verify that the ticket fee is updated.
	if ticket.Fees != fee {
		t.Errorf("expected ticket fee %f, got %f", fee, ticket.Fees)
	}

	// 6. Verify that the completed ticket remains stored.
	if _, ok := service.StoreTicket[ticket.TicketNo]; !ok {
		t.Error("expected completed ticket to remain in StoreTicket")
	}
}

func TestExitVehicle_InvalidTicket(t *testing.T) {

	// ARRANGE: Set up parking lot and service.

	carSpot := &ParkingSpot{
		ID:                   "C1",
		SupportedVehicleType: VehicleTypeCar,
	}

	floor := ParkingFloor{
		FloorNumber: 1,
		Spots:       []*ParkingSpot{carSpot},
	}

	lot := ParkingLot{
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	// Park a car and generate a valid active ticket.
	car := &Vehicle{
		RegistrationNumber: "MH01AB1234",
		Type:               VehicleTypeCar,
	}

	ticket, err := service.EnterVehicle(car)
	if err != nil {
		t.Fatalf("failed to park vehicle: %v", err)
	}

	// ACT: Attempt to exit with an invalid ticket number.
	invalidTicketNo := 999

	fee, err := service.ExitVehicle(invalidTicketNo)

	// ASSERT: Verify the results.

	// 1. An error should be returned.
	if err == nil {
		t.Fatal("expected error for invalid ticket number, got nil")
	}

	// 2. No parking fee should be returned.
	if fee != 0 {
		t.Errorf("expected fee 0 for invalid ticket, got %f", fee)
	}

	// 3. Original parking spot should remain occupied.
	if carSpot.ParkedVehicle == nil {
		t.Fatal("expected parking spot to remain occupied after failed exit")
	}

	// 4. Original ticket should remain active.
	if ticket.ExitTime != nil {
		t.Error("expected original ticket to remain active")
	}

	// 5. Original ticket should remain stored.
	if _, ok := service.StoreTicket[ticket.TicketNo]; !ok {
		t.Error("expected original ticket to remain in StoreTicket")
	}
}

func TestExitVehicle_AlreadyExited(t *testing.T) {

	// ARRANGE: Set up parking lot and service.

	carSpot := &ParkingSpot{
		ID:                   "C1",
		SupportedVehicleType: VehicleTypeCar,
	}

	floor := ParkingFloor{
		FloorNumber: 1,
		Spots:       []*ParkingSpot{carSpot},
	}

	lot := ParkingLot{
		Floors: []ParkingFloor{floor},
	}

	service := ParkingLotService{
		lot:         &lot,
		TicketRefNo: 0,
		StoreTicket: make(map[int]*Ticket),
	}

	// Park a car and generate a valid active ticket.
	car := &Vehicle{
		RegistrationNumber: "MH01AB1234",
		Type:               VehicleTypeCar,
	}

	ticket, err := service.EnterVehicle(car)
	if err != nil {
		t.Fatalf("failed to park vehicle: %v", err)
	}

	// ACT: First exit should succeed.
	_, err = service.ExitVehicle(ticket.TicketNo)
	if err != nil {
		t.Fatalf("failed to exit vehicle: %v", err)
	}

	// Save the original exit time.
	if ticket.ExitTime == nil {
		t.Fatal("expected ExitTime to be recorded after first exit")
	}

	firstExitTime := *ticket.ExitTime

	// Attempt to exit again using the same ticket.
	fee, err := service.ExitVehicle(ticket.TicketNo)

	// ASSERT: Verify the results.

	// 1. An error should be returned.
	if err == nil {
		t.Fatal("expected error for already exited ticket, got nil")
	}

	// 2. No parking fee should be returned.
	if fee != 0 {
		t.Errorf("expected fee 0 for already exited ticket, got %f", fee)
	}

	// 3. Parking spot should remain empty.
	if carSpot.ParkedVehicle != nil {
		t.Error("expected parking spot to remain empty after already exited ticket")
	}

	// 4. Exit time should remain unchanged.
	if ticket.ExitTime == nil {
		t.Fatal("expected ExitTime to remain set")
	}

	if !ticket.ExitTime.Equal(firstExitTime) {
		t.Error("ExitTime should not change after a repeated exit")
	}

	// 5. Ticket should remain stored.
	if _, ok := service.StoreTicket[ticket.TicketNo]; !ok {
		t.Error("expected ticket to remain in StoreTicket after already exited")
	}
}

func TestCalculateFees_Car(t *testing.T) {

	// ARRANGE: Create a ticket with a specific entry time.

	ticket := &Ticket{
		TicketNo:  1,
		EntryTime: time.Now().Add(-2 * time.Hour), // 2 hours ago
	}

	service := ParkingLotService{
		StoreTicket: make(map[int]*Ticket),
	}

	// ACT: Calculate fees for a car.
	fee, err := service.CalculateFees(ticket, time.Now())

	// ASSERT: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if fee != 100 {
		t.Errorf("expected fee 100 for 2 hours, got: %f", fee)
	}
}