#!/bin/sh

# Build script with optional conditional build
# Usage:
#   ./build.sh        - Force rebuild (clean first)
#   ./build.sh fast   - Only build if dist doesn't exist

FORCE_BUILD=true

# Check if "fast" mode is enabled
if [ "$1" = "fast" ]; then
    if [ -d "./dist" ]; then
        echo "Dist directory exists, skipping build..."
        exit 0
    fi
    FORCE_BUILD=false
fi

# Clean dist directory if force build
if [ "$FORCE_BUILD" = "true" ]; then
    echo "Cleaning dist directory..."
    rm -rf ./dist
fi

# Create dist directory structure
mkdir -p dist/backend
mkdir -p dist/frontend/static

# Get version from lzc-manifest.yml
APP_VERSION=$(grep '^version:' lzc-manifest.yml | sed 's/version: *//' | tr -d '\r')

# Get git commit information
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_COMMIT_FULL=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "========================================"
echo "Building GitLab CI Runner Manager"
echo "========================================"
echo "App version: $APP_VERSION"
echo "Git commit: $GIT_COMMIT (branch: $GIT_BRANCH)"
echo "Build time: $BUILD_TIME"
echo "========================================"

# Build backend
echo ""
echo "Building backend..."
cd backend

# Build Go binary with version information
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=$APP_VERSION \
              -X main.GitCommit=$GIT_COMMIT \
              -X main.GitCommitFull=$GIT_COMMIT_FULL \
              -X main.GitBranch=$GIT_BRANCH \
              -X main.BuildTime=$BUILD_TIME" \
    -o ../dist/backend/ci-runner-server \
    main.go

if [ $? -ne 0 ]; then
    echo "Backend build failed!"
    exit 1
fi

echo "Backend build completed successfully!"

# Copy frontend files
echo ""
echo "Copying frontend files..."
cd ..
cp -r frontend/static/* dist/frontend/static/

# Generate version info file for frontend
cat > dist/frontend/static/version.js << EOF
// Build version information
window._version_ = {
  version: "$APP_VERSION",
  gitCommit: "$GIT_COMMIT",
  gitCommitFull: "$GIT_COMMIT_FULL",
  gitBranch: "$GIT_BRANCH",
  buildTime: "$BUILD_TIME"
};
EOF

echo "Frontend files copied successfully!"

# Copy static files from backend
echo ""
echo "Copying backend static files..."
cp backend/start.sh dist/start.sh
cp backend/Dockerfile dist/Dockerfile
chmod +x dist/start.sh

echo "Backend static files copied successfully!"

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Output directory: ./dist"
echo "Binary size: $(du -h dist/backend/ci-runner-server 2>/dev/null | cut -f1 || echo 'N/A')"
echo ""
echo "To build Docker image:"
echo "  cd dist && docker build -t ci-runner:$APP_VERSION ."
echo ""
echo "To run locally:"
echo "  cd dist && ./start.sh"
echo "========================================"
