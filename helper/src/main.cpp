#include <iostream>
#include <vector>
#include <string>
#include <opencv2/opencv.hpp>

#include "scale.hpp"
#include "converter.hpp"

int main(int argc, char* argv[]) {

    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <operation> [args]..." << std::endl;
        return 1;
    }

    // 1. Read image bytes from stdin
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

    // 2. Initial Setup
    std::string targetFormat = detectFormat(buffer); // Default to input format
    cv::Mat image = cv::imdecode(buffer, cv::IMREAD_COLOR);
    
    if (image.empty()) {
        std::cerr << "Failed to decode image" << std::endl;
        return 1;
    }

    // 3. Process Arguments Pipeline
    int i = 1;
    while (i < argc) {
        std::string op = argv[i];

        if (op == "scale") {
            if (i + 2 >= argc) {
                std::cerr << "Error: 'scale' requires width and height." << std::endl;
                return 1;
            }
            int width = std::stoi(argv[++i]);
            int height = std::stoi(argv[++i]);
            
            cv::Mat output;
            if (!scaleImage(image, output, width, height)) {
                std::cerr << "Scaling failed" << std::endl;
                return 1;
            }
            image = output; // Overwrite current image state for the next step
        } 
        else if (op == "convert") {
            if (i + 1 >= argc) {
                std::cerr << "Error: 'convert' requires a format (e.g., webp, png, jpg)." << std::endl;
                return 1;
            }
            targetFormat = standardizeFormat(argv[++i]);
        }
        else {
            std::cerr << "Unknown operation: " << op << std::endl;
            return 1;
        }
        i++;
    }

    // 4. Encode and Output
    std::vector<uchar> encoded;
    std::vector<int> params = getEncodingParams(targetFormat);

    if (!cv::imencode(targetFormat, image, encoded, params)) {
        std::cerr << "Failed encoding image as " << targetFormat << std::endl;
        return 1;
    }

    std::cout.write(reinterpret_cast<const char*>(encoded.data()), encoded.size());
    std::cout.flush();

    return 0;
}