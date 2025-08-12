#!/bin/bash

# Script to convert SVG logo to PNG for social media sharing
# Requires ImageMagick to be installed

echo "Converting SVG logo to PNG for social media sharing..."

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "Error: ImageMagick is not installed."
    echo "Install it with: sudo apt-get install imagemagick (Ubuntu/Debian)"
    echo "Or: brew install imagemagick (macOS)"
    exit 1
fi

# Convert SVG to PNG with proper dimensions for social media
convert static/resource/full_logo.svg \
    -resize 1200x630 \
    -background white \
    -gravity center \
    -extent 1200x630 \
    static/resource/og-image.png

if [ $? -eq 0 ]; then
    echo "✅ Successfully created og-image.png (1200x630)"
    echo "📁 File saved to: static/resource/og-image.png"
    echo "🔍 You can now test your social media previews!"
else
    echo "❌ Error converting SVG to PNG"
    echo "💡 Alternative: Create a 1200x630 PNG image manually"
    echo "   and save it as 'static/resource/og-image.png'"
fi
