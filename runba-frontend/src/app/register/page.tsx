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

const registerSchema = z.object({
  username: z.string().min(3, '用户名至少3个字符').max(20, '用户名不能超过20个字符'),
  email: z.string().email('请输入有效的邮箱地址'),
  password: z.string().min(6, '密码至少6个字符'),
  confirmPassword: z.string().min(6, '请确认密码'),
}).refine((data) => data.password === data.confirmPassword, {
  message: '两次输入的密码不一致',
  path: ['confirmPassword'],
})

type RegisterForm = z.infer<typeof registerSchema>

export default function RegisterPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const { login } = useAuth()

  const form = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
  })

  const onSubmit = async (data: RegisterForm) => {
    try {
      setIsLoading(true)
      setSuccess(false)
      
      const { confirmPassword, ...registerData } = data
      await authAPI.register(registerData)
      
      // 注册成功，跳转到登录页面（与原始Vue项目保持一致）
      setSuccess(true)
      toast.success('注册成功！正在跳转到登录页面...', '注册成功')
      setTimeout(() => {
        window.location.href = '/login'
      }, 1500)
    } catch (err: any) {
      toast.error(err.message || '注册失败，请重试', '注册失败')
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
              创建您的账户
            </p>
          </div>
          
          <Card>
            <CardHeader>
              <CardTitle>注册</CardTitle>
              <CardDescription>
                填写以下信息来创建新账户
              </CardDescription>
            </CardHeader>
            
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)}>
                <CardContent className="space-y-4">
                  
                  <FormField
                    control={form.control}
                    name="username"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>用户名</FormLabel>
                        <FormControl>
                          <Input 
                            placeholder="请输入用户名" 
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
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>邮箱</FormLabel>
                        <FormControl>
                          <Input 
                            type="email"
                            placeholder="请输入邮箱地址" 
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
                  
                  <FormField
                    control={form.control}
                    name="confirmPassword"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>确认密码</FormLabel>
                        <FormControl>
                          <Input 
                            type="password" 
                            placeholder="请再次输入密码" 
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
                    disabled={isLoading || success}
                  >
                    {isLoading ? '注册中...' : success ? '注册成功' : '注册'}
                  </Button>
                  
                  <div className="text-center text-sm">
                    <span className="text-gray-600">已有账户？</span>
                    <Link 
                      href="/login" 
                      className="text-blue-600 hover:text-blue-500 ml-1"
                    >
                      登录
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