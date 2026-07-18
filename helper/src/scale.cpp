#include "scale.hpp"

bool scaleImage(
    const cv::Mat& input,
    cv::Mat& output,
    int width,
    int height
) {
    if (input.empty() || width <= 0 || height <= 0) {
        return false;
    }

    cv::resize(
        input,
        output,
        cv::Size(width, height),
        0,
        0,
        cv::INTER_AREA
    );

    return !output.empty();
}