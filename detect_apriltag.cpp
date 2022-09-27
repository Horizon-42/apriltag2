//
// Created by horizon on 2021/7/31.
//

#include "detect_apriltag.h"

//#define DEBUG

AprilTags::TagDetector *m_tagDetector(NULL);
AprilTags::TagCodes *m_tagCodes(NULL);

bool Init() {
    m_tagCodes = new AprilTags::TagCodes(AprilTags::tagCodes36h11);
    m_tagDetector = new AprilTags::TagDetector(*m_tagCodes);
    return true;
}

int HaveAprilTags(Mat frame) {
    if ((*frame).empty())
        return 0;
    cv::Mat image_gray;
    if ((*frame).channels() == 3)
        cv::cvtColor(*frame, image_gray, cv::COLOR_BGR2GRAY);
    else
        image_gray = (*frame).clone();

    vector<AprilTags::TagDetection> detections = m_tagDetector->extractTags(image_gray);
    image_gray.release();
    auto count = detections.size();
    vector<AprilTags::TagDetection>().swap(detections);
    return count;
}

void Close() {
    delete m_tagDetector, m_tagCodes;
}

TagDetector::TagDetector() {
    tagCodes = new AprilTags::TagCodes(AprilTags::tagCodes36h11);
    detector = new AprilTags::TagDetector(*tagCodes);
}

TagDetector::~TagDetector() {
    if (detector != NULL) {
        delete detector;
        detector = NULL;
    }
    if (tagCodes != NULL) {
        delete tagCodes;
        tagCodes = NULL;
    }
}

int TagDetector::CountTags(const cv::Mat &frame) {
    if (frame.empty())
        return 0;
    cv::Mat image_gray;
    if (frame.channels() == 3)
        cv::cvtColor(frame, image_gray, cv::COLOR_BGR2GRAY);
    else
        image_gray = frame.clone();

    vector<AprilTags::TagDetection> detections = detector->extractTags(image_gray);
    image_gray.release();
    auto count = detections.size();
    vector<AprilTags::TagDetection>().swap(detections);
    return count;
}

bool TagDetector::DetectTags(const cv::Mat &frame, cv::Mat &points, cv::Mat &ids, bool draw) {
    if (frame.empty())
        return false;
    cv::Mat image_gray;
    if (frame.channels() == 3)
        cv::cvtColor(frame, image_gray, cv::COLOR_BGR2GRAY);
    else
        image_gray = frame.clone();
    vector<AprilTags::TagDetection> detections = detector->extractTags(image_gray);
//    std::cout<<detections.size()<<"\n";

    bool ret = false;
    if (detections.size() > 0) {
        points = cv::Mat(detections.size() * 4, 2, CV_32F);
        ids = cv::Mat(detections.size(), 1, CV_32S);
        for (int i = 0; i < detections.size(); ++i) {
            for (int j = 0; j < 4; ++j) {
                points.at<float>(i * 4 + j, 0) = detections[i].p[j].first;
                points.at<float>(i * 4 + j, 1) = detections[i].p[j].second;
            }
            ids.at<int>(i, 0) = detections[i].id;
        }

//        // 计算亚像素角点
//        try {
//            cv::cornerSubPix(image_gray, points, {5, 5}, {-1, -1},
//                             cv::TermCriteria(cv::TermCriteria::EPS + cv::TermCriteria::MAX_ITER,
//                                              40, 0.001));
//        } catch (std::exception const &e) {
//            std::cout << "sub pix corners calculate failed, " << e.what() << ".\n";
//        }

        if (draw) {
            for (int i = 0; i < detections.size(); ++i) {
                cv::circle(frame, {int(detections[i].cxy.first), int(detections[i].cxy.second)}, 3,
                           {0, 255, 0}, -1);
                cv::circle(frame, {int(detections[i].p[0].first), int(detections[i].p[0].second)}, 3,
                           {255, 0, 0}, -1);
                cv::circle(frame, {int(detections[i].p[1].first), int(detections[i].p[1].second)}, 3,
                           {0, 255, 0}, -1);
                cv::circle(frame, {int(detections[i].p[2].first), int(detections[i].p[2].second)}, 3,
                           {0, 0, 255}, -1);

                cv::putText(frame, to_string(detections[i].id),
                            {int(detections[i].p[0].first), int(detections[i].p[0].second)}, 0,
                            .5, {0, 0, 255});
            }
        }
        ret = true;
    }
    image_gray.release();
    vector<AprilTags::TagDetection>().swap(detections);
    return ret;
}

TagDetectorPtr NewTagDetector() {
    return new TagDetector;
}

void ReleaseTagDetector(TagDetectorPtr *detector) {
    if (*detector != NULL) {
        delete (*detector);
        *detector = NULL;
    }
}

int CountTags(TagDetectorPtr detector, Mat frame) {
    return detector->CountTags(*frame);
}

bool DetectTags(TagDetectorPtr detector, Mat frame, Mat points, Mat ids, bool draw) {
    return detector->DetectTags(*frame, *points, *ids, draw);
}

bool IsEmpty(TagDetectorPtr detector) {
    return detector == NULL;
}