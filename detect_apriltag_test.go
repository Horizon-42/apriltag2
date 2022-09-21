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
			a.Origin = gocv.Point3f{
				X: 1,
				Y: 1,
				Z: 1,
			}
			a.P = gocv.Point3f{
				X: 1.3,
				Y: 1.4,
				Z: 1,
			}
			if !a.Empty() {
				a.SetAll3dPoints()
				log.Printf("%v", a.Corners)
				a.SetWorld3dPoints()
				log.Printf("%v", a.Corners)
			}
			win.IMShow(frame)
			if key := win.WaitKey(10); key == 27 {
				break
			}
		}
	}
}
