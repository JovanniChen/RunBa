import axios from 'axios'

// API配置管理
export const API_CONFIG = {
  LOCAL: '/api/v1',
  REMOTE: 'http://8.223.15.182:8080/api/v1'
}

// 获取当前API基础URL
const getBaseURL = (): string => {
  if (typeof window !== 'undefined') {
    const savedConfig = localStorage.getItem('api_config')
    if (savedConfig) {
      return savedConfig
    }
  }
  return API_CONFIG.LOCAL
}

// 设置API配置
export const setAPIConfig = (config: string) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('api_config', config)
    // 重新创建axios实例
    updateAxiosBaseURL(config)
    // 刷新页面以确保配置生效
    window.location.reload()
  }
}

// 初始化API配置
export const initializeAPIConfig = () => {
  if (typeof window !== 'undefined') {
    const savedConfig = getCurrentAPIConfig()
    updateAxiosBaseURL(savedConfig)
  }
}

// 获取当前API配置
export const getCurrentAPIConfig = (): string => {
  return getBaseURL()
}

const baseURL = getBaseURL()

// 创建axios实例
const request = axios.create({
  baseURL,
  timeout: 30000,
  // 自定义 JSON 解析器处理大整数
  transformResponse: [function (data) {
    if (typeof data === 'string') {
      try {
        // 使用正则表达式将大整数转换为字符串
        let processedData = data
        // 处理外层字段
        processedData = processedData.replace(/"steam_id":\s*(\d{15,})/g, (match, p1) => `"steam_id":"${p1}"`)
        processedData = processedData.replace(/"sender_steam_id":\s*(\d{15,})/g, (match, p1) => `"sender_steam_id":"${p1}"`)
        processedData = processedData.replace(/"receiver_steam_id":\s*(\d{15,})/g, (match, p1) => `"receiver_steam_id":"${p1}"`)
        processedData = processedData.replace(/"send_steam_id":\s*(\d{15,})/g, (match, p1) => `"send_steam_id":"${p1}"`)
        processedData = processedData.replace(/"target_steam_id":\s*(\d{15,})/g, (match, p1) => `"target_steam_id":"${p1}"`)
        processedData = processedData.replace(/"SteamID":\s*(\d{15,})/g, (match, p1) => `"SteamID":"${p1}"`)
        processedData = processedData.replace(/"steamid":\s*(\d{15,})/g, (match, p1) => `"steamid":"${p1}"`)
        processedData = processedData.replace(/"steam64":\s*(\d{15,})/g, (match, p1) => `"steam64":"${p1}"`)
        processedData = processedData.replace(/"account_id":\s*(\d{10,})/g, (match, p1) => `"account_id":"${p1}"`)
        processedData = processedData.replace(/"user_id":\s*(\d{10,})/g, (match, p1) => `"user_id":"${p1}"`)
        return JSON.parse(processedData)
      } catch (e) {
        return data
      }
    }
    return data
  }]
})

// 创建不需要认证的请求实例
const publicRequest = axios.create({
  baseURL,
  timeout: 30000,
  transformResponse: [function (data) {
    if (typeof data === 'string') {
      try {
        let processedData = data
        processedData = processedData.replace(/"steam_id":\s*(\d{15,})/g, (match, p1) => `"steam_id":"${p1}"`)
        processedData = processedData.replace(/"sender_steam_id":\s*(\d{15,})/g, (match, p1) => `"sender_steam_id":"${p1}"`)
        processedData = processedData.replace(/"receiver_steam_id":\s*(\d{15,})/g, (match, p1) => `"receiver_steam_id":"${p1}"`)
        processedData = processedData.replace(/"send_steam_id":\s*(\d{15,})/g, (match, p1) => `"send_steam_id":"${p1}"`)
        processedData = processedData.replace(/"target_steam_id":\s*(\d{15,})/g, (match, p1) => `"target_steam_id":"${p1}"`)
        processedData = processedData.replace(/"SteamID":\s*(\d{15,})/g, (match, p1) => `"SteamID":"${p1}"`)
        processedData = processedData.replace(/"steamid":\s*(\d{15,})/g, (match, p1) => `"steamid":"${p1}"`)
        processedData = processedData.replace(/"steam64":\s*(\d{15,})/g, (match, p1) => `"steam64":"${p1}"`)
        processedData = processedData.replace(/"account_id":\s*(\d{10,})/g, (match, p1) => `"account_id":"${p1}"`)
        processedData = processedData.replace(/"user_id":\s*(\d{10,})/g, (match, p1) => `"user_id":"${p1}"`)
        return JSON.parse(processedData)
      } catch (e) {
        return data
      }
    }
    return data
  }]
})

