FROM ubuntu:22.04

# 安装依赖包
RUN apt-get update && \
    apt-get install -y curl ca-certificates git && \
    rm -rf /var/lib/apt/lists/*

# 下载 GitLab Runner 二进制文件
RUN curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64 && \
    chmod +x /usr/local/bin/gitlab-runner

# 创建 GitLab Runner 用户
RUN useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash && \
    rm -rf /home/gitlab-runner/.bash_logout

# 安装并配置为服务
RUN gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner

# 创建必要的目录
RUN mkdir -p /home/gitlab-runner/.gitlab-runner && \
    chown -R gitlab-runner:gitlab-runner /home/gitlab-runner

# 设置工作目录
WORKDIR /app

# 复制应用文件
COPY --chmod=755 start.sh /app/start.sh
COPY --chmod=755 backend /app/backend
COPY frontend /app/frontend

# 以 root 用户运行,允许管理 gitlab-runner 进程
USER root

# 设置入口点
ENTRYPOINT ["/app/start.sh"]
