#ifndef SCALE_HPP
#define SCALE_HPP

#include <opencv2/opencv.hpp>

bool scaleImage(
    const cv::Mat& input,
    cv::Mat& output,
    int width,
    int height
);

#endif