// 动态更新axios实例的baseURL
export const updateAxiosBaseURL = (newBaseURL: string) => {
  request.defaults.baseURL = newBaseURL
  publicRequest.defaults.baseURL = newBaseURL
}

// 请求拦截器
request.interceptors.request.use(
  config => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token')
      if (token) {
        // 使用与Vue项目相同的token头格式
        config.headers.token = token
      }
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response?.status === 401) {
      if (typeof window !== 'undefined') {
        // 只有在非登录页面时才清除token和重定向
        const isLoginPage = window.location.pathname === '/login'
        if (!isLoginPage) {
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          window.location.href = '/login'
        }
      }
    }
    return Promise.reject(error.response?.data || error)
  }
)

// 公共请求拦截器
publicRequest.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 公共响应拦截器
publicRequest.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    return Promise.reject(error.response?.data || error)
  }
)

// 通用API响应类型
export interface ApiResponse<T = any> {
  code: number
  success: boolean
  message: string
  data: T
}

// 分页响应数据类型
export interface PaginatedData<T = any> {
  rows: T[]
  page: number
  page_size: number
  total: number
  totalpage: number
}

// 分页响应类型
export type PaginatedResponse<T = any> = ApiResponse<PaginatedData<T>>

// 类型定义
export interface LoginData {
  account: string
  password: string
}

export interface RegisterData {
  username: string
  email: string
  password: string
  nickname?: string
  phone?: string
}

export interface User {
  id: string
  username: string
  email: string
  nickname?: string
  phone?: string
  status: boolean
  last_login?: string
  created_at: string
  updated_at: string
}

export interface Certificate {
  id: number
  number?: string
  overall_length?: string
  blade_length?: string
  handle_length?: string
  blade_material?: string
  blade_thickness?: string
  kissan_length?: string
  hamon?: string
  mekugi?: number
  sakihaba?: string
  motohaba?: string
  saya_material?: string
  tsuba_material?: string
  habaki_material?: string
  ito_sageo_material?: string
  forge_id?: number
  swordsmith_id?: number
  forge_name?: string
  swordsmith_name?: string
  forge_image_url?: string
  swordsmith_image_url?: string
  date_completed?: string
  no?: string
  created_at?: string
}

export interface ActivationCode {
  id: number
  name: string
  code: string
  points: number
  status: number // 1=正常, 2=已使用, 3=已废弃
  used_by?: string
  used_at?: string
  created_at: string
  updated_at: string
}

export interface Proxy {
  id: string
  host: string
  port: number
  protocol: number // 1=HTTP, 2=HTTPS, 3=SOCKS5
  category: number // 1=登录, 2=支付
  username?: string
  password?: string
  status: number // 1=正常, 2=禁用
  remark?: string
  created_at: string
}

export interface Config {
  id: string
  key?: string
  value?: string
  description?: string
  is_active?: boolean
  tutorial_url?: string
  created_at?: string
  updated_at?: string
}

export interface Forge {
  id: number
  name: string
  image_url?: string
  created_at?: string
}

export interface Swordsmith {
  id: number
  name: string
  image_url?: string
  created_at?: string
}

