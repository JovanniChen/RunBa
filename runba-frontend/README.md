# runba-frontend

基于Next.js + TailwindCSS + Shadcn/ui构建的现代化Steam后台管理系统。

## 技术栈

- **框架**: Next.js 15 (App Router)
- **样式**: TailwindCSS + Shadcn/ui
- **状态管理**: React Context + useState
- **表单处理**: React Hook Form + Zod
- **HTTP客户端**: Axios
- **图标**: Lucide React

## 功能特性

- 🔐 完整的登录注册功能
- 👥 用户管理（增删改查、状态管理）
- 💳 Steam账号管理（积分操作、令牌管理）
- 📋 激活记录管理（查看、导出）
- 🔑 激活码管理（生成、验证）
- 🌐 代理管理（配置、状态控制）
- ⚙️ 配置管理（系统参数）
- 🎮 Steam CDK公共兑换接口
- 📱 响应式设计，支持移动端

## 开发环境要求

- Node.js >= 18
- npm >= 9

## 安装和运行

### 1. 安装依赖
```bash
npm install
```

### 2. 启动开发服务器
```bash
npm run dev
```

### 3. 构建生产版本
```bash
npm run build
```

### 4. 启动生产服务器
```bash
npm run start
```

## 项目结构

```
runba-frontend/
├── src/
│   ├── app/                    # Next.js App Router页面
│   │   ├── login/             # 登录页面
│   │   ├── register/          # 注册页面
│   │   ├── users/             # 用户管理
│   │   ├── certificates/      # 证书管理
│   │   ├── activation-record/ # 激活记录
│   │   ├── activation-code-management/ # 激活码管理
│   │   ├── proxy-management/  # 代理管理
│   │   ├── configuration-management/   # 配置管理
│   │   └── steam-cdk/         # Steam CDK兑换
│   ├── components/            # 可复用组件
│   │   ├── ui/               # Shadcn/ui组件
│   │   └── dashboard-layout.tsx # 仪表盘布局
│   └── lib/                  # 工具库
│       ├── api.ts            # API客户端
│       ├── auth.tsx          # 认证管理
│       └── utils.ts          # 工具函数
├── components.json           # Shadcn/ui配置
├── tailwind.config.ts       # TailwindCSS配置
└── next.config.js           # Next.js配置
```

## 后端集成

项目需要配合Steam后端API使用：

- 后端地址：http://localhost:8080
- API前缀：/api/v1

### 主要API接口

#### 认证相关
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `GET /api/v1/auth/profile` - 获取用户信息

#### 用户管理
- `GET /api/v1/users` - 获取用户列表
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

#### 账号管理
- `GET /api/v1/accounts` - 获取Steam账号列表
- `POST /api/v1/accounts/:id/points/add` - 增加积分
- `POST /api/v1/accounts/:id/points/deduct` - 扣除积分

#### Steam CDK (公共接口)
- `POST /api/v1/steam-cdk/exchange` - 提交兑换申请
- `GET /api/v1/steam-cdk/records` - 获取兑换记录

## 开发说明

### 认证流程
- JWT tokens存储在localStorage
- 路由守卫自动处理未认证访问
- 401响应自动触发登出

### Steam ID处理
项目包含专门的Steam ID处理逻辑，将大整数转换为字符串以防止精度丢失。

### 路由结构
- 公共路由：`/login`, `/register`, `/steam-cdk`
- 受保护路由：仪表盘下的所有管理页面
- 默认重定向：`/` → `/users`

## 部署

### 构建
```bash
npm run build
```

### 环境变量
创建 `.env.local` 文件：
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

## 开发者

本项目是 runba-frontend 的 Next.js 前端项目。
