'use client'

import { useState, useEffect } from 'react'
import { ProtectedRoute } from '@/lib/auth'
import { DashboardLayout } from '@/components/dashboard-layout'
import { userAPI, User, extractPaginatedData, PaginatedResponse } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { toast } from '@/components/ui/toast'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogFooter, 
  DialogHeader, 
  DialogTitle 
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { SmartPagination } from '@/components/smart-pagination'
import { Search, RefreshCw, Eye, Edit, Trash2, CheckCircle, XCircle } from 'lucide-react'

export default function UsersPage() {
  const [userList, setUserList] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  
  // 搜索表单状态
  const [searchForm, setSearchForm] = useState({
    keyword: '',
    status: '0' // 0=全部, 1=激活, 2=禁用
  })
  
  // 分页状态
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0
  })
  
  // 对话框状态
  const [dialogVisible, setDialogVisible] = useState(false)
  const [dialogTitle, setDialogTitle] = useState('')
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  
  // 删除确认对话框
  const [deleteDialog, setDeleteDialog] = useState<{
    open: boolean
    user: User | null
  }>({ open: false, user: null })
  
  // 状态切换确认对话框
  const [statusDialog, setStatusDialog] = useState<{
    open: boolean
    user: User | null
    newStatus: boolean
  }>({ open: false, user: null, newStatus: true })

  useEffect(() => {
    loadUsers()
  }, [pagination.page, pagination.pageSize])

  // 加载用户列表
  const loadUsers = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = {
        page: pagination.page,
        page_size: pagination.pageSize,
        keyword: searchForm.keyword || undefined,
        status: searchForm.status !== '0' ? parseInt(searchForm.status) : undefined
      }

      const response = await userAPI.getUsers(params) as unknown as PaginatedResponse<User>
      const { list, pagination: paginationData } = extractPaginatedData(response)
      
      setUserList(list)
      setPagination(prev => ({ 
        ...prev, 
        total: paginationData.total,
        totalPages: paginationData.totalPages
      }))
    } catch (err: any) {
      toast.error('加载用户列表失败: ' + (err.message || ''), '加载失败')
    } finally {
      setLoading(false)
    }
  }

  // 搜索
  const handleSearch = () => {
    setPagination(prev => ({ ...prev, page: 1 }))
    loadUsers()
  }

  // 重置
  const handleReset = () => {
    setSearchForm({ keyword: '', status: '0' })
    setPagination(prev => ({ ...prev, page: 1 }))
    setTimeout(() => {
      loadUsers()
    }, 0)
  }

  // 刷新数据
  const refreshData = () => {
    loadUsers()
    toast.info('数据已刷新', '刷新成功')
  }

  // 查看用户详情
  const handleView = (user: User) => {
    setCurrentUser(user)
    setDialogTitle('用户详情')
    setDialogVisible(true)
  }

  // 编辑用户
  const handleEdit = (user: User) => {
    // TODO: 实现编辑功能
    toast.info('编辑功能开发中...', '提醒')
  }

  // 切换用户状态
  const handleToggleStatus = (user: User) => {
    const newStatus = !user.status
    setStatusDialog({ open: true, user, newStatus })
  }

  const confirmStatusToggle = async () => {
    if (!statusDialog.user) return
    
    try {
      await userAPI.updateUserStatus(statusDialog.user.id, statusDialog.newStatus)
      toast.success(`用户已${statusDialog.newStatus ? '启用' : '禁用'}`, '状态更新')
      loadUsers()
      setStatusDialog({ open: false, user: null, newStatus: true })
    } catch (err: any) {
      toast.error(`${statusDialog.newStatus ? '启用' : '禁用'}失败: ` + (err.message || ''), '状态更新失败')
    }
  }

  // 删除用户
  const handleDelete = (user: User) => {
    setDeleteDialog({ open: true, user })
  }

  const confirmDelete = async () => {
    if (!deleteDialog.user) return
    
    try {
      await userAPI.deleteUser(deleteDialog.user.id)
      toast.success('用户已删除', '删除成功')
      loadUsers()
      setDeleteDialog({ open: false, user: null })
    } catch (err: any) {
      toast.error('删除失败: ' + (err.message || ''), '删除失败')
    }
  }

  // 分页大小变化
  const handleSizeChange = (size: number) => {
    setPagination(prev => ({ ...prev, pageSize: size, page: 1 }))
  }

  // 当前页变化
  const handleCurrentChange = (page: number) => {
    setPagination(prev => ({ ...prev, page }))
  }

  // 格式化日期
  const formatDate = (dateString: string) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`
  }

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-4">
          {/* 搜索和操作栏 */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  <label className="text-sm font-medium whitespace-nowrap">关键词</label>
                  <Input
                    value={searchForm.keyword}
                    onChange={(e) => setSearchForm(prev => ({ ...prev, keyword: e.target.value }))}
                    placeholder="搜索用户名、邮箱、昵称"
                    className="w-48"
                  />
                </div>
                
                <div className="flex items-center space-x-2">
                  <label className="text-sm font-medium whitespace-nowrap">状态</label>
                  <Select
                    value={searchForm.status}
                    onValueChange={(value) => setSearchForm(prev => ({ ...prev, status: value }))}
                  >
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="0">全部</SelectItem>
                      <SelectItem value="1">激活</SelectItem>
                      <SelectItem value="2">禁用</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                
                <Button onClick={handleSearch} className="shrink-0">
                  <Search className="h-4 w-4 mr-2" />
                  搜索
                </Button>
                
                <Button variant="outline" onClick={handleReset} className="shrink-0">
                  <RefreshCw className="h-4 w-4 mr-2" />
                  重置
                </Button>
              </div>
            </CardContent>
          </Card>


          {/* 用户列表 */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>用户列表</CardTitle>
              <Button onClick={refreshData} variant="outline">
                <RefreshCw className="h-4 w-4 mr-2" />
                刷新
              </Button>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-center py-8">
                  <p>加载中...</p>
                </div>
              ) : userList.length === 0 ? (
                <div className="text-center py-8">
                  <p className="text-muted-foreground">暂无用户数据</p>
                </div>
              ) : (
                <>
                  <div className="border rounded-md">
                    <div className="max-h-[60vh] overflow-auto">
                      <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-16 text-center">ID</TableHead>
                        <TableHead className="w-20 text-center">头像</TableHead>
                        <TableHead className="min-w-32 text-center">用户名</TableHead>
                        <TableHead className="min-w-44 text-center">邮箱</TableHead>
                        <TableHead className="min-w-32 text-center">昵称</TableHead>
                        <TableHead className="min-w-32 text-center">手机号</TableHead>
                        <TableHead className="w-20 text-center">状态</TableHead>
                        <TableHead className="min-w-40 text-center">最后登录</TableHead>
                        <TableHead className="min-w-40 text-center">创建时间</TableHead>
                        <TableHead className="w-36 text-center">操作</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {userList.map((user) => (
                        <TableRow key={user.id}>
                          <TableCell className="font-medium text-center">{user.id}</TableCell>
                          
                          <TableCell className="text-center">
                            <Avatar className="w-10 h-10 mx-auto">
                              <AvatarFallback>
                                {user.username.charAt(0).toUpperCase()}
                              </AvatarFallback>
                            </Avatar>
                          </TableCell>
                          
                          <TableCell className="font-medium text-center">{user.username}</TableCell>
                          <TableCell className="text-center">{user.email}</TableCell>
                          <TableCell className="text-center">{user.nickname || '-'}</TableCell>
                          <TableCell className="text-center">{user.phone || '-'}</TableCell>
                          
                          <TableCell className="text-center">
                            <Badge variant={user.status ? "default" : "destructive"}>
                              {user.status ? '激活' : '禁用'}
                            </Badge>
                          </TableCell>
                          
                          <TableCell className="text-center">
                            {user.last_login ? formatDate(user.last_login) : '-'}
                          </TableCell>
                          <TableCell className="text-center">{formatDate(user.created_at)}</TableCell>
                          
                          <TableCell className="text-center">
                            <div className="flex items-center justify-center space-x-1">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleView(user)}
                                className="h-8 w-8 p-0"
                              >
                                <Eye className="h-4 w-4" />
                              </Button>
                              
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleEdit(user)}
                                className="h-8 w-8 p-0"
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                              
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleToggleStatus(user)}
                                className="h-8 w-8 p-0"
                              >
                                {user.status ? (
                                  <XCircle className="h-4 w-4 text-red-500" />
                                ) : (
                                  <CheckCircle className="h-4 w-4 text-green-500" />
                                )}
                              </Button>
                              
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDelete(user)}
                                className="h-8 w-8 p-0"
                              >
                                <Trash2 className="h-4 w-4 text-red-500" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                      </Table>
                    </div>
                  </div>
                  
                  {/* 分页 */}
                  <div className="mt-4 flex justify-end">
                    <SmartPagination
                      currentPage={pagination.page}
                      totalPages={pagination.totalPages}
                      pageSize={pagination.pageSize}
                      totalItems={pagination.total}
                      onPageChange={handleCurrentChange}
                      onPageSizeChange={handleSizeChange}
                      showSizeChanger
                      showTotal
                      pageSizeOptions={[10, 20, 50, 100]}
                    />
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {/* 用户详情对话框 */}
          <Dialog open={dialogVisible} onOpenChange={setDialogVisible}>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>{dialogTitle}</DialogTitle>
              </DialogHeader>
              
              {currentUser && (
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">用户ID:</span>
                      <p className="font-medium">{currentUser.id}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">用户名:</span>
                      <p className="font-medium">{currentUser.username}</p>
                    </div>
                    <div className="col-span-2">
                      <span className="text-muted-foreground">邮箱:</span>
                      <p className="font-medium">{currentUser.email}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">昵称:</span>
                      <p className="font-medium">{currentUser.nickname || '-'}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">手机号:</span>
                      <p className="font-medium">{currentUser.phone || '-'}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">状态:</span>
                      <Badge variant={currentUser.status ? "default" : "destructive"}>
                        {currentUser.status ? '激活' : '禁用'}
                      </Badge>
                    </div>
                    <div>
                      <span className="text-muted-foreground">最后登录:</span>
                      <p className="font-medium">
                        {currentUser.last_login ? formatDate(currentUser.last_login) : '-'}
                      </p>
                    </div>
                    <div className="col-span-2">
                      <span className="text-muted-foreground">创建时间:</span>
                      <p className="font-medium">{formatDate(currentUser.created_at)}</p>
                    </div>
                  </div>
                </div>
              )}
              
              <DialogFooter>
                <Button variant="outline" onClick={() => setDialogVisible(false)}>
                  关闭
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>

          {/* 状态切换确认对话框 */}
          <Dialog open={statusDialog.open} onOpenChange={(open) => setStatusDialog(prev => ({ ...prev, open }))}>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>提示</DialogTitle>
                <DialogDescription>
                  确定要{statusDialog.newStatus ? '启用' : '禁用'}用户 "{statusDialog.user?.username}" 吗？
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <Button 
                  variant="outline" 
                  onClick={() => setStatusDialog({ open: false, user: null, newStatus: true })}
                >
                  取消
                </Button>
                <Button onClick={confirmStatusToggle}>
                  确定
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>

          {/* 删除确认对话框 */}
          <Dialog open={deleteDialog.open} onOpenChange={(open) => setDeleteDialog(prev => ({ ...prev, open }))}>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>危险操作</DialogTitle>
                <DialogDescription>
                  确定要删除用户 "{deleteDialog.user?.username}" 吗？此操作不可恢复！
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <Button 
                  variant="outline" 
                  onClick={() => setDeleteDialog({ open: false, user: null })}
                >
                  取消
                </Button>
                <Button variant="destructive" onClick={confirmDelete}>
                  确定删除
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}