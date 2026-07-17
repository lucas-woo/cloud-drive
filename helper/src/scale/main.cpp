#include <iostream>
#include <vector>
#include <string>
#include <opencv2/opencv.hpp>

#ifdef _WIN32
#include <fcntl.h>
#include <io.h>
#endif

int main(int argc, char* argv[]) {
    // 1. Validate command-line arguments
    if (argc != 3) {
        std::cerr << "Usage: " << argv[0] << " <width> <height>" << std::endl;
        return 1; // Non-zero exit code will be caught by Go's cmd.Wait()
    }

    int target_width = std::stoi(argv[1]);
    int target_height = std::stoi(argv[2]);

    if (target_width <= 0 || target_height <= 0) {
        std::cerr << "Error: Width and height must be positive integers." << std::endl;
        return 1;
    }

    // 2. Force stdin and stdout to binary mode 
    // (Crucial if this ever runs on Windows to prevent \n translation corruption)
#ifdef _WIN32
    _setmode(_fileno(stdin), _O_BINARY);
    _setmode(_fileno(stdout), _O_BINARY);
#endif

    // 3. Read the incoming gRPC binary stream from stdin into a buffer
    std::vector<uchar> in_buffer;
    char chunk[4096];
    
    while (std::cin.read(chunk, sizeof(chunk))) {
        in_buffer.insert(in_buffer.end(), chunk, chunk + std::cin.gcount());
    }
    // Catch the final partial chunk
    if (std::cin.gcount() > 0) {
        in_buffer.insert(in_buffer.end(), chunk, chunk + std::cin.gcount());
    }

    if (in_buffer.empty()) {
        std::cerr << "Error: No image data received on stdin." << std::endl;
        return 1;
    }

    // 4. Decode the image (OpenCV automatically detects JPEG, PNG, etc.)
    cv::Mat img = cv::imdecode(in_buffer, cv::IMREAD_COLOR);
    if (img.empty()) {
        std::cerr << "Error: Failed to decode image. Data may be corrupted." << std::endl;
        return 1;
    }

    // 5. Resize the image
    // INTER_AREA is the best interpolation method for shrinking images down
    cv::Mat resized_img;
    cv::resize(img, resized_img, cv::Size(target_width, target_height), 0, 0, cv::INTER_AREA);

    // 6. Encode the resized image back to a byte buffer
    std::vector<uchar> out_buffer;
    std::vector<int> encode_params = {cv::IMWRITE_JPEG_QUALITY, 85}; 
    
    // Note: We default to outputting as JPEG. 
    if (!cv::imencode(".jpg", resized_img, out_buffer, encode_params)) {
        std::cerr << "Error: Failed to encode resized image." << std::endl;
        return 1;
    }

    // 7. Stream the finished bytes to stdout (which Go pipes to S3)
    std::cout.write(reinterpret_cast<const char*>(out_buffer.data()), out_buffer.size());
    std::cout.flush();

    return 0; // Success!
}