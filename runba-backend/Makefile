# Makefile for Steam Backend API

# 变量定义
APP_NAME := steam-backend
BUILD_DIR := build
MAIN_FILE := main.go
DOCKER_IMAGE := steam-backend:latest

# Go 相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
BINARY_NAME := $(APP_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix

# 默认目标
.PHONY: all
all: clean build

# 构建应用
.PHONY: build
build:
	@echo "构建应用..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 运行应用
.PHONY: run
run:
	@echo "启动应用..."
	$(GOCMD) run $(MAIN_FILE)

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...