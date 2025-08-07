#!/bin/bash

# Test script for video thumbnail generation
echo "Testing video thumbnail generation..."

# Check if ffmpeg is available in Docker
echo "Testing ffmpeg in Docker container..."
docker exec -it olhourbano2-backend-1 ffmpeg -version | head -1

# Test thumbnail generation with a sample video (if exists)
if [ -f "test_video.mp4" ]; then
    echo "Generating thumbnail for test_video.mp4..."
    docker exec -it olhourbano2-backend-1 ffmpeg \
        -i /olhourbano2/test_video.mp4 \
        -ss 00:00:01 \
        -vframes 1 \
        -vf "scale=150x150:force_original_aspect_ratio=decrease,pad=150x150:(ow-iw)/2:(oh-ih)/2" \
        -y /olhourbano2/uploads/thumbnails/test_video_thumb.jpg 2>/dev/null
    
    if [ -f "uploads/thumbnails/test_video_thumb.jpg" ]; then
        echo "✅ Video thumbnail generated successfully"
        ls -la uploads/thumbnails/test_video_thumb.jpg
    else
        echo "❌ Video thumbnail generation failed"
    fi
else
    echo "⚠️  No test video found. To test video thumbnails:"
    echo "   1. Place a test video file named 'test_video.mp4' in the project root"
    echo "   2. Run this script again"
fi

echo "Video thumbnail test completed!"
