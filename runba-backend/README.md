# Steam Backend API

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)
![Gin Version](https://img.shields.io/badge/Gin-1.10.1-green.svg)
![GORM Version](https://img.shields.io/badge/GORM-1.30.1-orange.svg)
![MySQL Version](https://img.shields.io/badge/MySQL-8.0+-yellow.svg)
![License](https://img.shields.io/badge/License-MIT-red.svg)

**基于 Golang + Gin + GORM + MySQL 构建的现代化 Steam 集成后端 API**

提供完整的用户认证、Steam 平台集成、积分系统等企业级功能

[项目特性](#项目特性) • [技术栈](#技术栈) • [快速开始](#快速开始) • [API文档](#api-接口) • [部署指南](#部署说明)

</div>

---

## 📋 目录

- [项目概述](#项目概述)
- [项目特性](#项目特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [API接口](#api-接口)
- [数据库设计](#数据库设计)
- [任务调度](#任务调度)
- [开发指南](#开发指南)
- [部署说明](#部署说明)
- [安全考虑](#安全考虑)
- [重构亮点](#重构亮点)
- [常见问题](#常见问题)
- [贡献指南](#贡献指南)

## 🎯 项目概述

Steam Backend API 是一个基于 Go 语言构建的现代化 Web 后端框架，采用**依赖注入容器**和**分层架构**设计，提供完整的用户认证、Steam 平台集成、积分系统和权限管理等功能。项目遵循 RESTful API 设计原则，支持高并发、高可用的企业级应用需求。

