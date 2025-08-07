#!/bin/bash

# Test script for thumbnail generation
echo "Testing thumbnail generation system..."

# Check if required tools are installed
echo "Checking dependencies..."

if command -v ffmpeg &> /dev/null; then
    echo "✅ ffmpeg is installed"
else
    echo "❌ ffmpeg is not installed"
    exit 1
fi

if command -v convert &> /dev/null; then
    echo "✅ ImageMagick (convert) is installed"
else
    echo "❌ ImageMagick (convert) is not installed"
    exit 1
fi

# Create test directories
echo "Creating test directories..."
mkdir -p uploads/thumbnails

# Test video thumbnail generation
echo "Testing video thumbnail generation..."
if [ -f "test_video.mp4" ]; then
    ffmpeg -i test_video.mp4 -ss 00:00:01 -vframes 1 -vf "scale=150x150:force_original_aspect_ratio=decrease,pad=150x150:(ow-iw)/2:(oh-ih)/2" -y uploads/thumbnails/test_video_thumb.jpg 2>/dev/null
    if [ -f "uploads/thumbnails/test_video_thumb.jpg" ]; then
        echo "✅ Video thumbnail generated successfully"
    else
        echo "❌ Video thumbnail generation failed"
    fi
else
    echo "⚠️  No test video found, skipping video test"
fi

# Test PDF thumbnail generation
echo "Testing PDF thumbnail generation..."
if [ -f "test_document.pdf" ]; then
    convert -density 150 -resize 150x150 -background white -alpha remove -alpha off "test_document.pdf[0]" uploads/thumbnails/test_document_thumb.jpg 2>/dev/null
    if [ -f "uploads/thumbnails/test_document_thumb.jpg" ]; then
        echo "✅ PDF thumbnail generated successfully"
    else
        echo "❌ PDF thumbnail generation failed"
    fi
else
    echo "⚠️  No test PDF found, skipping PDF test"
fi

echo "Thumbnail generation test completed!"
