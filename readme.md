# 1aides

[English](README.md) | 中文

[![Release](https://img.shields.io/github/v/release/tongque0/1aides)](https://github.com/tongque0/1aides/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/tongque0/1aides)](https://goreportcard.com/report/github.com/tongque0/1aides)
[![OpenIssue](https://img.shields.io/github/issues/tongque0/1aides)](https://github.com/tongque0/1aides/issues)
[![ClosedIssue](https://img.shields.io/github/issues-closed/tongque0/1aides)](https://github.com/tongque0/1aides/issues?q=is%3Aissue+is%3Aclosed)
![Stars](https://img.shields.io/github/stars/tongque0/1aides)

## 快速开始

要启动该微信机器人，只需进入 `deploy` 目录并运行以下命令：
```bash
cd deploy && docker-compose up -d
```

## 特点

- 自带后台管理页面
- 自带命令行管理
- 支持定时任务
- 支持多模型
- 支持记忆，多轮对话

## 贡献代码

提交代码需要符合提交规范，提交附带标记会自动更新此版本或主版本，如：
- **fix\***: 自动将次版本号加一，v1.0.1 -> v1.1.0
- **feat\*\***: 自动将主版本号加一 v1.0.1 -> v2.0.0

### 规范分类

- **feat**: 新功能（feature）
- **fix**: 错误修复
- **docs**: 文档更改（documentation）
- **style**: 格式（不影响代码含义的更改，空格、格式、缺少分号等）
- **refactor**: 重构（即不是新功能，也不是修补bug的代码变动）
- **perf**: 优化（提高性能的代码更改）
- **test**: 测试（添加缺失的测试或更正现有测试）
- **chore**: 构建过程或辅助工具的变动
- **revert**: 还原以前的提交

## 开源许可

1aides 基于 [Apache License 2.0](https://github.com/tongque0/1aides/LICENSE) 许可证，其依赖的三方组件的开源许可见 [Here](https://github.com/eatmoreapple/openwechat/blob/master/LICENSE)。

