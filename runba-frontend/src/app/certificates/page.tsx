'use client'

import { useEffect, useState } from 'react'
import { ProtectedRoute } from '@/lib/auth'
import { DashboardLayout } from '@/components/dashboard-layout'
import { certificateAPI, Certificate, extractPaginatedData, PaginatedResponse, forgeAPI, Forge, swordsmithAPI, Swordsmith } from '@/lib/api'
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
import { Search, RefreshCw, Plus, Eye, Pencil, Trash2 } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'

const certificateFormSchema = z.object({
  number: z.string().min(1, '请输入证书编号'),
  overall_length: z.string().optional(),
  blade_length: z.string().optional(),
  handle_length: z.string().optional(),
  blade_material: z.string().optional(),
  blade_thickness: z.string().optional(),
  kissan_length: z.string().optional(),
  hamon: z.string().optional(),
  mekugi: z.string().optional(),
  sakihaba: z.string().optional(),
  motohaba: z.string().optional(),
  saya_material: z.string().optional(),
  tsuba_material: z.string().optional(),
  habaki_material: z.string().optional(),
  ito_sageo_material: z.string().optional(),
  forge_id: z.string().min(1, '请选择锻刀所'),
  swordsmith_id: z.string().min(1, '请选择刀匠'),
  date_completed: z.string().optional(),
  no: z.string().optional()
})

type CertificateFormData = z.infer<typeof certificateFormSchema>

const defaultFormValues: CertificateFormData = {
  number: 'SC0001',
  overall_length: '110 cm / 43.3 in',
  blade_length: '76.2 cm / 30 in',
  handle_length: '27 cm / 10.6 in',
  blade_material: 'FOLDED STEEL',
  blade_thickness: '7.5 mm / 0.29 in',
  kissan_length: '5.5 cm / 2.17 in',
  hamon: 'Nature',
  mekugi: '2',
  sakihaba: '3.2 cm / 1.26 in',
  motohaba: '2.2 cm / 0.86 in',
  saya_material: 'HARD WOOD',
  tsuba_material: 'BRASS',
  habaki_material: 'BRASS',
  ito_sageo_material: 'LEATHER+SYNTHETIC SILK',
  forge_id: '',
  swordsmith_id: '',
  date_completed: '乙巳年己丑月吉日',
  no: 'YFSC2026124578'
}

