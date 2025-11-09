# ci-runner

安装 gitlab runner cli

```bash
# Download the binary for your system
sudo curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64

# Give it permission to execute
sudo chmod +x /usr/local/bin/gitlab-runner

# Create a GitLab Runner user
sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash

# Install and run as a service
sudo gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
sudo gitlab-runner start
```

注册应用

```bash
gitlab-runner register  --url http://gitlab.minicat.heiyu.space  --token glrt-NGgBudPC5wqHKcWYuixy
gitlab-runner start
```

## 构建项目

### 方式一：使用构建脚本（推荐）

```bash
# 完整构建（清理后重新构建）
./build.sh

# 快速构建（如果 dist 目录存在则跳过）
./build.sh fast
```

构建完成后，会在 `dist/` 目录生成：

- `backend/ci-runner-server` - 编译好的 Go 二进制文件
- `frontend/static/` - 前端静态文件
- `start.sh` - 启动脚本
- `Dockerfile` - 用于构建 Docker 镜像
- `README.md` - 发布说明

### 方式二：Docker 构建

```bash
cd dist
docker build -t ci-runner:0.0.1 .
docker run -d -p 8098:8098 --name ci-runner ci-runner:0.0.1
```

### 方式三：手动构建后端

```bash
cd backend
go mod download
go build -o ci-runner-server main.go
./ci-runner-server
```

访问 `http://localhost:8098` 即可使用 Web 界面。

## 开发

### 运行后端服务

```bash
cd backend
go run main.go
```

### 修改前端

前端使用原生 HTML/CSS/JavaScript，直接编辑 `frontend/static/` 下的文件即可。

## 功能特性

- **注册 Runner**: 通过 Web 界面注册新的 GitLab Runner
- **查看状态**: 实时查看每个 Runner 的运行状态（运行中/已停止）
- **查看日志**: 查看单个 Runner 的运行日志（最近 1000 行）
- **单独重启**: 重启指定的 Runner（kill 进程后重新启动）
- **后台运行**: 使用 nohup 在后台运行 Runner，日志保存到文件
- **进程管理**: 自动管理 PID 文件，跟踪每个 Runner 的进程状态
- **删除 Runner**: 删除不需要的 Runner

## 技术实现

### Runner 进程管理

- 每个 Runner 使用 `nohup gitlab-runner run` 在后台运行
- PID 文件保存在 `/home/gitlab-runner/.gitlab-runner/pids/<runner-name>.pid`
- 日志文件保存在 `/home/gitlab-runner/.gitlab-runner/logs/<runner-name>.log`
- 通过检查 PID 文件和进程状态来判断 Runner 是否运行

### 重启机制

1. 查找 Runner 的 PID 文件
2. 使用 `kill` 命令停止进程
3. 使用 `nohup gitlab-runner run --service <name>` 重新启动
4. 保存新的 PID 到文件

## API 端点

- `GET /api/version` - 获取版本信息
- `POST /api/runners/register` - 注册新 runner
- `GET /api/runners` - 获取所有 runners（包含状态信息）
- `POST /api/runners/delete` - 删除指定 runner（需要 token）
- `POST /api/runners/restart` - 重启指定 runner（需要 name，会 kill 进程并重新启动）
- `GET /api/runners/logs?name=<runner-name>` - 获取指定 runner 的日志（最近 1000 行）
