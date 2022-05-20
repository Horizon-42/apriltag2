//
// Created by horizon on 2021/7/31.
//

#ifndef APRILTAG_DETECT_APRILTAG_H
#define APRILTAG_DETECT_APRILTAG_H

#include <stdbool.h>

#ifdef __cplusplus

#include <opencv2/opencv.hpp>
#include "apriltags/TagDetector.h"

class TagDetector {
    AprilTags::TagDetector *detector;
    AprilTags::TagCodes *tagCodes;
public:
    TagDetector();

    ~TagDetector();

    int CountTags(cv::Mat const &frame);

    bool DetectTags(cv::Mat const &frame, cv::Mat &points, cv::Mat &ids, bool draw);
};

typedef cv::Mat *Mat;
using namespace std;

typedef TagDetector *TagDetectorPtr;

#else
typedef void *Mat;
typedef void *TagDetectorPtr;
#endif

#ifdef __cplusplus


extern "C"
{
#include "apriltags/Tag16h5.h"
#include "apriltags/Tag25h7.h"
#include "apriltags/Tag25h9.h"
#include "apriltags/Tag36h9.h"
#include "apriltags/Tag36h11.h"
#endif

TagDetectorPtr NewTagDetector();
void ReleaseTagDetector(TagDetectorPtr *detector);
int CountTags(TagDetectorPtr detector, Mat frame);
bool DetectTags(TagDetectorPtr detector, Mat frame, Mat points, Mat ids, bool draw);
bool IsEmpty(TagDetectorPtr detector);

bool Init();

int HaveAprilTags(Mat frame);

void Close();

#ifdef __cplusplus
};
#endif

#endif //APRILTAG_DETECT_APRILTAG_H
