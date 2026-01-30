#!/bin/bash

# ArLOG Desktop - Icon Generation Script
# This script creates all required icon formats from a master PNG

echo "🎨 Creating ArLOG Desktop Icons..."

# Check if source icon exists
if [ ! -f "resources/icon.png" ]; then
    echo "❌ Error: resources/icon.png not found!"
    echo "Please create a 1024x1024 PNG icon first."
    echo "See ICON_SETUP.md for instructions."
    exit 1
fi

# Create build directory
mkdir -p build/icons

echo "📦 Creating PNG icons..."

# Generate PNG icons for all sizes
sips -z 16 16     resources/icon.png --out build/icons/16x16.png
sips -z 32 32     resources/icon.png --out build/icons/32x32.png
sips -z 64 64     resources/icon.png --out build/icons/64x64.png
sips -z 128 128   resources/icon.png --out build/icons/128x128.png
sips -z 256 256   resources/icon.png --out build/icons/256x256.png
sips -z 512 512   resources/icon.png --out build/icons/512x512.png
sips -z 1024 1024 resources/icon.png --out build/icons/1024x1024.png

# Copy 512x512 as main Linux icon
cp build/icons/512x512.png build/icon.png

echo "🍎 Creating macOS .icns icon..."

# Create iconset for macOS
mkdir -p icon.iconset
cp build/icons/16x16.png    icon.iconset/icon_16x16.png
cp build/icons/32x32.png    icon.iconset/icon_16x16@2x.png
cp build/icons/32x32.png    icon.iconset/icon_32x32.png
cp build/icons/64x64.png    icon.iconset/icon_32x32@2x.png
cp build/icons/128x128.png  icon.iconset/icon_128x128.png
cp build/icons/256x256.png  icon.iconset/icon_128x128@2x.png
cp build/icons/256x256.png  icon.iconset/icon_256x256.png
cp build/icons/512x512.png  icon.iconset/icon_256x256@2x.png
cp build/icons/512x512.png  icon.iconset/icon_512x512.png
cp build/icons/1024x1024.png icon.iconset/icon_512x512@2x.png

# Convert to .icns
iconutil -c icns icon.iconset -o build/icon.icns

# Cleanup
rm -rf icon.iconset

echo "🪟 Creating Windows .ico icon..."

# For Windows icon, you'll need to convert online or use ImageMagick
if command -v convert &> /dev/null; then
    convert resources/icon.png -define icon:auto-resize=256,128,96,64,48,32,16 build/icon.ico
    echo "✅ Windows .ico created"
else
    echo "⚠️  ImageMagick not found. Please convert to .ico online:"
    echo "   Visit: https://icoconvert.com/"
    echo "   Upload: resources/icon.png"
    echo "   Download and save as: build/icon.ico"
fi

echo ""
echo "✅ Icons created successfully!"
echo ""
echo "📁 Icon files:"
echo "   macOS:   build/icon.icns"
echo "   Windows: build/icon.ico (may need manual conversion)"
echo "   Linux:   build/icon.png"
echo ""
echo "🚀 Ready to build app with icons!"
echo "   Run: npm run build:mac"
echo ""



