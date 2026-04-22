.PHONY: build test clean format docs

# 输出目录
OUTPUT_DIR=output

# 构建目标的设置
all: test build

# 构建 example 下的所有示例并将二进制可执行文件输出到 output/
build:
	@echo "=> 构建 example 项目到 $(OUTPUT_DIR)/ 目录..."
	@mkdir -p $(OUTPUT_DIR)
	@cd example && go build -o ../$(OUTPUT_DIR)/go-kit-example main.go
	@echo "=> 构建完成!"

# 执行所有的单元测试
test:
	@echo "=> 正在执行单元测试..."
	@go test ./... -v

# 清理输出历史和缓存
clean:
	@echo "=> 清理构建输出和临时文件..."
	@rm -rf $(OUTPUT_DIR)/*
	@go clean -cache

# 初始化并补充缺失得必要目录 (根据项目约束)
docs:
	@echo "=> 初始化文档记录目录..."
	@mkdir -p docs/changelog
	@touch docs/changelog/CHANGELOG.md
