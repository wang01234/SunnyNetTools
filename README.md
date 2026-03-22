# DS 大神工具 - 网络抓包工具

[English](./README_EN.md) | 中文

## 简介

DS 大神工具是一款功能强大的网络抓包和调试工具，基于 [SunnyNet](https://github.com/qtgolang/SunnyNet) 开发，使用 Wails + Vue 3 构建的跨平台桌面应用。

支持 HTTP、HTTPS、TCP、UDP、WebSocket 等多种网络协议的抓包和分析。

## 功能特性

### 核心功能

- **多协议抓包** - 支持 HTTP、HTTPS、TCP、UDP、WebSocket 协议
- **进程拦截** - 按进程名或 PID 精准拦截指定应用的网络请求
- **证书管理** - 内置 SSL 证书生成，支持 Windows、macOS、Android、iOS 全平台
- **Hosts 映射** - 本地 Hosts 规则配置
- **请求/响应替换** - 支持请求和响应内容的动态替换
- **Go 脚本** - 支持使用 Go 脚本编写自定义抓包逻辑
- **代理设置** - 全局代理和上游代理配置
- **强制 TCP** - 强制指定域名走 TCP 协议

### 特色功能

- **进程名持久化** - 拦截的进程名列表自动保存，无需每次手动配置
- **多主题支持** - 支持明暗主题切换
- **账号管理** - 支持多用户登录和权限管理
- **响应高亮** - 支持 JSON、XML、HTML 等格式自动高亮显示
- **文本对比** - 支持请求/响应内容的对比分析

## 技术栈

- **后端**: Go + [SunnyNet](https://github.com/qtgolang/SunnyNet)
- **前端**: Vue 3 + TypeScript + Element Plus
- **桌面框架**: Wails
- **编辑器**: Monaco Editor

## 构建

### 环境要求

- Go 1.18+
- Node.js 16+
- Wails CLI

### 构建命令

```bash
# 安装依赖
go mod download

# 前端依赖
cd frontend && npm install

# 开发模式
wails dev

# 生产构建 (macOS)
wails build -platform darwin/amd64

# 生产构建 (Windows)
wails build -platform windows/amd64

# 生产构建 (Linux)
wails build -platform linux/amd64
```

## 下载安装

### macOS

[下载 macOS 版本](https://wwxa.lanzouj.com/b0ciopv1c) 密码: 2oxf

### Windows

[下载 Windows 版本](https://wwxa.lanzouj.com/b0cior9kb) 密码: 2brf

### 从源码构建

参见上方 [构建](#构建) 部分

## 界面预览

<img src="./img/1.jpg">
<img src="./img/2.jpg">

## 许可证

MIT License

## 致谢

- [SunnyNet](https://github.com/qtgolang/SunnyNet) - 核心网络中间件
- [Wails](https://wails.io/) - 桌面应用框架
- [Vue 3](https://vuejs.org/) - 前端框架
- [Element Plus](https://element-plus.org/) - UI 组件库
