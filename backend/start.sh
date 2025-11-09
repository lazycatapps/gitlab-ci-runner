#!/bin/sh
set -e

# Ensure config directory exists
mkdir -p /home/gitlab-runner/.gitlab-runner

# Start the backend server
cd /app/backend
exec ./ci-runner-server
