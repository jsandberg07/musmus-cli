package main

// change if you want a million things printed or not
const verbose bool = false

type Permission int

const (
	PermissionActivateInactivate = iota
	PermissionDeactivateReactivate
	PermissionAddOrder
	PermissionReceiveOrder
	PermissionRunQueries
	PermissionProtocol
	PermissionStaff
	PermissionReminders
)

const DefaultFlagMsg string = "An allowed flag was unhandled by the switch: "
const CancelMsg string = "Exiting..."
const CancelError string = "cancel"
