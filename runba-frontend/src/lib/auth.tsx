'use client'

import React, { createContext, useContext, useEffect, useState } from 'react'
import { authAPI, initializeAPIConfig } from './api'

interface User {
  id: string
  username: string
  email: string
  [key: string]: any
}

interface AuthContextType {
  user: User | null
  token: string | null
  isLoggedIn: boolean
  login: (userInfo: User, token: string) => void
  logout: () => void
  updateUser: (userInfo: Partial<User>) => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  const login = (userInfo: User, token: string) => {
    setUser(userInfo)
    setToken(token)
    setIsLoggedIn(true)
    
    // 持久化到本地存储
    localStorage.setItem('userInfo', JSON.stringify(userInfo))
    localStorage.setItem('token', token)
  }

  const logout = () => {
    setUser(null)
    setToken(null)
    setIsLoggedIn(false)
    
    // 清除本地存储
    localStorage.removeItem('userInfo')
    localStorage.removeItem('token')
  }

  const updateUser = (userInfo: Partial<User>) => {
    if (user) {
      const updatedUser = { ...user, ...userInfo }
      setUser(updatedUser)
      localStorage.setItem('userInfo', JSON.stringify(updatedUser))
    }
  }

  // 初始化用户信息
  useEffect(() => {
    const initAuth = () => {
      try {
        // 初始化API配置
        initializeAPIConfig()
        const storedUserInfo = localStorage.getItem('userInfo')
        const storedToken = localStorage.getItem('token')
        
        if (storedUserInfo && storedToken) {
          const userInfo = JSON.parse(storedUserInfo)
          setUser(userInfo)
          setToken(storedToken)
          setIsLoggedIn(true)
        }
      } catch (error) {
        console.error('恢复用户信息失败:', error)
        logout()
      } finally {
        setIsLoading(false)
      }
    }

    initAuth()
  }, [])

  const value = {
    user,
    token,
    isLoggedIn,
    login,
    logout,
    updateUser,
    isLoading
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

// 路由守卫组件
export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isLoggedIn, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading && !isLoggedIn) {
      window.location.href = '/login'
    }
  }, [isLoggedIn, isLoading])

  if (isLoading) {
    return <div className="flex items-center justify-center h-screen">加载中...</div>
  }

  if (!isLoggedIn) {
    return null
  }

  return <>{children}</>
}

// 公共路由守卫组件（已登录用户重定向）
export function PublicRoute({ children }: { children: React.ReactNode }) {
  const { isLoggedIn, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading && isLoggedIn) {
      window.location.href = '/users'
    }
  }, [isLoggedIn, isLoading])

  if (isLoading) {
    return <div className="flex items-center justify-center h-screen">加载中...</div>
  }

  if (isLoggedIn) {
    return null
  }

  return <>{children}</>
}