'use client'

import React, { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { 
  Search, 
  Filter, 
  X, 
  RotateCcw,
  ChevronDown,
  ChevronUp
} from "lucide-react"

export interface FilterConfig {
  key: string
  label: string
  type: 'text' | 'select' | 'multiSelect'
  placeholder?: string
  options?: { value: string; label: string }[]
  defaultValue?: any
}

interface SearchFiltersProps {
  filters: FilterConfig[]
  values: Record<string, any>
  onValuesChange: (values: Record<string, any>) => void
  onSearch: () => void
  onReset: () => void
  loading?: boolean
  showAdvanced?: boolean
  className?: string
}

export function SearchFilters({
  filters,
  values,
  onValuesChange,
  onSearch,
  onReset,
  loading = false,
  showAdvanced = true,
  className
}: SearchFiltersProps) {
  const [showAdvancedFilters, setShowAdvancedFilters] = useState(false)
  
  // 分离基础搜索和高级过滤器
  const basicFilters = filters.slice(0, 1) // 通常第一个是搜索框
  const advancedFilters = filters.slice(1)
  
  const hasActiveFilters = Object.values(values).some(value => 
    value !== undefined && value !== '' && value !== null
  )

  const activeFilterCount = Object.entries(values).filter(([key, value]) => 
    value !== undefined && value !== '' && value !== null
  ).length

  const updateValue = (key: string, value: any) => {
    onValuesChange({
      ...values,
      [key]: value
    })
  }

  const clearValue = (key: string) => {
    const newValues = { ...values }
    delete newValues[key]
    onValuesChange(newValues)
  }

  const renderFilter = (filter: FilterConfig) => {
    const value = values[filter.key]

    if (filter.type === 'text') {
      return (
        <div key={filter.key} className="relative">
          <Input
            placeholder={filter.placeholder || filter.label}
            value={value || ''}
            onChange={(e) => updateValue(filter.key, e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && onSearch()}
            className="pr-8"
          />
          {value && (
            <button
              onClick={() => clearValue(filter.key)}
              className="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
            >
              <X className="h-4 w-4" />
            </button>
          )}
        </div>
      )
    }

    if (filter.type === 'select') {
      return (
        <Select
          key={filter.key}
          value={value || ''}
          onValueChange={(newValue) => updateValue(filter.key, newValue)}
        >
          <SelectTrigger>
            <SelectValue placeholder={filter.placeholder || `选择${filter.label}`} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">全部</SelectItem>
            {filter.options?.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      )
    }

    return null
  }

  return (
    <Card className={className}>
      <CardContent className="p-4">
        {/* 基础搜索行 */}
        <div className="flex items-center space-x-3">
          <div className="flex-1 flex items-center space-x-3">
            {basicFilters.map(renderFilter)}
            
            <Button onClick={onSearch} disabled={loading} className="shrink-0">
              <Search className="h-4 w-4 mr-2" />
              搜索
            </Button>
          </div>

          {/* 高级过滤和重置按钮 */}
          <div className="flex items-center space-x-2">
            {hasActiveFilters && (
              <Badge variant="secondary" className="px-2">
                {activeFilterCount} 个筛选
              </Badge>
            )}

            {showAdvanced && advancedFilters.length > 0 && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowAdvancedFilters(!showAdvancedFilters)}
                className="shrink-0"
              >
                <Filter className="h-4 w-4 mr-1" />
                高级筛选
                {showAdvancedFilters ? (
                  <ChevronUp className="h-4 w-4 ml-1" />
                ) : (
                  <ChevronDown className="h-4 w-4 ml-1" />
                )}
              </Button>
            )}

            {hasActiveFilters && (
              <Button
                variant="ghost"
                size="sm"
                onClick={onReset}
                className="shrink-0 text-gray-500 hover:text-gray-700"
              >
                <RotateCcw className="h-4 w-4 mr-1" />
                重置
              </Button>
            )}
          </div>
        </div>

        {/* 高级过滤器 */}
        {showAdvancedFilters && advancedFilters.length > 0 && (
          <div className="mt-4 pt-4 border-t">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
              {advancedFilters.map(renderFilter)}
            </div>
            
            {/* 活跃过滤器显示 */}
            {hasActiveFilters && (
              <div className="mt-4 pt-4 border-t">
                <div className="text-sm text-muted-foreground mb-2">当前筛选条件：</div>
                <div className="flex flex-wrap gap-2">
                  {Object.entries(values).map(([key, value]) => {
                    if (!value) return null
                    const filter = filters.find(f => f.key === key)
                    if (!filter) return null
                    
                    let displayValue = value
                    if (filter.type === 'select' && filter.options) {
                      const option = filter.options.find(opt => opt.value === value)
                      displayValue = option ? option.label : value
                    }

                    return (
                      <Badge key={key} variant="outline" className="gap-1">
                        <span className="text-xs">{filter.label}:</span>
                        <span>{displayValue}</span>
                        <button
                          onClick={() => clearValue(key)}
                          className="ml-1 hover:text-red-500"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </Badge>
                    )
                  })}
                </div>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  )
}