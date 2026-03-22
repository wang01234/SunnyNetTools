# DS Big God Tools - Network Packet Capture Tool

[English](./README_EN.md) | [中文](./README.md)

## Introduction

DS Big God Tools is a powerful network packet capture and debugging tool, developed based on [SunnyNet](https://github.com/qtgolang/SunnyNet) and built with Wails + Vue 3 as a cross-platform desktop application.

Supports packet capture and analysis of HTTP, HTTPS, TCP, UDP, WebSocket and other network protocols.

## Features

### Core Features

- **Multi-protocol Capture** - Support for HTTP, HTTPS, TCP, UDP, WebSocket protocols
- **Process Interception** - Precisely intercept network requests from specified applications by process name or PID
- **Certificate Management** - Built-in SSL certificate generation, supports Windows, macOS, Android, iOS
- **Hosts Mapping** - Local Hosts rules configuration
- **Request/Response Replacement** - Support dynamic replacement of request and response content
- **Go Scripting** - Support custom packet capture logic using Go scripts
- **Proxy Settings** - Global proxy and upstream proxy configuration
- **Force TCP** - Force specified domains to use TCP protocol

### Special Features

- **Process Name Persistence** - Intercepted process names are automatically saved, no manual configuration needed each time
- **Multi-theme Support** - Light/Dark theme switching
- **Account Management** - Multi-user login and permission management
- **Response Highlighting** - Automatic syntax highlighting for JSON, XML, HTML and other formats
- **Text Comparison** - Support for comparing request/response content

## Tech Stack

- **Backend**: Go + [SunnyNet](https://github.com/qtgolang/SunnyNet)
- **Frontend**: Vue 3 + TypeScript + Element Plus
- **Desktop Framework**: Wails
- **Editor**: Monaco Editor

## Building

### Requirements

- Go 1.18+
- Node.js 16+
- Wails CLI

### Build Commands

```bash
# Install dependencies
go mod download

# Frontend dependencies
cd frontend && npm install

# Development mode
wails dev

# Production build (macOS)
wails build -platform darwin/amd64

# Production build (Windows)
wails build -platform windows/amd64

# Production build (Linux)
wails build -platform linux/amd64
```

## Download

### macOS

[Download macOS version](https://wwxa.lanzouj.com/b0ciopv1c) Password: 2oxf

### Windows

[Download Windows version](https://wwxa.lanzouj.com/b0cior9kb) Password: 2brf

### Build from Source

See [Building](#building) section above

## Screenshots

<img src="./img/1.jpg">
<img src="./img/2.jpg">

## License

MIT License

## Acknowledgments

- [SunnyNet](https://github.com/qtgolang/SunnyNet) - Core network middleware
- [Wails](https://wails.io/) - Desktop application framework
- [Vue 3](https://vuejs.org/) - Frontend framework
- [Element Plus](https://element-plus.org/) - UI component library
