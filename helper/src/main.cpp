#include <iostream>
#include <vector>
#include <string>

#include <opencv2/opencv.hpp>

#include "scale.hpp"


int main(int argc, char* argv[]) {

    /*
        Usage:

        processor scale <width> <height>

    */

    if (argc < 2) {
        std::cerr << "Usage: "
                  << argv[0]
                  << " <operation> [args]"
                  << std::endl;

        return 1;
    }


    std::string operation = argv[1];


    // Read image bytes from stdin
    std::vector<uchar> buffer;

    char chunk[4096];

    while (std::cin.read(chunk, sizeof(chunk))) {
        buffer.insert(
            buffer.end(),
            chunk,
            chunk + std::cin.gcount()
        );
    }

    if (std::cin.gcount() > 0) {
        buffer.insert(
            buffer.end(),
            chunk,
            chunk + std::cin.gcount()
        );
    }


    if (buffer.empty()) {
        std::cerr << "No image data received" << std::endl;
        return 1;
    }


    // Decode image
    cv::Mat image = cv::imdecode(
        buffer,
        cv::IMREAD_COLOR
    );


    if (image.empty()) {
        std::cerr << "Failed to decode image" << std::endl;
        return 1;
    }


    cv::Mat output;


    if (operation == "scale") {

        if (argc != 4) {
            std::cerr
                << "Usage: processor scale <width> <height>"
                << std::endl;

            return 1;
        }


        int width = std::stoi(argv[2]);
        int height = std::stoi(argv[3]);


        if (!scaleImage(image, output, width, height)) {
            std::cerr << "Scaling failed" << std::endl;
            return 1;
        }

    } else {

        std::cerr
            << "Unknown operation: "
            << operation
            << std::endl;

        return 1;
    }


    // Encode output JPEG
    std::vector<uchar> encoded;

    std::vector<int> params = {
        cv::IMWRITE_JPEG_QUALITY,
        85
    };


    if (!cv::imencode(
            ".jpg",
            output,
            encoded,
            params
        )) {

        std::cerr << "Failed encoding image"
                  << std::endl;

        return 1;
    }


    // Send bytes back to Go/S3
    std::cout.write(
        reinterpret_cast<const char*>(encoded.data()),
        encoded.size()
    );

    std::cout.flush();


    return 0;
}