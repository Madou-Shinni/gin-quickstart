# 获取当前版本号，如果没有则默认为 v0.0.0
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# 解析版本号
VERSION_MAJOR := $(shell echo $(CURRENT_VERSION) | cut -d. -f1 | sed 's/v//')
VERSION_MINOR := $(shell echo $(CURRENT_VERSION) | cut -d. -f2)
VERSION_PATCH := $(shell echo $(CURRENT_VERSION) | cut -d. -f3)

# 计算下一个版本号
NEXT_PATCH := $(shell echo $$(($(VERSION_PATCH) + 1)))
NEXT_VERSION := v$(VERSION_MAJOR).$(VERSION_MINOR).$(NEXT_PATCH)

.PHONY: tag
# 使用方式：
# make tag                 # 自动递增补丁版本 (v0.0.0 → v0.0.1)
# make tag VERSION=v1.2.3  # 手动指定版本号
tag:
ifndef VERSION
	@echo "当前版本: $(CURRENT_VERSION)"
	@echo "自动生成下一个版本: $(NEXT_VERSION)"
	$(eval VERSION = $(NEXT_VERSION))
else
	@echo "使用手动指定版本号: $(VERSION)"
endif
	@# 校验版本号格式
	@echo $(VERSION) | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' >/dev/null || (echo "版本号格式错误，需符合 vX.Y.Z"; exit 1)
	@# 创建并推送标签
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

.PHONY: init
# initialize env
init:
	go install github.com/swaggo/swag/cmd/swag@v1.7.9
	go install github.com/Madou-Shinni/gctl@latest

.PHONY: build
# build
# 使用方式：
# make build              # 使用当前版本号构建
# make build VERSION=v1.2.3 # 使用指定版本号构建
build:
ifndef VERSION
	$(eval VERSION = $(CURRENT_VERSION))
	@echo "使用当前版本号构建: $(VERSION)"
else
	@echo "使用指定版本号构建: $(VERSION)"
endif
	@echo "使用 Docker 构建应用程序..."
	docker build -t gin-quickstart:$(VERSION) .

.PHONY: api-sync
# sync api
api-sync:
	swag init && go run cmd/auto/main.go
