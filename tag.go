package apriltag

import (
	"fmt"
	"gocv.io/x/gocv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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

func (tp *TagPoint) String() string {
	return fmt.Sprintf("%f, %f,\t%f, %f,\t,%f, %f,\t%f, %f\t\t",
		tp.points[0].X, tp.points[0].Y,
		tp.points[1].X, tp.points[1].Y,
		tp.points[2].X, tp.points[2].Y,
		tp.points[3].X, tp.points[3].Y)
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

func (a *AprilTag) SetDetected(detected bool) {
	a.detected = detected
}

func (a *AprilTag) SetPoint3d(id int) {
	// 公式计算3d点坐标
	pt := a.Corners[id]
	if pt == nil {
		return
	}

}
