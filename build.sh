#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
WEB_DIR="$PROJECT_ROOT/web"
SERVER_DIR="$PROJECT_ROOT/server"
OUTPUT_DIR="$PROJECT_ROOT/build"
BINARY_NAME="todo-zero"

echo "===== Todo-Zero Build Script ====="
echo ""

# 1. 构建前端
echo "[1/3] 构建前端 (Vue3 + Vite)..."
cd "$WEB_DIR"
if [ ! -d "node_modules" ]; then
  echo "  安装前端依赖..."
  npm install
fi
npm run build
echo "  前端构建完成 -> $SERVER_DIR/dist/"
echo ""

# 2. 编译后端
echo "[2/3] 编译后端 (Go)..."
cd "$SERVER_DIR"
mkdir -p "$OUTPUT_DIR"
CGO_ENABLED=0 go build -ldflags="-s -w" -o "$OUTPUT_DIR/$BINARY_NAME" .
echo "  后端编译完成 -> $OUTPUT_DIR/$BINARY_NAME"
echo ""

# 3. 输出结果
echo "[3/3] 构建结果:"
ls -lh "$OUTPUT_DIR/$BINARY_NAME"
echo ""
echo "===== 构建成功! ====="
echo "运行方式: $OUTPUT_DIR/$BINARY_NAME"
echo "示例: $OUTPUT_DIR/$BINARY_NAME -port 8888"
