'use client'

import React from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Checkbox } from "@/components/ui/checkbox"
import { Switch } from "@/components/ui/switch"
import { Loader2 } from "lucide-react"

export interface FormFieldConfig {
  name: string
  label: string
  type: 'text' | 'email' | 'password' | 'number' | 'textarea' | 'select' | 'checkbox' | 'switch'
  placeholder?: string
  description?: string
  options?: { value: string | number; label: string }[]
  required?: boolean
  disabled?: boolean
  min?: number
  max?: number
  rows?: number
}

interface FormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  fields: FormFieldConfig[]
  schema: z.ZodSchema
  defaultValues?: Record<string, any>
  onSubmit: (data: any) => Promise<void>
  submitLabel?: string
  cancelLabel?: string
  loading?: boolean
  className?: string
}

export function FormDialog({
  open,
  onOpenChange,
  title,
  description,
  fields,
  schema,
  defaultValues = {},
  onSubmit,
  submitLabel = '确定',
  cancelLabel = '取消',
  loading = false,
  className
}: FormDialogProps) {
  const form = useForm({
    resolver: zodResolver(schema as any),
    defaultValues
  })

  const handleSubmit = async (data: any) => {
    try {
      await onSubmit(data)
      form.reset()
      onOpenChange(false)
    } catch (error) {
      console.error('Form submission error:', error)
    }
  }

  const renderField = (field: FormFieldConfig) => {
    const { name, label, type, placeholder, description, options, required, disabled, min, max, rows } = field

    return (
      <FormField
        key={name}
        control={form.control}
        name={name}
        render={({ field: formField }) => (
          <FormItem>
            <FormLabel>
              {label}
              {required && <span className="text-red-500 ml-1">*</span>}
            </FormLabel>
            <FormControl>
              {type === 'text' || type === 'email' || type === 'password' ? (
                <Input
                  type={type}
                  placeholder={placeholder}
                  disabled={disabled || loading}
                  {...formField}
                />
              ) : type === 'number' ? (
                <Input
                  type="number"
                  placeholder={placeholder}
                  disabled={disabled || loading}
                  min={min}
                  max={max}
                  {...formField}
                  onChange={(e) => formField.onChange(Number(e.target.value) || 0)}
                />
              ) : type === 'textarea' ? (
                <Textarea
                  placeholder={placeholder}
                  disabled={disabled || loading}
                  rows={rows}
                  {...formField}
                />
              ) : type === 'select' ? (
                <Select
                  disabled={disabled || loading}
                  value={formField.value?.toString()}
                  onValueChange={formField.onChange}
                >
                  <SelectTrigger>
                    <SelectValue placeholder={placeholder} />
                  </SelectTrigger>
                  <SelectContent>
                    {options?.map((option) => (
                      <SelectItem key={option.value} value={option.value.toString()}>
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : type === 'checkbox' ? (
                <div className="flex items-center space-x-2">
                  <Checkbox
                    checked={formField.value}
                    onCheckedChange={formField.onChange}
                    disabled={disabled || loading}
                  />
                  {placeholder && <span className="text-sm">{placeholder}</span>}
                </div>
              ) : type === 'switch' ? (
                <div className="flex items-center space-x-2">
                  <Switch
                    checked={formField.value}
                    onCheckedChange={formField.onChange}
                    disabled={disabled || loading}
                  />
                  {placeholder && <span className="text-sm">{placeholder}</span>}
                </div>
              ) : null}
            </FormControl>
            {description && <FormDescription>{description}</FormDescription>}
            <FormMessage />
          </FormItem>
        )}
      />
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className={`sm:max-w-md ${className}`}>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description && <DialogDescription>{description}</DialogDescription>}
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <div className="space-y-4">
              {fields.map(renderField)}
            </div>

            <div className="flex justify-end space-x-2 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={loading}
              >
                {cancelLabel}
              </Button>
              <Button type="submit" disabled={loading}>
                {loading && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                {submitLabel}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}