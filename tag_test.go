package apriltag

import (
	"log"
	"testing"
)

func TestNewAprilTag(t *testing.T) {
	tag := NewAprilTag(NewConfig("april.yaml"))
	if tag != nil {
		log.Printf("%v", tag)
	} else {
		t.Fatalf("failed.")
	}
}