export default function CertificatesPage() {
  const [certificateList, setCertificateList] = useState<Certificate[]>([])
  const [loading, setLoading] = useState(false)
  const [forgeOptions, setForgeOptions] = useState<Forge[]>([])
  const [swordsmithOptions, setSwordsmithOptions] = useState<Swordsmith[]>([])
  const [optionsLoading, setOptionsLoading] = useState(false)

  // 搜索表单状态
  const [searchForm, setSearchForm] = useState({
    keyword: ''
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
  const [currentCertificate, setCurrentCertificate] = useState<Certificate | null>(null)

  // 编辑对话框状态
  const [editDialogVisible, setEditDialogVisible] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [saving, setSaving] = useState(false)

  // 删除确认对话框
  const [deleteDialog, setDeleteDialog] = useState<{
    open: boolean
    certificate: Certificate | null
  }>({ open: false, certificate: null })

  const form = useForm<CertificateFormData>({
    resolver: zodResolver(certificateFormSchema as any),
    defaultValues: defaultFormValues
  })

  useEffect(() => {
    loadCertificates()
  }, [pagination.page, pagination.pageSize])

  useEffect(() => {
    loadReferenceOptions()
  }, [])

  // 加载证书列表
  const loadCertificates = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = {
        page: pagination.page,
        page_size: pagination.pageSize,
        keyword: searchForm.keyword || undefined
      }

      const response = await certificateAPI.getCertificates(params) as unknown as PaginatedResponse<Certificate>
      const { list, pagination: paginationData } = extractPaginatedData(response)

      setCertificateList(list)
      setPagination(prev => ({
        ...prev,
        total: paginationData.total,
        totalPages: paginationData.totalPages
      }))
    } catch (err: any) {
      toast.error('加载证书列表失败: ' + (err.message || ''), '加载失败')
    } finally {
      setLoading(false)
    }
  }

  // 加载锻刀所与刀匠选项
  const loadReferenceOptions = async () => {
    setOptionsLoading(true)
    try {
      const [forgeResponse, swordsmithResponse] = await Promise.all([
        forgeAPI.getForges({ page: 1, page_size: 100 }),
        swordsmithAPI.getSwordsmiths({ page: 1, page_size: 100 })
      ])

      const { list: forgeList } = extractPaginatedData(forgeResponse as unknown as PaginatedResponse<Forge>)
      const { list: swordsmithList } = extractPaginatedData(swordsmithResponse as unknown as PaginatedResponse<Swordsmith>)

      setForgeOptions(forgeList)
      setSwordsmithOptions(swordsmithList)
    } catch (err: any) {
      toast.error('加载锻刀所或刀匠失败: ' + (err.message || ''), '加载失败')
    } finally {
      setOptionsLoading(false)
    }
  }

  // 搜索
  const handleSearch = () => {
    setPagination(prev => ({ ...prev, page: 1 }))
    loadCertificates()
  }

  // 重置
  const handleReset = () => {
    setSearchForm({ keyword: '' })
    setPagination(prev => ({ ...prev, page: 1 }))
    setTimeout(() => {
      loadCertificates()
    }, 0)
  }

  // 刷新数据
  const refreshData = () => {
    loadCertificates()
    toast.info('数据已刷新', '刷新成功')
  }

  // 查看证书详情
  const handleView = (certificate: Certificate) => {
    setCurrentCertificate(certificate)
    setDialogVisible(true)
  }

  // 添加证书
  const handleAdd = () => {
    setIsEditing(false)
    setCurrentCertificate(null)
    form.reset(defaultFormValues)
    setEditDialogVisible(true)
  }

  // 编辑证书
  const handleEdit = (certificate: Certificate) => {
    setIsEditing(true)
    setCurrentCertificate(certificate)
    form.reset({
      number: certificate.number || '',
      overall_length: certificate.overall_length || '',
      blade_length: certificate.blade_length || '',
      handle_length: certificate.handle_length || '',
      blade_material: certificate.blade_material || '',
      blade_thickness: certificate.blade_thickness || '',
      kissan_length: certificate.kissan_length || '',
      hamon: certificate.hamon || '',
      mekugi: certificate.mekugi !== undefined && certificate.mekugi !== null ? String(certificate.mekugi) : '',
      sakihaba: certificate.sakihaba || '',
      motohaba: certificate.motohaba || '',
      saya_material: certificate.saya_material || '',
      tsuba_material: certificate.tsuba_material || '',
      habaki_material: certificate.habaki_material || '',
      ito_sageo_material: certificate.ito_sageo_material || '',
      forge_id: certificate.forge_id ? String(certificate.forge_id) : '',
      swordsmith_id: certificate.swordsmith_id ? String(certificate.swordsmith_id) : '',
      date_completed: certificate.date_completed || '',
      no: certificate.no || ''
    })
    setEditDialogVisible(true)
  }

  const buildSubmitData = (data: CertificateFormData) => {
    const payload: Record<string, any> = {}
    const stringFields: Array<keyof CertificateFormData> = [
      'number',
      'overall_length',
      'blade_length',
      'handle_length',
      'blade_material',
      'blade_thickness',
      'kissan_length',
      'hamon',
      'sakihaba',
      'motohaba',
      'saya_material',
      'tsuba_material',
      'habaki_material',
      'ito_sageo_material',
      'date_completed',
      'no'
    ]

    stringFields.forEach((field) => {
      const value = data[field]
      if (typeof value === 'string' && value.trim() !== '') {
        payload[field] = value.trim()
      }
    })

    if (data.mekugi && data.mekugi.trim() !== '') {
      const parsed = Number(data.mekugi)
      if (Number.isNaN(parsed)) {
        throw new Error('目钉数必须是数字')
      }
      payload.mekugi = parsed
    }

    if (data.forge_id && data.forge_id.trim() !== '') {
      const parsed = Number(data.forge_id)
      if (Number.isNaN(parsed) || parsed <= 0) {
        throw new Error('请选择有效的锻刀所')
      }
      payload.forge_id = parsed
    }

    if (data.swordsmith_id && data.swordsmith_id.trim() !== '') {
      const parsed = Number(data.swordsmith_id)
      if (Number.isNaN(parsed) || parsed <= 0) {
        throw new Error('请选择有效的刀匠')
      }
      payload.swordsmith_id = parsed
    }

    return payload
  }

  // 保存证书
  const handleSaveCertificate = async (data: CertificateFormData) => {
    try {
      setSaving(true)
      const submitData = buildSubmitData(data)

      if (!submitData.number) {
        toast.error('请输入证书编号', '参数缺失')
        return
      }
      if (!submitData.forge_id) {
        toast.error('请选择锻刀所', '参数缺失')
        return
      }
      if (!submitData.swordsmith_id) {
        toast.error('请选择刀匠', '参数缺失')
        return
      }

      if (isEditing) {
        if (!currentCertificate) {
          toast.error('未选择要编辑的证书', '更新失败')
          return
        }
        await certificateAPI.updateCertificate(currentCertificate.id, submitData)
        toast.success('证书已更新', '更新成功')
      } else {
        await certificateAPI.createCertificate(submitData)
        toast.success('证书已添加', '添加成功')
      }

      setEditDialogVisible(false)
      loadCertificates()
    } catch (err: any) {
      toast.error(`${isEditing ? '更新' : '添加'}失败: ` + (err.message || ''), `${isEditing ? '更新' : '添加'}失败`)
    } finally {
      setSaving(false)
    }
  }

  // 删除证书
  const handleDelete = (certificate: Certificate) => {
    setDeleteDialog({ open: true, certificate })
  }

  const confirmDelete = async () => {
    if (!deleteDialog.certificate) return

    try {
      await certificateAPI.deleteCertificate(deleteDialog.certificate.id)
      toast.success('证书已删除', '删除成功')
      loadCertificates()
      setDeleteDialog({ open: false, certificate: null })
    } catch (err: any) {
      toast.error('删除失败: ' + (err.message || ''), '删除失败')
    }
  }

  // 分页处理
  const handleSizeChange = (size: number) => {
    setPagination(prev => ({ ...prev, pageSize: size, page: 1 }))
  }

  const handleCurrentChange = (page: number) => {
    setPagination(prev => ({ ...prev, page }))
  }

  // 格式化日期
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

  const formatValue = (value?: string | number | null) => {
    if (value === null || value === undefined || value === '') return '-'
    return String(value)
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
                    placeholder="搜索证书编号、序号、锻刀所、刀匠"
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

          {/* 证书列表 */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>证书列表</CardTitle>
              <div className="flex space-x-2">
                <Button onClick={handleAdd}>
                  <Plus className="h-4 w-4 mr-2" />
                  添加证书
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
              ) : certificateList.length === 0 ? (
                <div className="text-center py-8">
                  <p className="text-muted-foreground">暂无证书数据</p>
                </div>
              ) : (
                <>
                  <div className="border rounded-md">
                    <div className="max-h-[60vh] overflow-auto">
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead className="w-16 text-center">ID</TableHead>
                            <TableHead className="min-w-40 text-center">证书编号</TableHead>
                            <TableHead className="min-w-40 text-center">证书序号</TableHead>
                            <TableHead className="min-w-40 text-center">创建时间</TableHead>
                            <TableHead className="w-32 text-center">操作</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {certificateList.map((certificate) => (
                            <TableRow key={certificate.id}>
                              <TableCell className="font-medium text-center">{certificate.id}</TableCell>
                              <TableCell className="text-center">{formatValue(certificate.number)}</TableCell>
                              <TableCell className="text-center">{formatValue(certificate.no)}</TableCell>
                              <TableCell className="text-center">{formatDate(certificate.created_at)}</TableCell>
                              <TableCell className="text-center">
                                <div className="flex items-center justify-center space-x-1">
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handleView(certificate)}
                                    className="h-8 w-8 p-0"
                                  >
                                    <Eye className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handleEdit(certificate)}
                                    className="h-8 w-8 p-0"
                                  >
                                    <Pencil className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handleDelete(certificate)}
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

          {/* 证书详情对话框 */}
          <Dialog open={dialogVisible} onOpenChange={setDialogVisible}>
            <DialogContent className="max-w-3xl">
              <DialogHeader>
                <DialogTitle>证书详情</DialogTitle>
              </DialogHeader>

              {currentCertificate && (
                <div className="space-y-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">ID:</span>
                      <p className="font-medium">{currentCertificate.id}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">证书编号:</span>
                      <p className="font-medium">{formatValue(currentCertificate.number)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">全长:</span>
                      <p className="font-medium">{formatValue(currentCertificate.overall_length)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">刀身长:</span>
                      <p className="font-medium">{formatValue(currentCertificate.blade_length)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">柄长:</span>
                      <p className="font-medium">{formatValue(currentCertificate.handle_length)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">刀身材质:</span>
                      <p className="font-medium">{formatValue(currentCertificate.blade_material)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">刀身厚:</span>
                      <p className="font-medium">{formatValue(currentCertificate.blade_thickness)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">切先长度:</span>
                      <p className="font-medium">{formatValue(currentCertificate.kissan_length)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">刃纹:</span>
                      <p className="font-medium">{formatValue(currentCertificate.hamon)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">目钉数:</span>
                      <p className="font-medium">{formatValue(currentCertificate.mekugi)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">先幅:</span>
                      <p className="font-medium">{formatValue(currentCertificate.sakihaba)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">元幅:</span>
                      <p className="font-medium">{formatValue(currentCertificate.motohaba)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">鞘材质:</span>
                      <p className="font-medium">{formatValue(currentCertificate.saya_material)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">鍔材质:</span>
                      <p className="font-medium">{formatValue(currentCertificate.tsuba_material)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">鎺材质:</span>
                      <p className="font-medium">{formatValue(currentCertificate.habaki_material)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">柄卷/下绪材质:</span>
                      <p className="font-medium">{formatValue(currentCertificate.ito_sageo_material)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">锻刀所:</span>
                      <p className="font-medium">{formatValue(currentCertificate.forge_name || currentCertificate.forge_id)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">刀匠:</span>
                      <p className="font-medium">{formatValue(currentCertificate.swordsmith_name || currentCertificate.swordsmith_id)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">完成日期:</span>
                      <p className="font-medium">{formatValue(currentCertificate.date_completed)}</p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">证书序号:</span>
                      <p className="font-medium">{formatValue(currentCertificate.no)}</p>
                    </div>
                    <div className="col-span-1 sm:col-span-2">
                      <span className="text-muted-foreground">创建时间:</span>
                      <p className="font-medium">{formatDate(currentCertificate.created_at)}</p>
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

          {/* 添加/编辑证书对话框 */}
          <Dialog open={editDialogVisible} onOpenChange={setEditDialogVisible}>
            <DialogContent className="max-w-4xl">
              <DialogHeader>
                <DialogTitle>{isEditing ? '编辑证书' : '添加证书'}</DialogTitle>
                <DialogDescription>
                  请填写证书信息，带 * 为必填项。
                </DialogDescription>
              </DialogHeader>

              <Form {...form}>
                <form onSubmit={form.handleSubmit(handleSaveCertificate)} className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField
                      control={form.control}
                      name="number"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>证书编号 *</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入证书编号" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="overall_length"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>全长</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入全长" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="blade_length"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>刀身长</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入刀身长" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="handle_length"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>柄长</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入柄长" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="blade_material"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>刀身材质</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入刀身材质" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="blade_thickness"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>刀身厚</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入刀身厚" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="kissan_length"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>切先长度</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入切先长度" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="hamon"
                      render={({ field }) => (
                        <FormItem className="md:col-span-2">
                          <FormLabel>刃纹</FormLabel>
                          <FormControl>
                            <Textarea placeholder="请输入刃纹" rows={3} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="mekugi"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>目钉数</FormLabel>
                          <FormControl>
                            <Input type="number" min={0} step={1} placeholder="请输入目钉数" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="sakihaba"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>先幅</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入先幅" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="motohaba"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>元幅</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入元幅" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="saya_material"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>鞘材质</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入鞘材质" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="tsuba_material"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>鍔材质</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入鍔材质" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="habaki_material"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>鎺材质</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入鎺材质" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="ito_sageo_material"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>柄卷/下绪材质</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入柄卷/下绪材质" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="forge_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>锻刀所 *</FormLabel>
                          <Select value={field.value} onValueChange={field.onChange}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder={optionsLoading ? '加载中...' : '请选择锻刀所'} />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {forgeOptions.length === 0 ? (
                                <SelectItem value="__empty__" disabled>
                                  {optionsLoading ? '加载中...' : '暂无锻刀所'}
                                </SelectItem>
                              ) : (
                                forgeOptions.map((forge) => (
                                  <SelectItem key={forge.id} value={String(forge.id)}>
                                    {forge.name}
                                  </SelectItem>
                                ))
                              )}
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="swordsmith_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>刀匠 *</FormLabel>
                          <Select value={field.value} onValueChange={field.onChange}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder={optionsLoading ? '加载中...' : '请选择刀匠'} />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {swordsmithOptions.length === 0 ? (
                                <SelectItem value="__empty__" disabled>
                                  {optionsLoading ? '加载中...' : '暂无刀匠'}
                                </SelectItem>
                              ) : (
                                swordsmithOptions.map((swordsmith) => (
                                  <SelectItem key={swordsmith.id} value={String(swordsmith.id)}>
                                    {swordsmith.name}
                                  </SelectItem>
                                ))
                              )}
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="date_completed"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>完成日期</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入交付日期" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="no"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>证书序号</FormLabel>
                          <FormControl>
                            <Input placeholder="请输入证书序号" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <DialogFooter>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => setEditDialogVisible(false)}
                    >
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

          {/* 删除确认对话框 */}
          <Dialog open={deleteDialog.open} onOpenChange={(open) => setDeleteDialog(prev => ({ ...prev, open }))}>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>危险操作</DialogTitle>
                <DialogDescription>
                  确定要删除证书 "{deleteDialog.certificate?.number || deleteDialog.certificate?.id}" 吗？此操作不可恢复！
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <Button
                  variant="outline"
                  onClick={() => setDeleteDialog({ open: false, certificate: null })}
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
