#include "converter.hpp"

std::string detectFormat(const std::vector<uchar>& buf) {
    if (buf.size() >= 4) {
        // JPEG: FF D8
        if (buf[0] == 0xFF && buf[1] == 0xD8) return ".jpg";
        // PNG: 89 50 4E 47
        if (buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47) return ".png";
        // WebP: RIFF ... WEBP
        if (buf.size() >= 12 && 
            buf[0] == 0x52 && buf[1] == 0x49 && buf[2] == 0x46 && buf[3] == 0x46 &&
            buf[8] == 0x57 && buf[9] == 0x45 && buf[10] == 0x42 && buf[11] == 0x50) return ".webp";
    }
    return ".jpg"; 
}

std::vector<int> getEncodingParams(const std::string& formatExt) {
    if (formatExt == ".jpg" || formatExt == ".jpeg") {
        return { cv::IMWRITE_JPEG_QUALITY, 85 };
    } else if (formatExt == ".png") {
        return { cv::IMWRITE_PNG_COMPRESSION, 4 };
    } else if (formatExt == ".webp") {
        return { cv::IMWRITE_WEBP_QUALITY, 85 };
    }
    return {};
}

std::string standardizeFormat(const std::string& format) {
    std::string ext = format;
    if (ext.empty()) return ".jpg";
    if (ext[0] != '.') ext = "." + ext; // Ensure it starts with a dot
    if (ext == ".jpeg") return ".jpg";
    return ext;
}