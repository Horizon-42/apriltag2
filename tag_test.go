package apriltag

import (
	"log"
	"testing"
)

func TestNewAprilTag(t *testing.T) {
	a := []int{1, 2, 3, 8, 5}
	a[3] += 1
	for _, num := range a {
		log.Printf("%d\t", num)
	}
	println()
	tag := NewAprilTag(NewConfig("april.yaml"))
	if tag != nil {
		log.Printf("%v", tag)
	} else {
		t.Fatalf("failed.")
	}
}
