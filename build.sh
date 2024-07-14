#!/bin/bash

go mod tidy

go build easyssh.go

# 执行命令并将输出保存到变量中
output=$(./easyssh version)

# 使用awk提取版本号
version=$(echo "$output" | awk -F ': ' '/Version/ {print $2}')

# 输出结果
echo "版本号是: $version"

echo "开始构建Linux版本"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easyssh_linux_amd64_${version} easyssh.go

#echo "开始构建Windows版本"
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o easyssh_windows_amd64_${version}.exe easyssh.go

echo "开始构建Mac版本"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o easyssh_darwin_amd64_${version} easyssh.go

rm -f easyssh

echo "构建完成"
