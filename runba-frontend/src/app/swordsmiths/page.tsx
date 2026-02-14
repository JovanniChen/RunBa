'use client'

import { useEffect, useState } from 'react'
import { ProtectedRoute } from '@/lib/auth'
import { DashboardLayout } from '@/components/dashboard-layout'
import { swordsmithAPI, Swordsmith, extractPaginatedData, PaginatedResponse } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { toast } from '@/components/ui/toast'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { SmartPagination } from '@/components/smart-pagination'
import { Search, RefreshCw, Plus, Pencil, Trash2 } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'

const swordsmithFormSchema = z.object({
  name: z.string().min(1, '请输入刀匠名称'),
  image_url: z.string().optional()
})

type SwordsmithFormData = z.infer<typeof swordsmithFormSchema>

const defaultFormValues: SwordsmithFormData = {
  name: '',
  image_url: ''
}

const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  if (Number.isNaN(date.getTime())) return dateString
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`
}

export default function SwordsmithsPage() {
  const [swordsmithList, setSwordsmithList] = useState<Swordsmith[]>([])
  const [loading, setLoading] = useState(false)
  const [searchForm, setSearchForm] = useState({ keyword: '' })
  const [pagination, setPagination] = useState({ page: 1, pageSize: 10, total: 0, totalPages: 0 })

  const [editDialogVisible, setEditDialogVisible] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [saving, setSaving] = useState(false)
  const [currentSwordsmith, setCurrentSwordsmith] = useState<Swordsmith | null>(null)

  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; swordsmith: Swordsmith | null }>({
    open: false,
    swordsmith: null
  })

  const form = useForm<SwordsmithFormData>({
    resolver: zodResolver(swordsmithFormSchema as any),
    defaultValues: defaultFormValues
  })

  useEffect(() => {
    loadSwordsmiths()
  }, [pagination.page, pagination.pageSize])

  const loadSwordsmiths = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = {
        page: pagination.page,
        page_size: pagination.pageSize,
        keyword: searchForm.keyword || undefined
      }

      const response = await swordsmithAPI.getSwordsmiths(params) as unknown as PaginatedResponse<Swordsmith>
      const { list, pagination: paginationData } = extractPaginatedData(response)
      setSwordsmithList(list)
      setPagination(prev => ({
        ...prev,
        total: paginationData.total,
        totalPages: paginationData.totalPages
      }))
    } catch (err: any) {
      toast.error('加载刀匠列表失败: ' + (err.message || ''), '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = () => {
    setPagination(prev => ({ ...prev, page: 1 }))
    loadSwordsmiths()
  }

  const handleReset = () => {
    setSearchForm({ keyword: '' })
    setPagination(prev => ({ ...prev, page: 1 }))
    setTimeout(() => {
      loadSwordsmiths()
    }, 0)
  }

  const refreshData = () => {
    loadSwordsmiths()
    toast.info('数据已刷新', '刷新成功')
  }

  const handleAdd = () => {
    setIsEditing(false)
    setCurrentSwordsmith(null)
    form.reset(defaultFormValues)
    setEditDialogVisible(true)
  }

  const handleEdit = (swordsmith: Swordsmith) => {
    setIsEditing(true)
    setCurrentSwordsmith(swordsmith)
    form.reset({
      name: swordsmith.name,
      image_url: swordsmith.image_url || ''
    })
    setEditDialogVisible(true)
  }

  const handleSaveSwordsmith = async (data: SwordsmithFormData) => {
    try {
      setSaving(true)
      const submitData = {
        name: data.name.trim(),
        image_url: data.image_url?.trim() || ''
      }

      if (isEditing) {
        if (!currentSwordsmith) {
          toast.error('未选择要编辑的刀匠', '更新失败')
          return
        }
        await swordsmithAPI.updateSwordsmith(currentSwordsmith.id, submitData)
        toast.success('刀匠已更新', '更新成功')
      } else {
        await swordsmithAPI.createSwordsmith(submitData)
        toast.success('刀匠已添加', '添加成功')
      }

      setEditDialogVisible(false)
      loadSwordsmiths()
    } catch (err: any) {
      toast.error(`${isEditing ? '更新' : '添加'}失败: ` + (err.message || ''), `${isEditing ? '更新' : '添加'}失败`)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = (swordsmith: Swordsmith) => {
    setDeleteDialog({ open: true, swordsmith })
  }

  const confirmDelete = async () => {
    if (!deleteDialog.swordsmith) return

    try {
      await swordsmithAPI.deleteSwordsmith(deleteDialog.swordsmith.id)
      toast.success('刀匠已删除', '删除成功')
      loadSwordsmiths()
      setDeleteDialog({ open: false, swordsmith: null })
    } catch (err: any) {
      toast.error('删除失败: ' + (err.message || ''), '删除失败')
    }
  }

  const handleSizeChange = (size: number) => {
    setPagination(prev => ({ ...prev, pageSize: size, page: 1 }))
  }

  const handleCurrentChange = (page: number) => {
    setPagination(prev => ({ ...prev, page }))
  }

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-4">
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  <label className="text-sm font-medium whitespace-nowrap">关键词</label>
                  <Input
                    value={searchForm.keyword}
                    onChange={(event) => setSearchForm(prev => ({ ...prev, keyword: event.target.value }))}
                    placeholder="搜索刀匠名称"
                    className="w-64"
                  />
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

          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>刀匠列表</CardTitle>
              <div className="flex space-x-2">
                <Button onClick={handleAdd}>
                  <Plus className="h-4 w-4 mr-2" />
                  添加刀匠
                </Button>
                <Button onClick={refreshData} variant="outline">
                  <RefreshCw className="h-4 w-4 mr-2" />
                  刷新
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-center py-8">
                  <p>加载中...</p>
                </div>
              ) : swordsmithList.length === 0 ? (
                <div className="text-center py-8">
                  <p className="text-muted-foreground">暂无刀匠数据</p>
                </div>
              ) : (
                <>
                  <div className="border rounded-md">
                    <div className="max-h-[60vh] overflow-auto">
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead className="w-16 text-center">ID</TableHead>
                            <TableHead className="min-w-40 text-center">名称</TableHead>
                            <TableHead className="min-w-40 text-center">图片</TableHead>
                            <TableHead className="min-w-40 text-center">创建时间</TableHead>
                            <TableHead className="w-32 text-center">操作</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {swordsmithList.map((swordsmith) => (
                            <TableRow key={swordsmith.id}>
                              <TableCell className="font-medium text-center">{swordsmith.id}</TableCell>
                              <TableCell className="text-center">{swordsmith.name}</TableCell>
                              <TableCell className="text-center">
                                {swordsmith.image_url ? (
                                  <img
                                    src={swordsmith.image_url}
                                    alt={swordsmith.name}
                                    className="mx-auto h-10 w-10 rounded-md object-cover"
                                  />
                                ) : (
                                  <span className="text-muted-foreground">-</span>
                                )}
                              </TableCell>
                              <TableCell className="text-center">{formatDate(swordsmith.created_at)}</TableCell>
                              <TableCell className="text-center">
                                <div className="flex items-center justify-center space-x-1">
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handleEdit(swordsmith)}
                                    className="h-8 w-8 p-0"
                                  >
                                    <Pencil className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handleDelete(swordsmith)}
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

          <Dialog open={editDialogVisible} onOpenChange={setEditDialogVisible}>
            <DialogContent className="max-w-lg">
              <DialogHeader>
                <DialogTitle>{isEditing ? '编辑刀匠' : '添加刀匠'}</DialogTitle>
                <DialogDescription>请填写刀匠信息。</DialogDescription>
              </DialogHeader>

              <Form {...form}>
                <form onSubmit={form.handleSubmit(handleSaveSwordsmith)} className="space-y-4">
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>名称 *</FormLabel>
                        <FormControl>
                          <Input placeholder="请输入刀匠名称" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="image_url"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>图片地址</FormLabel>
                        <FormControl>
                          <Input placeholder="请输入图片地址" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter>
                    <Button type="button" variant="outline" onClick={() => setEditDialogVisible(false)}>
                      取消
                    </Button>
                    <Button type="submit" disabled={saving}>
                      {isEditing ? '更新' : '添加'}
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </DialogContent>
          </Dialog>

          <Dialog open={deleteDialog.open} onOpenChange={(open) => setDeleteDialog(prev => ({ ...prev, open }))}>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>危险操作</DialogTitle>
                <DialogDescription>
                  确定要删除刀匠 "{deleteDialog.swordsmith?.name}" 吗？此操作不可恢复！
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <Button variant="outline" onClick={() => setDeleteDialog({ open: false, swordsmith: null })}>
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