// 认证相关API
export const authAPI = {
  login: (data: LoginData) => request.post('/auth/login', data),
  register: (data: RegisterData) => request.post('/auth/register', data),
  getProfile: () => request.get('/auth/profile'),
  updateProfile: (data: Record<string, unknown>) => request.put('/auth/profile', data),
  changePassword: (data: { oldPassword: string; newPassword: string }) => request.post('/auth/change-password', data)
}

// 用户管理相关API
export const userAPI = {
  getUsers: (params?: Record<string, unknown>) => request.get('/users/getUsers', { params }),
  getUser: (id: string) => request.get('/users/getUserByID', { params: { id } }),
  createUser: (data: Partial<User>) => request.post('/users/createUser', data),
  updateUser: (id: string, data: Partial<User>) => request.post('/users/updateUser', { ...data, id }),
  deleteUser: (id: string) => request.get('/users/deleteUserByID', { params: { id } }),
  updateUserStatus: (id: string, status: boolean) =>
    request.patch(`/users/${id}/status`, { status }),
  resetPassword: (id: string) => request.post(`/users/${id}/reset-password`),
  getUserStats: () => request.get('/users/stats')
}

// 证书管理相关API
export const certificateAPI = {
  getCertificates: (params?: Record<string, unknown>) => request.get('/certificates/getCertificates', { params }),
  getCertificate: (id: string | number) => request.get('/certificates/getCertificateByID', { params: { id } }),
  createCertificate: (data: Partial<Certificate>) => request.post('/certificates/createCertificate', data),
  updateCertificate: (id: string | number, data: Partial<Certificate>) =>
    request.post('/certificates/updateCertificate', { ...data, id }),
  deleteCertificate: (id: string | number) =>
    request.get('/certificates/deleteCertificateByID', { params: { id } })
}

// 证书公共查询API
export const certificatePublicAPI = {
  getCertificateByNo: (no: string) =>
    publicRequest.get<ApiResponse<Certificate>, ApiResponse<Certificate>>(
      '/certificates/public',
      { params: { no } }
    )
}

// 锻刀所管理相关API
export const forgeAPI = {
  getForges: (params?: Record<string, unknown>) => request.get('/forges/getForges', { params }),
  getForge: (id: string | number) => request.get('/forges/getForgeByID', { params: { id } }),
  createForge: (data: Partial<Forge>) => request.post('/forges/createForge', data),
  updateForge: (id: string | number, data: Partial<Forge>) =>
    request.post('/forges/updateForge', { ...data, id }),
  deleteForge: (id: string | number) => request.get('/forges/deleteForgeByID', { params: { id } })
}

// 刀匠管理相关API
export const swordsmithAPI = {
  getSwordsmiths: (params?: Record<string, unknown>) => request.get('/swordsmiths/getSwordsmiths', { params }),
  getSwordsmith: (id: string | number) => request.get('/swordsmiths/getSwordsmithByID', { params: { id } }),
  createSwordsmith: (data: Partial<Swordsmith>) => request.post('/swordsmiths/createSwordsmith', data),
  updateSwordsmith: (id: string | number, data: Partial<Swordsmith>) =>
    request.post('/swordsmiths/updateSwordsmith', { ...data, id }),
  deleteSwordsmith: (id: string | number) =>
    request.get('/swordsmiths/deleteSwordsmithByID', { params: { id } })
}

// 工具函数：处理分页数据响应
export function extractPaginatedData<T>(response: PaginatedResponse<T>) {
  if (!response.success) {
    throw new Error(response.message || '请求失败')
  }
  return {
    list: response.data.rows,
    pagination: {
      page: response.data.page,
      pageSize: response.data.page_size,
      total: response.data.total,
      totalPages: response.data.totalpage
    }
  }
}

// 工具函数：处理普通API响应
export function extractApiData<T>(response: ApiResponse<T>) {
  if (!response.success) {
    throw new Error(response.message || '请求失败')
  }
  return response.data
}

export default request
