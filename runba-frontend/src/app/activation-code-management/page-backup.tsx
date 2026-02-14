'use client'

import { ProtectedRoute } from '@/lib/auth'
import { DashboardLayout } from '@/components/dashboard-layout'

export default function ActivationCodeManagementPage() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold">激活码管理</h1>
              <p className="text-muted-foreground">管理系统激活码</p>
            </div>
          </div>
          
          <div className="text-center py-8">
            <p className="text-muted-foreground">激活码管理功能正在开发中...</p>
          </div>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}