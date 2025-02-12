package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func TestCheckPermissions(t *testing.T) {
	allowedPosition := database.Position{
		ID:                uuid.New(),
		Title:             "Allowed Position",
		CanActivate:       true,
		CanDeactivate:     true,
		CanAddOrders:      true,
		CanReceiveOrders:  true,
		CanQuery:          true,
		CanChangeProtocol: true,
		CanAddStaff:       true,
		CanAddReminders:   true,
	}
	deniedPosition := database.Position{
		ID:                uuid.New(),
		Title:             "Denied Position",
		CanActivate:       false,
		CanDeactivate:     false,
		CanAddOrders:      false,
		CanReceiveOrders:  false,
		CanQuery:          false,
		CanChangeProtocol: false,
		CanAddStaff:       false,
		CanAddReminders:   false,
	}

	permissions := []Permission{
		PermissionActivateInactivate,
		PermissionDeactivateReactivate,
		PermissionAddOrder,
		PermissionReceiveOrder,
		PermissionRunQueries,
		PermissionProtocol,
		PermissionStaff,
		PermissionReminders,
	}

	// should all be allowed ie err == nil
	for i, p := range permissions {
		err := checkPermission(&allowedPosition, p)
		if err != nil {
			fmt.Println(err)
			t.Fatalf("Allowed failed test %v", i)
		}
	}

	// should all be denied ie err != nil
	for i, p := range permissions {
		err := checkPermission(&deniedPosition, p)
		// fmt.Println(err)
		if err == nil {
			t.Fatalf("Denied failed test %v", i)
		}
	}
}
