#include <iostream>
#include <vector>
#include <string>

#include <opencv2/opencv.hpp>
#include "scale.hpp"

// Helper function to detect image format from magic bytes
std::string detectFormat(const std::vector<uchar>& buf) {
    if (buf.size() >= 4) {
        // JPEG magic bytes: FF D8
        if (buf[0] == 0xFF && buf[1] == 0xD8) {
            return ".jpg";
        }
        // PNG magic bytes: 89 50 4E 47
        if (buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47) {
            return ".png";
        }
        // WebP magic bytes: RIFF (0-3) ... WEBP (8-11)
        if (buf.size() >= 12 && 
            buf[0] == 0x52 && buf[1] == 0x49 && buf[2] == 0x46 && buf[3] == 0x46 &&
            buf[8] == 0x57 && buf[9] == 0x45 && buf[10] == 0x42 && buf[11] == 0x50) {
            return ".webp";
        }
    }
    // Fallback if unknown
    return ".jpg"; 
}

int main(int argc, char* argv[]) {

    /*
        Usage:
        processor scale <width> <height>
    */

    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <operation> [args]" << std::endl;
        return 1;
    }

    std::string operation = argv[1];

    // Read image bytes from stdin
    std::vector<uchar> buffer;
    char chunk[4096];

    while (std::cin.read(chunk, sizeof(chunk))) {
        buffer.insert(buffer.end(), chunk, chunk + std::cin.gcount());
    }

    if (std::cin.gcount() > 0) {
        buffer.insert(buffer.end(), chunk, chunk + std::cin.gcount());
    }

    if (buffer.empty()) {
        std::cerr << "No image data received" << std::endl;
        return 1;
    }

    // Detect format BEFORE decoding
    std::string formatExt = detectFormat(buffer);

    // Decode image
    cv::Mat image = cv::imdecode(buffer, cv::IMREAD_COLOR);

    if (image.empty()) {
        std::cerr << "Failed to decode image" << std::endl;
        return 1;
    }

    cv::Mat output;

    if (operation == "scale") {
        if (argc != 4) {
            std::cerr << "Usage: processor scale <width> <height>" << std::endl;
            return 1;
        }

        int width = std::stoi(argv[2]);
        int height = std::stoi(argv[3]);

        if (!scaleImage(image, output, width, height)) {
            std::cerr << "Scaling failed" << std::endl;
            return 1;
        }
    } else {
        std::cerr << "Unknown operation: " << operation << std::endl;
        return 1;
    }

    // Set encoding parameters based on the detected format
    std::vector<uchar> encoded;
    std::vector<int> params;

    if (formatExt == ".jpg") {
        params = { cv::IMWRITE_JPEG_QUALITY, 85 };
    } else if (formatExt == ".png") {
        params = { cv::IMWRITE_PNG_COMPRESSION, 4 };
    } else if (formatExt == ".webp") {
        params = { cv::IMWRITE_WEBP_QUALITY, 85 };
    }

    if (!cv::imencode(formatExt, output, encoded, params)) {
        std::cerr << "Failed encoding image as " << formatExt << std::endl;
        return 1;
    }


    std::cout.write(reinterpret_cast<const char*>(encoded.data()), encoded.size());
    std::cout.flush();

    return 0;
}