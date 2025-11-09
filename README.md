# GitLab CI Runner 管理器

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Lazycat Platform](https://img.shields.io/badge/Platform-Lazycat-green.svg)](https://lazycat.cloud)

## 项目简介

本项目是一个为 [GitLab Runner](https://docs.gitlab.com/runner/) 设计的 Web 管理界面，让用户可以在 Lazycat 平台上轻松部署和管理多个 GitLab CI/CD Runner。

GitLab Runner 是 GitLab CI/CD 的核心组件,负责执行持续集成和持续交付任务。通过本项目提供的 Web 界面，您可以方便地注册、启动、停止 Runner，查看运行状态和日志，无需复杂的命令行操作。

## 主要功能

- 🚀 **简单易用的 Web 界面** - 通过浏览器即可管理所有 Runner
- 📝 **快速注册 Runner** - 支持通过界面输入 URL 和 Token 快速注册新 Runner
- 📊 **实时状态监控** - 查看每个 Runner 的运行状态（运行中/已停止）
- 📋 **日志查看** - 实时查看 Runner 的运行日志（最近 1000 行）
- 🔄 **单独重启** - 支持单独重启指定的 Runner
- 🗑️ **删除管理** - 可以删除不需要的 Runner
- 💾 **数据持久化** - Runner 配置自动保存，重启不丢失
- 🔒 **安全隔离** - 基于 Lazycat 平台的容器隔离技术

## 技术特性

### 后端实现

- **语言**: Go 语言开发,性能优异
- **进程管理**: 使用 `nohup` 在后台运行 Runner
- **状态追踪**: 通过 PID 文件跟踪每个 Runner 的进程状态
- **日志管理**: 每个 Runner 的日志独立保存到文件
- **API 接口**: 提供 RESTful API 供前端调用

### 前端实现

- **技术栈**: 原生 HTML/CSS/JavaScript
- **界面风格**: 简洁现代的 UI 设计
- **响应式布局**: 支持桌面和移动端访问

### 目录结构

```
/home/gitlab-runner/.gitlab-runner/
├── config.toml              # GitLab Runner 主配置文件
├── pids/                    # PID 文件目录
│   └── <runner-name>.pid   # 每个 Runner 的进程 ID
└── logs/                    # 日志文件目录
    └── <runner-name>.log   # 每个 Runner 的运行日志
```

## API 端点

- `GET /api/version` - 获取版本信息
- `POST /api/runners/register` - 注册新 Runner
- `GET /api/runners` - 获取所有 Runners 列表及状态
- `POST /api/runners/delete` - 删除指定 Runner
- `POST /api/runners/restart` - 重启指定 Runner
- `GET /api/runners/logs?name=<runner-name>` - 获取 Runner 日志

## 使用限制

> **注意**: 当前版本仅支持 **shell 模式** 的 GitLab Runner。
>
> Docker executor、Kubernetes executor 等其他执行器模式暂不支持，敬请期待后续版本更新！

## 致谢

本项目基于开源社区的杰出贡献：

- **GitLab 团队**: 感谢 [GitLab](https://about.gitlab.com/) 和 [GitLab Runner](https://docs.gitlab.com/runner/) 项目团队开发和维护这个优秀的 CI/CD 平台
- **Go 语言社区**: 感谢 Go 语言及其丰富的生态系统
- **开源社区**: 感谢所有为 GitLab 相关项目做出贡献的开发者
- **Lazycat 平台**: 提供便捷的应用部署和容器管理能力

## 版权说明

- 本仓库的代码和配置文件采用 [Apache License 2.0](LICENSE)
- GitLab Runner 软件本身采用 [MIT License](https://gitlab.com/gitlab-org/gitlab-runner/-/blob/main/LICENSE)

## 相关链接

- 项目仓库: https://github.com/lazycatapps/gitlab-ci-runner
- GitLab Runner 官方文档: https://docs.gitlab.com/runner/
- GitLab CI/CD 文档: https://docs.gitlab.com/ee/ci/
- Lazycat 平台: https://lazycat.cloud

## 开发者信息

- 作者: xiao
- 维护: LazyCat Apps 团队
- 问题反馈: [GitHub Issues](https://github.com/lazycatapps/gitlab-ci-runner/issues)

---

Made with ❤️ for the Lazycat Platform
