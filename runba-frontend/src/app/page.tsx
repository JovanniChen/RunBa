'use client'

import { useEffect } from 'react'
import { useAuth } from '@/lib/auth'

export default function Home() {
  const { isLoggedIn, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading) {
      if (isLoggedIn) {
        window.location.href = '/users'
      } else {
        window.location.href = '/login'
      }
    }
  }, [isLoggedIn, isLoading])

  return (
    <div className="flex items-center justify-center h-screen">
      <div className="text-center">
        <h1 className="text-2xl font-bold mb-4">Steam 管理系统</h1>
        <p className="text-gray-600">加载中...</p>
      </div>
    </div>
  )
}
