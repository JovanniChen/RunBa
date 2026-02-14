'use client'

import { useState } from 'react'
import { certificatePublicAPI, Certificate, extractApiData } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { toast } from '@/components/ui/toast'

const titleJpClass =
  'text-[28px] leading-none tracking-[0.16em] text-[#15243b] sm:text-[38px]'
const titleEnClass =
  'mt-3 text-[11px] font-semibold tracking-[0.16em] text-[#2d3d52] sm:text-[14px]'
const fieldJpClass = 'text-[18px] leading-none text-[#15243b] sm:text-[22px]'
const fieldEnClass =
  'mt-1 text-[9px] uppercase tracking-[0.08em] text-[#3f4f65] sm:text-[10px]'
const fieldValueClass =
  'pt-[1px] text-[18px] font-medium leading-[1.05] tracking-[0.01em] text-[#1f2b3d] sm:text-[21px]'

const formatValue = (value?: string | number | null) => {
  if (value === null || value === undefined || value === '') return '-'
  return String(value)
}

type FieldRowProps = {
  jp: string
  en: string
  value?: string | number | null
  noWrapJp?: boolean
}

function FieldRow({ jp, en, value, noWrapJp }: FieldRowProps) {
  return (
    <div className="grid grid-cols-[120px_1fr] items-start gap-4 sm:grid-cols-[152px_1fr] sm:gap-6">
      <div>
        <p className={`${fieldJpClass}${noWrapJp ? ' whitespace-nowrap' : ''}`}>
          {jp}
        </p>
        <p className={fieldEnClass}>{en}</p>
      </div>
      <p className={fieldValueClass}>{formatValue(value)}</p>
    </div>
  )
}

type ImageFieldRowProps = {
  jp: string
  en: string
  imageUrl?: string
  alt?: string
}

type FieldCell =
  | { kind: 'text'; data: FieldRowProps }
  | { kind: 'image'; data: ImageFieldRowProps }

function ImageFieldRow({ jp, en, imageUrl, alt }: ImageFieldRowProps) {
  return (
    <div className="grid grid-cols-[120px_1fr] items-start gap-4 sm:grid-cols-[152px_1fr] sm:gap-6">
      <div>
        <p className={fieldJpClass}>{jp}</p>
        <p className={fieldEnClass}>{en}</p>
      </div>
      {imageUrl ? (
        <div className="flex min-h-[88px] items-start sm:min-h-[96px]">
          <div>
            <img
              src={imageUrl}
              alt={alt || jp}
              className="h-[74px] w-[74px] object-contain sm:h-[84px] sm:w-[84px]"
            />
          </div>
        </div>
      ) : (
        <div className="h-[74px] w-[74px] border border-dashed border-[#b9c0cb] bg-white/40 sm:h-[84px] sm:w-[84px]" />
      )}
    </div>
  )
}

function FieldCellRenderer({ cell }: { cell: FieldCell }) {
  if (cell.kind === 'text') {
    return (
      <FieldRow
        jp={cell.data.jp}
        en={cell.data.en}
        value={cell.data.value}
        noWrapJp={cell.data.noWrapJp}
      />
    )
  }

  return (
    <ImageFieldRow
      jp={cell.data.jp}
      en={cell.data.en}
      imageUrl={cell.data.imageUrl}
      alt={cell.data.alt}
    />
  )
}

function CornerDecor() {
  return (
    <div className="pointer-events-none absolute inset-0">
      <span className="absolute left-4 top-4 h-5 w-5 border-l border-t border-[#d2d2d2] sm:left-7 sm:top-7 sm:h-7 sm:w-7" />
      <span className="absolute right-4 top-4 h-5 w-5 border-r border-t border-[#d2d2d2] sm:right-7 sm:top-7 sm:h-7 sm:w-7" />
      <span className="absolute bottom-4 left-4 h-5 w-5 border-b border-l border-[#d2d2d2] sm:bottom-7 sm:left-7 sm:h-7 sm:w-7" />
      <span className="absolute bottom-4 right-4 h-5 w-5 border-b border-r border-[#d2d2d2] sm:bottom-7 sm:right-7 sm:h-7 sm:w-7" />
    </div>
  )
}

function BrandMark() {
  return (
    <img
      src="https://i.ibb.co/C580CyVP/1500-500.png"
      alt="ソードカスタム"
      className="h-auto w-full max-w-[430px] object-contain"
      loading="lazy"
    />
  )
}

