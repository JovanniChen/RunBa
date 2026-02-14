'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import Link from 'next/link'
import { useAuth, PublicRoute } from '@/lib/auth'
import { authAPI } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { toast } from '@/components/ui/toast'

const loginSchema = z.object({
  account: z.string().min(1, '请输入账户名'),
  password: z.string().min(1, '请输入密码'),
})

type LoginForm = z.infer<typeof loginSchema>

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false)
  const { login } = useAuth()

  const form = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      account: '',
      password: '',
    },
  })

  const onSubmit = async (data: LoginForm) => {
    try {
      setIsLoading(true)
      
      const response = await authAPI.login(data) as any
      console.log('Login response:', response) // 添加调试日志
      
      // 按照Vue项目的处理方式：直接访问 response.data.token 和 response.data.user
      if (response && response.data && response.data.token && response.data.user) {
        const token = response.data.token
        const user = response.data.user
        login(user, token)
        window.location.href = '/users'
      } else {
        toast.error('登录响应格式错误', '登录失败')
        console.error('Unexpected response format:', response)
      }
    } catch (err: any) {
      console.error('Login error:', err) // 添加调试日志
      toast.error(err.message || '登录失败，请检查账户名和密码', '登录失败')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <PublicRoute>
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div className="text-center">
            <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
              Steam 管理系统
            </h2>
            <p className="mt-2 text-sm text-gray-600">
              请登录您的账户
            </p>
          </div>
          
          <Card>
            <CardHeader>
              <CardTitle>登录</CardTitle>
              <CardDescription>
                输入您的账户名和密码来访问系统
              </CardDescription>
            </CardHeader>
            
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)}>
                <CardContent className="space-y-4">
                  
                  <FormField
                    control={form.control}
                    name="account"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>账户名</FormLabel>
                        <FormControl>
                          <Input 
                            placeholder="请输入账户名" 
                            {...field}
                            disabled={isLoading}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>密码</FormLabel>
                        <FormControl>
                          <Input 
                            type="password" 
                            placeholder="请输入密码" 
                            {...field}
                            disabled={isLoading}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </CardContent>
                
                <CardFooter className="flex flex-col space-y-4">
                  <Button 
                    type="submit" 
                    className="w-full"
                    disabled={isLoading}
                  >
                    {isLoading ? '登录中...' : '登录'}
                  </Button>
                  
                  <div className="text-center text-sm">
                    <span className="text-gray-600">还没有账户？</span>
                    <Link 
                      href="/register" 
                      className="text-blue-600 hover:text-blue-500 ml-1"
                    >
                      注册
                    </Link>
                  </div>
                  
                </CardFooter>
              </form>
            </Form>
          </Card>
        </div>
      </div>
    </PublicRoute>
  )
}
