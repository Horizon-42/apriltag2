package apriltag

/*
//#cgo CPPFLAGS: -I /usr/local/include -I /usr/local/include/opencv4 -I /usr/local/include/apriltags -I /usr/include/eigen3
//#cgo LDFLAGS: -L /usr/local/lib -L /usr/local/opencv4/lib -L /home/horizon/works/calibrate_extrinsics_go/apriltag -lopencv_core -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_dnn -lethz_apriltag2 -lopencv_cudev -lopencv_cudafilters
#include <stdlib.h>
#include "detect_apriltag.h"
*/
import "C"
import (
	//"fmt"
	"gocv.io/x/gocv"
)

type Mat struct {
	p C.Mat
}

func (m *Mat) MatSetP(mat gocv.Mat) {
	m.p = C.Mat(mat.Ptr())
}

func (m *Mat) Ptr() C.Mat {
	return m.p
}

type TagDetector struct {
	detecorPtr C.TagDetectorPtr
	config     *Config
}

func (td *TagDetector) CountAprilTags(frame gocv.Mat) int {
	src := Mat{}
	src.MatSetP(frame)
	return int(C.CountTags(td.detecorPtr, src.Ptr()))
}

func (td *TagDetector) DetectAprilTags(frame gocv.Mat, draw bool) *AprilTag {
	src := Mat{}
	src.MatSetP(frame)
	points := gocv.NewMat()
	ptsIds := gocv.NewMat()

	ptsOut := Mat{}
	ptsOut.MatSetP(points)
	idsOut := Mat{}
	idsOut.MatSetP(ptsIds)

	ret := bool(C.DetectTags(td.detecorPtr, src.Ptr(), ptsOut.Ptr(), idsOut.Ptr(), C._Bool(draw)))

	tag := NewAprilTag(td.config)
	tag.SetDetected(ret)
	if ret {
		pointsNum := ptsIds.Rows()
		for i := 0; i < pointsNum; i++ {
			tagPoint := NewTagPoint()
			for j := 0; j < 4; j++ {
				tagPoint.points[j] = gocv.Point2f{
					X: points.GetFloatAt(i*4+j, 0),
					Y: points.GetFloatAt(i*4+j, 1),
				}
			}
			id := int(ptsIds.GetIntAt(i, 0))
			tag.Corners[id] = &tagPoint
			tag.Ids = append(tag.Ids, id)
		}
	}
	return tag
}

func (td *TagDetector) IsEmpty() bool {
	return bool(C.IsEmpty(td.detecorPtr))
}

func NewTagDetector(configPath string) *TagDetector {
	return &TagDetector{detecorPtr: C.NewTagDetector(), config: NewConfig(configPath)}
}

func DestoryTagDetector(detector *TagDetector) {
	C.ReleaseTagDetector(&detector.detecorPtr)
}
