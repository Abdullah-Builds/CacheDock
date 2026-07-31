package cache

import (
	"testing"
	"time"
)

func TestCleanup(t *testing.T) {
	c := New()
	c.Set("short-lived", "value", 15*time.Millisecond)
	time.Sleep(25 * time.Millisecond)
	c.cleanupExpired()

	if _, ok := c.Get("short-lived"); ok {
		t.Fatal("expected cleanup to remove expired entry")
	}
}
