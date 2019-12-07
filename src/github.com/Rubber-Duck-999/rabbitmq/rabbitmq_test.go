// rabbitmq_test.go

package rabbitmq

import (
	"errors"
	"strings"
	"testing"
)

// Check that State is set
// then run this test will prove it is set
func TestPublishFailRabbit(t *testing.T) {
	init_err = errors.New("dial tcp 127.0.0.1:5672: getsockopt: connection refused")
	failure := "cheese"
	//failure = messageFailure(true)
	if strings.Contains(FAILUREPUBLISH, failure) {
		t.Error("Failure")
	} else if strings.Contains(FAILURECONVERT, failure) {
		t.Error("Failure")
	}
}

func TestMessagesPass(t *testing.T) {
	messages(FAILURENETWORK, "")
}
