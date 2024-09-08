# 1aides

[English](README.md) | 中文

[![Release](https://img.shields.io/github/v/release/tongque0/1aides)](https://github.com/tongque0/1aides/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/tongque0/1aides)](https://goreportcard.com/report/github.com/tongque0/1aides)
[![OpenIssue](https://img.shields.io/github/issues/tongque0/1aides)](https://github.com/tongque0/1aides/issues)
[![ClosedIssue](https://img.shields.io/github/issues-closed/tongque0/1aides)](https://github.com/tongque0/1aides/issues?q=is%3Aissue+is%3Aclosed)
![Stars](https://img.shields.io/github/stars/tongque0/1aides)


## 快速开始

确保您的系统中已安装 Docker

`用户名: admin 密码: Aides123.`
### 使用 Bash 脚本启动

```bash
curl -sL https://raw.githubusercontent.com/tongque0/1aides/main/start.sh -o start.sh && chmod +x start.sh && ./start.sh
```

### 使用 Docker Compose 启动

要通过 Docker Compose 启动服务，请先下载 `/deploy目录`下`docker-compose.yaml` 文件，然后使用 Docker Compose 来启动服务。执行以下步骤：

   ```bash
   docker-compose up -d
   ```


## 特点

- 自带后台管理页面
- 自带命令行管理
- 支持定时任务
- 支持多模型
- 支持记忆，多轮对话


## 开源许可

1aides 基于 [Apache License 2.0](https://github.com/tongque0/1aides/LICENSE) 许可证，其依赖的三方组件的开源许可见 [Here](https://github.com/eatmoreapple/openwechat/blob/master/LICENSE)。

