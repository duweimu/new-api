# 一键重打补丁镜像脚本设计

## 背景

当前 `20pro.cn` 上的 `new-api` 采用 Docker 镜像部署。为了让本地补丁在后续官方镜像升级后可以快速重放，需要一套从本地 Windows 电脑发起的“一键重打补丁镜像”脚本。

这个仓库的后端二进制会嵌入前端产物：

- `web/default/dist`
- `web/classic/dist`

因此脚本不能只执行 `go build`，还必须确保前端构建产物存在且与当前代码一致。

## 目标

提供一套面向本地 Windows 环境的单命令部署流程，完成以下工作：

1. 校验本地依赖和配置
2. 构建前端产物
3. 执行补丁相关的最小测试
4. 编译 Linux `amd64` 后端二进制
5. 上传二进制到目标服务器
6. 在目标服务器基于当前运行镜像重打一个新的补丁镜像
7. 更新 `docker-compose.yml` 中 `new-api` 服务的镜像标签
8. 重启 `new-api` 容器并等待健康检查
9. 失败时自动回滚 compose 中的镜像标签

## 非目标

- 不负责首次安装服务器环境
- 不负责初始化数据库、域名或证书
- 不负责自动同步 Git 远端仓库
- 不把服务器密码写入版本库
- 不替代完整 CI/CD

## 推荐方案

采用“`PowerShell` 入口 + `Python(paramiko)` 远程执行”的双文件方案：

- `bin/rebuild-patched-image.ps1`
- `bin/rebuild_patched_image.py`

原因：

1. 用户在 Windows 上直接运行 `.ps1` 更顺手
2. `paramiko` 已经在当前工作流中验证可用，适合密码 SSH 场景
3. Python 更适合处理上传、远程执行、日志采集和失败回滚

## 文件设计

### 1. 入口脚本

文件：`bin/rebuild-patched-image.ps1`

职责：

- 解析参数
- 定位仓库根目录
- 检查 `bun`、`go`、`python` 可用性
- 调用前端构建
- 调用测试
- 调用 Go 编译
- 调用 Python 部署脚本

### 2. 远程部署脚本

文件：`bin/rebuild_patched_image.py`

职责：

- 读取配置文件
- 通过 `paramiko` 连接目标服务器
- 上传二进制
- 生成远端临时 `Dockerfile`
- 构建补丁镜像
- 备份并更新 `docker-compose.yml`
- 执行 `docker compose up -d new-api`
- 轮询健康检查
- 失败时回滚镜像标签并再次拉起服务

### 3. 配置模板

文件：`bin/rebuild-patched-image.example.json`

包含字段：

- `host`
- `port`
- `username`
- `password`
- `deploy_dir`
- `service_name`
- `status_url`
- `base_image`
- `remote_build_root`

### 4. 本地私有配置

文件：`bin/rebuild-patched-image.local.json`

说明：

- 不进入 Git
- 用户本地存放真实连接信息
- 脚本优先读取这个文件

## 命令行设计

### PowerShell 脚本参数

- `-ConfigPath`
- `-ImageTag`
- `-SkipFrontendBuild`
- `-SkipTests`
- `-SkipStatusCheck`
- `-DryRun`

默认行为：

- 若未指定 `-ImageTag`，自动生成形如 `new-api:patched-YYYYMMDD-HHmmss`

## 处理流程

### 本地阶段

1. 读取配置
2. 检查依赖
3. 如未指定 `-SkipFrontendBuild`：
   - `web/default` 执行 `bun install` 与 `bun run build`
   - `web/classic` 执行 `bun install` 与 `bun run build`
4. 如未指定 `-SkipTests`：
   - 运行面向当前补丁的 smoke tests
5. 执行 Linux `amd64` 编译，生成临时二进制

### 远端阶段

1. 创建远端临时目录
2. 上传本地编译好的 `new-api` 二进制
3. 生成远端 `Dockerfile`
4. 基于 `base_image` 或当前 compose 中正在使用的 `new-api` 镜像构建新镜像
5. 备份 `docker-compose.yml`
6. 仅替换 `new-api` 服务的 `image` 字段
7. 重启 `new-api` 服务
8. 轮询容器健康状态
9. 请求 `status_url` 进行外部可用性检查

## 回滚策略

当以下任一阶段失败时触发自动回滚：

- 远端镜像构建失败
- `docker compose up -d new-api` 失败
- 容器在限定时间内未进入 `healthy`
- `status_url` 返回异常

回滚步骤：

1. 使用备份文件恢复 `docker-compose.yml`
2. 执行 `docker compose up -d new-api`
3. 输出失败阶段和回滚结果

## 日志与输出

脚本输出应覆盖：

- 当前镜像标签
- 新镜像标签
- 本地编译产物路径
- 远端临时目录
- compose 备份路径
- 健康检查状态
- 外部状态检查结果

`-DryRun` 模式下只打印将要执行的步骤，不执行上传、构建和重启。

## 安全要求

- 不把密码、域名、私有配置写入版本库
- 示例配置只保留占位信息
- 本地配置文件加入 `.gitignore`
- PowerShell 脚本打印配置时隐藏密码
- 远端临时目录应可重复使用，但不删除用户其他目录

## 测试策略

至少覆盖以下几类验证：

1. 参数解析正常
2. 缺少配置时给出明确报错
3. `-DryRun` 不执行真实部署
4. 指定 `-SkipFrontendBuild` 时跳过前端构建
5. 远端健康检查失败时触发回滚

本次实现先以手工验证为主，自动化测试优先覆盖 Python 中纯函数部分，例如配置解析和 compose 镜像替换逻辑。

## 与当前补丁链的关系

该脚本默认服务于当前分支上的两轮补丁：

1. `Fix Codex responses streaming compatibility`
2. `Fix Claude stream termination with cached usage`

但脚本本身不绑定这两次补丁内容。只要当前工作区能成功构建，脚本就可以作为通用的“本地编译并远端重打 `new-api` 镜像”的工具使用。

## 成功标准

用户执行一条 PowerShell 命令后，可以在本地完成：

1. 构建当前仓库版本
2. 把二进制部署到目标服务器
3. 将 `new-api` 容器切换到新的补丁镜像
4. 看到清晰的成功或失败输出
5. 在失败时自动回滚到原镜像