export default function CertificateLookupPage() {
  const [searchNo, setSearchNo] = useState('')
  const [loading, setLoading] = useState(false)
  const [certificate, setCertificate] = useState<Certificate | null>(null)

  const handleSearch = async () => {
    const trimmed = searchNo.trim()
    if (!trimmed) {
      toast.error('証明書番号を入力してください', '入力エラー')
      return
    }

    try {
      setLoading(true)
      const response = await certificatePublicAPI.getCertificateByNo(trimmed)
      const data = extractApiData<Certificate>(response)
      setCertificate(data)
    } catch (err: any) {
      setCertificate(null)
      toast.error(err.message || '証明書情報が見つかりません', '照会失敗')
    } finally {
      setLoading(false)
    }
  }

  const leftFields: FieldRowProps[] = certificate
    ? [
        { jp: '刀の番号', en: 'とうろくばんごう', value: certificate.number },
        { jp: '全長', en: 'ぜんちょう', value: certificate.overall_length },
        { jp: '刀長', en: 'とうちょう', value: certificate.blade_length },
        {
          jp: '捲柄の長',
          en: 'つかのながさ',
          value: certificate.handle_length
        },
        {
          jp: '刀の材質',
          en: 'はざい',
          value: certificate.blade_material
        },
        {
          jp: '刀の厚さ',
          en: 'あつみ',
          value: certificate.blade_thickness
        },
        {
          jp: '切先の長さ',
          en: 'きっさきのながさ',
          value: certificate.kissan_length
        },
        { jp: '刃文', en: 'はもん', value: certificate.hamon },
        { jp: '目釘', en: 'めくぎ', value: certificate.mekugi }
      ]
    : []

  const rightTextFields: FieldRowProps[] = certificate
    ? [
        { jp: '先幅', en: 'さきはば', value: certificate.sakihaba },
        { jp: '元幅', en: 'もとはば', value: certificate.motohaba },
        {
          jp: '鞘の材質',
          en: 'さやざい',
          value: certificate.saya_material
        },
        {
          jp: '鍔の材質',
          en: 'つばざい',
          value: certificate.tsuba_material
        },
        {
          jp: '刀鎺の材質',
          en: 'はばきざい',
          value: certificate.habaki_material
        },
        {
          jp: '上、下緒の材質',
          en: 'いと・さげおざい',
          value: certificate.ito_sageo_material,
          noWrapJp: true
        }
      ]
    : []

  const rightCells: FieldCell[] = certificate
    ? [
        ...rightTextFields.map(
          (field) => ({ kind: 'text', data: field }) as FieldCell
        ),
        {
          kind: 'image',
          data: {
            jp: '鍛刀所',
            en: 'たんとうしょ',
            imageUrl: certificate.forge_image_url,
            alt: certificate.forge_name || '鍛刀所'
          }
        },
        {
          kind: 'image',
          data: {
            jp: '刀匠',
            en: 'とうしょう',
            imageUrl: certificate.swordsmith_image_url,
            alt: certificate.swordsmith_name || '刀匠'
          }
        },
        {
          kind: 'text',
          data: {
            jp: '交付日',
            en: 'こうふび',
            value: certificate.date_completed
          }
        }
      ]
    : []

  const pairedRows = leftFields.map((leftField, index) => ({
    left: { kind: 'text', data: leftField } as FieldCell,
    right: rightCells[index]
  }))

  return (
    <div
      className="min-h-screen bg-[radial-gradient(circle_at_top,_#f8f8f8_0%,_#efefee_55%,_#ececeb_100%)] text-[#1a2435]"
      style={{
        fontFamily: '"Noto Serif JP","Shippori Mincho","Times New Roman",serif'
      }}
    >
      <div className="mx-auto max-w-[1720px] px-3 py-8 sm:px-4 sm:py-10 lg:px-5">
        <div className="space-y-2">
          <p className="text-xs uppercase tracking-[0.42em] text-[#9aa3af]">
            登録証照会
          </p>
          <h1 className="text-2xl font-semibold text-[#1f2937] sm:text-3xl">
            刀剣登録証照会
          </h1>
          <p className="text-sm text-[#667085]">
            証明書番号を入力して、刀剣登録証情報を照会してください。
          </p>
        </div>

        <Card className="mt-6 border border-[#d7d7d7] bg-white/80 p-5 shadow-none sm:p-6">
          <form
            className="flex flex-col gap-4 sm:flex-row sm:items-end"
            onSubmit={(event) => {
              event.preventDefault()
              handleSearch()
            }}
          >
            <div className="flex-1">
              <label className="text-sm font-medium text-[#445064]">
                証明書番号
              </label>
              <Input
                value={searchNo}
                onChange={(event) => setSearchNo(event.target.value)}
                placeholder="証明書番号を入力"
                className="mt-2 border-[#d4d4d4] bg-white"
              />
            </div>
            <Button
              type="submit"
              disabled={loading}
              className="sm:w-32 bg-[#15243b] text-white hover:bg-[#0f1b2c]"
            >
              {loading ? '照会中...' : '照会'}
            </Button>
          </form>
        </Card>

        {certificate && (
          <section className="relative mt-8 rounded-sm border border-[#d1d1d1] bg-[#f6f6f5] px-4 py-7 sm:px-7 sm:py-10 lg:px-8 lg:py-12">
            <CornerDecor />

            <div className="text-center">
              <h2 className={titleJpClass}>刀剣の登録証</h2>
              <p className={titleEnClass}>真正性証明書</p>
            </div>

            <div className="mt-8 grid gap-8 xl:mt-10 xl:grid-cols-[minmax(0,1.5fr)_minmax(250px,0.58fr)] xl:gap-8">
              <div className="space-y-6 sm:space-y-7">
                {pairedRows.map((row, index) => (
                  <div
                    key={`${row.left.kind}-${index}`}
                    className="grid grid-cols-1 gap-6 md:grid-cols-2 md:gap-8"
                  >
                    <FieldCellRenderer cell={row.left} />
                    {row.right ? <FieldCellRenderer cell={row.right} /> : null}
                  </div>
                ))}
              </div>

              <div className="flex min-w-0 flex-col xl:items-center">
                <div className="flex w-full flex-1 items-center justify-center py-7 sm:py-8">
                  <BrandMark />
                </div>

                <div className="mt-4 w-full border-t border-[#d8d8d8] pt-4 text-right">
                  <p className="text-[22px] leading-none text-[#1d2a3d] sm:text-[26px]">
                    登録証番号
                  </p>
                  <p className="mt-1 text-[10px] tracking-[0.04em] text-[#49566b] sm:text-[11px]">
                    とうろくしょうばんごう
                  </p>
                  <p className="mt-3 text-[30px] font-semibold leading-none text-[#1a2537] sm:text-[40px]">
                    {formatValue(certificate.no)}
                  </p>
                </div>
              </div>
            </div>
          </section>
        )}
      </div>
    </div>
  )
}
