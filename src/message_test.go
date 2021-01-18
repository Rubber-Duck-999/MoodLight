// message_test.go

package main

import (
	"testing"
)

// Check that a test email will
// to send as no details have been inputted
func TestSendFailureEmail(t *testing.T) {
	if TestEmail() == false {
		t.Error("Failure")
	}
}

// Check that State is set and an Email is attempted
// then run this test will prove it is set
func TestEmailState(t *testing.T) {
	if TestEmail() == false {
		t.Error("Failure")
	}
}

func TestCheckCanSend(t *testing.T) {
	setDate()
	if checkCanSend() == false {
		t.Error("Failure")
	}
}
