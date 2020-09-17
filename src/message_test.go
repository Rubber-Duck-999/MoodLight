// message_test.go

package main

import (
	"testing"
)

// Check that State is set
// then run this test will prove it is set
func TestStateSetTrue(t *testing.T) {
	SetState(false)
	if getState() == true {
		t.Error("Failure")
	}
}

// Check that a test email will
// to send as no details have been inputted
func TestSendFailureEmail(t *testing.T) {
	SetState(true)
	if TestEmail() == false {
		t.Error("Failure")
	}
}

// Check that a test email will
// to send as no details have been inputted
func TestSendAttachedFailureEmail(t *testing.T) {
	SetState(true)
	if sendAttachedEmail("Test", "../config.yml-fail") == false {
		t.Error("Failure")
	}
}

// Check that State is set and an Email is attempted
// then run this test will prove it is set
func TestEmailState(t *testing.T) {
	SetState(false)
	if getState() != false {
		t.Error("Failure")
	}
	if TestEmail() != false {
		t.Error("Failure")
	}
}

// Check that State is set and an Email is attempted
// then run this test will prove it is set
func TestEmailStateRoutine(t *testing.T) {
	SetState(false)
	if getState() != false {
		t.Error("Failure")
	}
	if SendEmailRoutine("", "") != false {
		t.Error("Failure")
	}
}

func TestCheckCanSend(t *testing.T) {
	setDate()
	if checkCanSend() == false {
		t.Error("Failure")
	}
}
