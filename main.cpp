//
// Created by liudx on 2022/10/10.
//
#include <opencv2/opencv.hpp>
#include "detect_apriltag.h"
#include <filesystem>

namespace fs = std::filesystem;

int main()
{
    std::string data_dir = "C:/Users/liudx/Desktop/record_C08";
    cv::namedWindow("show", cv::WINDOW_GUI_NORMAL);
    auto tag_detector = NewTagDetector();
    for (auto &&v : fs::directory_iterator(data_dir))
    {
        if (fs::is_regular_file(v) && v.path().extension().string() == ".png")
        {
            std::cout << v.path() << "\n";
            auto frame = cv::imread(v.path().string());
            cv::imshow("show", frame);
            auto is_blur = tag_detector->is_blur(frame);
            char key = cv::waitKey(5);
            if (is_blur)
                cv::waitKey(0);
            if (key == 27)
                break;
        }
    }
    return 0;
}