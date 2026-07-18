#ifndef CONVERTER_HPP
#define CONVERTER_HPP

#include <opencv2/opencv.hpp>
#include <vector>
#include <string>

// Detects the original format of the image from raw bytes
std::string detectFormat(const std::vector<uchar>& buffer);

// Generates the appropriate encoding parameters for OpenCV
std::vector<int> getEncodingParams(const std::string& formatExt);

// Validates and standardizes format extensions (e.g., "webp" -> ".webp")
std::string standardizeFormat(const std::string& format);

#endif