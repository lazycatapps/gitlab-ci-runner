FROM golang:1.25.4-bookworm

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Download the GitLab Runner binary
RUN curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64 && \
    chmod +x /usr/local/bin/gitlab-runner

# Create a GitLab Runner user
RUN useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash

# Install and configure as a service
RUN gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner

# Switch to gitlab-runner user
USER gitlab-runner
WORKDIR /home/gitlab-runner
