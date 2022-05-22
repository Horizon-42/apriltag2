package apriltag

import (
	"fmt"
	"gocv.io/x/gocv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
)

type Config struct {
	TargetType string  `yaml:"target_type"`
	TagCols    int     `yaml:"tagCols"`
	TagRows    int     `yaml:"tagRows"`
	TagSize    float32 `yaml:"tagSize"`
	TagSpacing float32 `yaml:"tagSpacing"`
}

func NewConfig(configPath string) *Config {
	var config Config
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取AprilTag配置文件失败, %v", err)
	}
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("解析配置文件失败, %v", err)
	}
	return &config
}

type TagPoint struct {
	points   []gocv.Point2f
	points3d []gocv.Point3f
}

func (tp *TagPoint) String2d() string {
	return fmt.Sprintf("%f, %f,\t%f, %f,\t,%f, %f,\t%f, %f\t\t",
		tp.points[0].X, tp.points[0].Y,
		tp.points[1].X, tp.points[1].Y,
		tp.points[2].X, tp.points[2].Y,
		tp.points[3].X, tp.points[3].Y)
}

func (tp *TagPoint) String() string {
	return fmt.Sprintf("%f, %f, %f\t%f, %f, %f\t%f, %f, %f\t%f, %f, %f\t\t",
		tp.points3d[0].X, tp.points3d[0].Y, tp.points3d[0].Z,
		tp.points3d[1].X, tp.points3d[1].Y, tp.points3d[1].Z,
		tp.points3d[2].X, tp.points3d[2].Y, tp.points3d[2].Z,
		tp.points3d[3].X, tp.points3d[3].Y, tp.points3d[3].Z)
}

func NewTagPoint() TagPoint {
	return TagPoint{
		points:   make([]gocv.Point2f, 4),
		points3d: make([]gocv.Point3f, 4),
	}
}

type AprilTag struct {
	Corners []*TagPoint
	Ids     []int

	Origin gocv.Point3f
	P      gocv.Point3f

	config *Config

	detected bool
}

func NewAprilTag(config *Config) *AprilTag {
	if config == nil {
		return nil
	}
	return &AprilTag{
		Corners: make([]*TagPoint, config.TagCols*config.TagRows),
		config:  config,
		Ids:     make([]int, 0),
	}
}

func (a *AprilTag) Empty() bool {
	return !a.detected
}
func (a *AprilTag) Full() bool {
	return len(a.Ids) == a.config.TagCols*a.config.TagRows
}

func (a *AprilTag) SetDetected(detected bool) {
	a.detected = detected
}

func (a *AprilTag) SetPoint3d(id int) {
	// 公式计算3d点坐标
	pt := a.Corners[id]
	if pt == nil {
		return
	}
	row := float32(id / a.config.TagCols)
	col := float32(id % a.config.TagCols)
	pt.points3d[0] = gocv.Point3f{
		X: col*(a.config.TagSize*(1+a.config.TagSpacing)) + a.config.TagSize*a.config.TagSpacing,
		Y: row*(a.config.TagSize*(1+a.config.TagSpacing)) + a.config.TagSize*a.config.TagSpacing,
		Z: 0,
	}

	pt.points3d[1] = gocv.Point3f{
		X: pt.points3d[0].X + a.config.TagSize,
		Y: pt.points3d[0].Y,
		Z: 0,
	}

	pt.points3d[2] = gocv.Point3f{
		X: pt.points3d[1].X,
		Y: pt.points3d[1].Y + a.config.TagSize,
		Z: 0,
	}

	pt.points3d[3] = gocv.Point3f{
		X: pt.points3d[0].X,
		Y: pt.points3d[2].Y,
		Z: 0,
	}
}

func (a *AprilTag) SetAll3dPoints() {
	for _, id := range a.Ids {
		a.SetPoint3d(id)
	}
}

func (a *AprilTag) SetWorld3dPoints() {
	// 计算旋转角度
	p := gocv.Point3f{
		X: a.P.X - a.Origin.X,
		Y: a.P.Y - a.Origin.Y,
		Z: a.P.Z - a.Origin.Z,
	}
	rotateMat := gocv.NewMatWithSizeFromScalar(gocv.Scalar{
		Val1: 0,
		Val2: 0,
		Val3: 0,
		Val4: 0,
	}, 3, 3, gocv.MatTypeCV32FC1)
	rotateMat.SetFloatAt(2, 2, 1)
	pMod := float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
	cosTheta := p.X / pMod
	sinTheta := p.Y / pMod
	rotateMat.SetFloatAt(0, 0, cosTheta)
	rotateMat.SetFloatAt(0, 1, -sinTheta)
	rotateMat.SetFloatAt(1, 0, sinTheta)
	rotateMat.SetFloatAt(1, 1, cosTheta)

	for _, id := range a.Ids {
		a.SetPoint3d(id)
		tp := a.Corners[id]
		if tp == nil {
			continue
		}
		for i := 0; i < 4; i++ {
			point := tp.points3d[i]
			point.X = rotateMat.GetFloatAt(0, 0)*point.X + rotateMat.GetFloatAt(0, 1)*point.Y +
				rotateMat.GetFloatAt(0, 2)*point.Z + a.Origin.X
			point.Y = rotateMat.GetFloatAt(1, 0)*point.X + rotateMat.GetFloatAt(1, 1)*point.Y +
				rotateMat.GetFloatAt(1, 2)*point.Z + a.Origin.Y
			point.Z += a.Origin.Z
			tp.points3d[i] = point
		}
	}
}

func (a *AprilTag) GetCornerPoints() (pts2d []gocv.Point2f, pts3d []gocv.Point3f) {
	pts2d = make([]gocv.Point2f, len(a.Ids)*4)
	pts3d = make([]gocv.Point3f, len(a.Ids)*4)
	for i, id := range a.Ids {
		corner := a.Corners[id]

		for j := 0; j < 4; j++ {
			pts2d[i*4+j] = corner.points[j]
			pts3d[i*4+j] = corner.points3d[j]
		}
	}
	return pts2d, pts3d
}
