# 项目名

## 程序要求

- **`go v1.24.3`**
- **`jzero 1.8.3`**
- **`mysql 8.0`**
- **`redis 6.*：`** 不要使用超过6.\*的版本

## 部署

### 原生部署

#### 1. 配置 `etc/etc.yaml`

- 设置mysql连接配置
- 设置redis连接配置
- 填写deepseek apiKey
- 填写AI主持人提示词

#### 2. 配置 `.jzero.yaml`

- 配置 `migrate` 数据库连接串

#### 3. 创建 `desc/story/story.json` 文件

- 准备导入数据，数据格式如下：

```json
[
    {
        "title": "",
        "surface": "",
        "bottom": ""
    }
]
```

#### 4. 构建二进制

- `go build -o turtlesoupd main.go`
- `chmod a+x turtlesoupd`

#### 5. 执行初始化命令

- **初始化数据库：**`jzero migrate up`
- **初始化数据：**`turtlesoupd tools story add`

#### 6. （可选）生成API文档

- `jzero gen swagger`
- 启动服务，访问链接后加上 `/swagger` 后缀即可

#### 7. 启动后台服务

- `turtlesoupd server --config etc/etc.yaml`

访问页面链接 `http://localhost:8001/`

### docker-compose

#### 1. 配置 `etc/etc.docker.yaml`

- 填写deepseek apiKey
- 填写AI主持人提示词

#### 2. 配置 `.jzero.docker.yaml`

#### 3. 创建 `desc/story/story.json` 文件

- 准备导入数据，数据格式如下：

```json
[
    {
        "title": "",
        "surface": "",
        "bottom": ""
    }
]
```

#### 4. 执行部署脚本

- `bash deploy.sh`

#### 5. 维护

- **服务初始化容器 `--profile init_app`：**
    - **初始化mysql容器：**`docker compose --profile init_app run init_mysql`
    - **初始化app数据容器：**`docker compose --profile init_app run init_app_db`
    - **初始化app容器：**`docker compose --profile init_app run init_app`
- **服务启动：**
    - `docker compose up -d`

## 程序接口流程

### 1. 获取活动token

- 活动token有过期时间
- 重新访问接口会新生成一个活动token，已生成的活动token不会立即过期

### 2. 携带活动token派生会话token并获取第一个故事

- 会话token过期时间与活动token隐性绑定，活动token过期，派生的会话token也会过期
- 获取会话token的同时获取第1个故事

### 3. 携带会话token回答问题

- 回答问题共有10次机会，10次回答没有答出答案会自动派发新的故事
- 一个会话token有3次获取故事的次数，第一个故事包含在内
- 获取故事的次数为0则游戏结束

## 开发者

### @GHog
  - Discord @fbjmk97
  - Github @Makia9879
  - Twitter @Makia981
  - EVM network Address：0x26e64ccaf82cecf582c2c4ac537fb7ef4c60507f

### 特别感谢
  - Discord @huahua08261
  - Discord @461605841

