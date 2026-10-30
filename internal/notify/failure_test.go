package notify

import (
	"errors"
	"net/url"
	"strings"
	"testing"
)

func TestFailureOmitsToken(t *testing.T) {
	err := &url.Error{
		Op:  "Post",
		URL: "https://api.telegram.org/botSECRET/sendMessage",
		Err: errors.New("timeout"),
	}
	got := failure(err)
	if strings.Contains(got, "SECRET") || strings.Contains(got, "api.telegram.org") {
		t.Fatalf("logged url: %s", got)
	}
	if got != "timeout" {
		t.Fatalf("got %s", got)
	}
}
