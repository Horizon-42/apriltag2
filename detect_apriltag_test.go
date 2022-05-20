package apriltag

import (
	"gocv.io/x/gocv"
	"log"
	"testing"
)

func TestTagDetector_DetectAprilTags(t *testing.T) {
	capture, err := gocv.OpenVideoCapture(0)
	defer capture.Close()
	if err != nil {
		t.Fatal(err)
	}
	frame := gocv.NewMat()
	td := NewTagDetector("C:/works/calibrate_instrinsic/april.yaml")
	defer DestoryTagDetector(td)
	win := gocv.NewWindow("frame")
	defer win.Close()
	for {
		if ok := capture.Read(&frame); ok {
			a := td.DetectAprilTags(frame, true)
			if !a.Empty() {
				log.Printf("%v", a.Corners)
			}
			win.IMShow(frame)
			if key := win.WaitKey(10); key == 27 {
				break
			}
		}
	}
}
