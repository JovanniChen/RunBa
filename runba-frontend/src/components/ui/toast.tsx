'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'
import { CheckCircle, XCircle, Info, AlertTriangle, X } from 'lucide-react'

interface ToastProps {
  type?: 'success' | 'error' | 'warning' | 'info'
  title?: string
  message: string
  duration?: number
  onClose?: () => void
  className?: string
}

const toastVariants = {
  success: {
    bgColor: 'bg-green-50 border-green-200',
    textColor: 'text-green-800',
    iconColor: 'text-green-600',
    icon: CheckCircle
  },
  error: {
    bgColor: 'bg-red-50 border-red-200',
    textColor: 'text-red-800',
    iconColor: 'text-red-600',
    icon: XCircle
  },
  warning: {
    bgColor: 'bg-yellow-50 border-yellow-200',
    textColor: 'text-yellow-800',
    iconColor: 'text-yellow-600',
    icon: AlertTriangle
  },
  info: {
    bgColor: 'bg-blue-50 border-blue-200',
    textColor: 'text-blue-800',
    iconColor: 'text-blue-600',
    icon: Info
  }
}

export function Toast({ 
  type = 'info', 
  title, 
  message, 
  duration = 3000, 
  onClose,
  className 
}: ToastProps) {
  const [isVisible, setIsVisible] = React.useState(true)
  const variant = toastVariants[type]
  const Icon = variant.icon

  React.useEffect(() => {
    if (duration > 0) {
      const timer = setTimeout(() => {
        setIsVisible(false)
        setTimeout(() => onClose?.(), 300) // 等待动画完成
      }, duration)

      return () => clearTimeout(timer)
    }
  }, [duration, onClose])

  if (!isVisible) return null

  return (
    <div className={cn(
      'fixed top-4 right-4 z-50 w-auto max-w-sm transform transition-all duration-300 ease-in-out',
      isVisible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'
    )}>
      <div className={cn(
        'flex items-start gap-3 rounded-lg border p-4 shadow-lg backdrop-blur-sm',
        variant.bgColor,
        variant.textColor,
        'animate-in slide-in-from-right-full',
        className
      )}>
        <Icon className={cn('h-5 w-5 flex-shrink-0 mt-0.5', variant.iconColor)} />
        
        <div className="flex-1 min-w-0">
          {title && (
            <div className="font-semibold text-sm mb-1">
              {title}
            </div>
          )}
          <div className="text-sm">
            {message}
          </div>
        </div>

        {onClose && (
          <button
            onClick={() => {
              setIsVisible(false)
              setTimeout(() => onClose(), 300)
            }}
            className={cn(
              'flex-shrink-0 rounded-md p-1 hover:bg-black/10 transition-colors',
              variant.iconColor
            )}
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
    </div>
  )
}

// Toast 容器组件
export function ToastContainer() {
  const [toasts, setToasts] = React.useState<Array<{
    id: string
    type?: 'success' | 'error' | 'warning' | 'info'
    title?: string
    message: string
    duration?: number
  }>>([])

  const removeToast = (id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id))
  }

  // 全局 toast 函数
  React.useEffect(() => {
    const showToast = (options: Omit<ToastProps, 'onClose'>) => {
      const id = Math.random().toString(36).substr(2, 9)
      const newToast = { id, ...options }
      setToasts(prev => [...prev, newToast])
    }

    // 挂载到全局
    (window as any).showToast = showToast

    return () => {
      delete (window as any).showToast
    }
  }, [])

  return (
    <>
      {toasts.map((toast, index) => (
        <div key={toast.id} style={{ top: 16 + index * 80 }} className="fixed right-4 z-50">
          <Toast
            {...toast}
            onClose={() => removeToast(toast.id)}
          />
        </div>
      ))}
    </>
  )
}

// 工具函数
export const toast = {
  success: (message: string, title?: string) => {
    (window as any).showToast?.({ type: 'success', message, title })
  },
  error: (message: string, title?: string) => {
    (window as any).showToast?.({ type: 'error', message, title })
  },
  warning: (message: string, title?: string) => {
    (window as any).showToast?.({ type: 'warning', message, title })
  },
  info: (message: string, title?: string) => {
    (window as any).showToast?.({ type: 'info', message, title })
  }
}