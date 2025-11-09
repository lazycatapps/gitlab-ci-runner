FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl ca-certificates git && \
    rm -rf /var/lib/apt/lists/*

# Download the GitLab Runner binary
RUN curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64 && \
    chmod +x /usr/local/bin/gitlab-runner

# Create a GitLab Runner user
RUN useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash
RUN rm -rf /home/gitlab-runner/.bash_logout

# Install and configure as a service
RUN gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner

# Switch to gitlab-runner user
USER 0
WORKDIR /app

COPY --chmod=755 start.sh /app/start.sh
COPY --chmod=755 backend /app/backend
COPY frontend /app/frontend

# Set entrypoint
ENTRYPOINT ["/app/start.sh"]
