'use client'

import React from 'react'
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Badge } from "@/components/ui/badge"
import { 
  ChevronDown, 
  Trash2, 
  Download, 
  Upload, 
  FileText, 
  RefreshCw,
  X,
  CheckSquare,
  Square
} from "lucide-react"

export interface BatchAction {
  key: string
  label: string
  icon?: React.ReactNode
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  onClick: (selectedIds: string[]) => void
  disabled?: boolean
}

interface BatchActionsProps {
  selectedIds: string[]
  totalCount: number
  onClearSelection: () => void
  onSelectAll: () => void
  actions: BatchAction[]
  className?: string
}

export function BatchActions({
  selectedIds,
  totalCount,
  onClearSelection,
  onSelectAll,
  actions,
  className
}: BatchActionsProps) {
  const hasSelection = selectedIds.length > 0
  const isAllSelected = selectedIds.length === totalCount

  if (!hasSelection) {
    return (
      <div className={`flex items-center justify-between ${className}`}>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={onSelectAll}
            className="h-8"
          >
            <CheckSquare className="h-4 w-4 mr-1" />
            全选
          </Button>
          <span className="text-sm text-muted-foreground">
            共 {totalCount} 项
          </span>
        </div>
      </div>
    )
  }

  return (
    <div className={`flex items-center justify-between bg-blue-50 dark:bg-blue-950 p-3 rounded-lg ${className}`}>
      <div className="flex items-center space-x-3">
        <Badge variant="secondary" className="bg-blue-100 text-blue-800 dark:bg-blue-800 dark:text-blue-100">
          已选择 {selectedIds.length} 项
        </Badge>
        
        {!isAllSelected && (
          <Button
            variant="ghost"
            size="sm"
            onClick={onSelectAll}
            className="h-7 text-blue-600 hover:text-blue-700"
          >
            <CheckSquare className="h-4 w-4 mr-1" />
            全选
          </Button>
        )}

        <Button
          variant="ghost"
          size="sm"
          onClick={onClearSelection}
          className="h-7 text-gray-500 hover:text-gray-700"
        >
          <X className="h-4 w-4 mr-1" />
          取消选择
        </Button>
      </div>

      <div className="flex items-center space-x-2">
        {/* 主要操作按钮 */}
        {actions.slice(0, 2).map((action) => (
          <Button
            key={action.key}
            variant={action.variant || 'outline'}
            size="sm"
            onClick={() => action.onClick(selectedIds)}
            disabled={action.disabled}
            className="h-8"
          >
            {action.icon}
            <span className="ml-1">{action.label}</span>
          </Button>
        ))}

        {/* 更多操作下拉菜单 */}
        {actions.length > 2 && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-8">
                更多操作
                <ChevronDown className="h-4 w-4 ml-1" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {actions.slice(2).map((action, index) => (
                <div key={action.key}>
                  {index === 0 && actions.length > 2 && <DropdownMenuSeparator />}
                  <DropdownMenuItem
                    onClick={() => action.onClick(selectedIds)}
                    disabled={action.disabled}
                    className={action.variant === 'destructive' ? 'text-red-600' : ''}
                  >
                    {action.icon}
                    <span className="ml-2">{action.label}</span>
                  </DropdownMenuItem>
                </div>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>
    </div>
  )
}

// 常用的批量操作
export const commonBatchActions = {
  delete: (onDelete: (ids: string[]) => void): BatchAction => ({
    key: 'delete',
    label: '批量删除',
    icon: <Trash2 className="h-4 w-4" />,
    variant: 'destructive' as const,
    onClick: onDelete
  }),
  
  export: (onExport: (ids: string[]) => void): BatchAction => ({
    key: 'export',
    label: '导出选中',
    icon: <Download className="h-4 w-4" />,
    variant: 'outline' as const,
    onClick: onExport
  }),

  activate: (onActivate: (ids: string[]) => void): BatchAction => ({
    key: 'activate',
    label: '批量启用',
    icon: <CheckSquare className="h-4 w-4" />,
    variant: 'default' as const,
    onClick: onActivate
  }),

  deactivate: (onDeactivate: (ids: string[]) => void): BatchAction => ({
    key: 'deactivate',
    label: '批量禁用',
    icon: <Square className="h-4 w-4" />,
    variant: 'outline' as const,
    onClick: onDeactivate
  })
}