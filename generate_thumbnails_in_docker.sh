#!/bin/bash

# Script to generate thumbnails for existing image files inside Docker container
echo "Generating thumbnails for existing image files in Docker container..."

# Run the thumbnail generation inside the Docker container
docker exec -w /olhourbano2 olhourbano2-backend-1 bash -c '
echo "Checking for ImageMagick in container..."
if command -v convert &> /dev/null; then
    echo "✅ ImageMagick is available in container"
else
    echo "❌ ImageMagick is not available in container"
    exit 1
fi

echo "Creating thumbnails directory if it does not exist..."
mkdir -p uploads/thumbnails

echo "Finding all image files in uploads directory..."
find uploads/ -maxdepth 1 -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.gif" -o -iname "*.webp" -o -iname "*.bmp" \) | while read -r image_file; do
    # Skip if it is already a thumbnail
    if [[ "$image_file" == *"_thumb."* ]]; then
        continue
    fi
    
    # Extract filename and hash
    filename=$(basename "$image_file")
    extension="${filename##*.}"
    hash="${filename%.*}"
    thumbnail_path="uploads/thumbnails/${hash}_thumb.jpg"
    
    # Check if thumbnail already exists
    if [ -f "$thumbnail_path" ]; then
        echo "✅ Thumbnail already exists for: $filename"
        continue
    fi
    
    # Generate thumbnail
    echo "Generating thumbnail for: $filename"
    convert "$image_file" \
        -resize 150x150 \
        -background white \
        -gravity center \
        -extent 150x150 \
        "$thumbnail_path"
    
    if [ $? -eq 0 ]; then
        echo "✅ Generated thumbnail: $thumbnail_path"
    else
        echo "❌ Failed to generate thumbnail for: $filename"
    fi
done

echo "Thumbnail generation completed in Docker container!"
'
