package apriltag

import (
	"log"
	"testing"
)

func TestNewAprilTag(t *testing.T) {
	tag := NewAprilTag(NewConfig("C:/works/calibrate_extrinsics_go/april.yaml"))
	if tag != nil {
		log.Printf("%v", tag)
	} else {
		t.Fatalf("failed.")
	}
}
