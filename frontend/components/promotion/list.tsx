"use client"

import React, { useState, useEffect, useMemo } from 'react'
import { useRouter } from 'next/navigation'
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from '@tanstack/react-table'
import { Button } from "@/components/ui/button"
import { PlusIcon } from "@radix-ui/react-icons"
import PromotionLayout from './layout'

interface Promotion {
  id: string
  name: string
  percentDiscount: number
  maxDiscount: number
  usageLimit: string
  specificTime: string
  paymentMethod: string
  minSpend: number
  maxSpend: number
  termsAndConditions: string
}

const PromotionList: React.FC = () => {
  const [promotions, setPromotions] = useState<Promotion[]>([])
  const router = useRouter()

  useEffect(() => {
    // Fetch promotions data from API
    // For now, we'll use mock data
    const mockData: Promotion[] = [
      {
        id: '1',
        name: 'สมัครใหม่',
        percentDiscount: 50,
        maxDiscount: 5000,
        usageLimit: '1 ครั้ง',
        specificTime: 'ทั้งหมด',
        paymentMethod: '-',
        minSpend: 1500,
        maxSpend: 30000,
        termsAndConditions: '(ทุน+โบนัส)*1500%',
      },
      // ... add more mock data based on your image
    ]
    setPromotions(mockData)
  }, [])

  const columnHelper = createColumnHelper<Promotion>()

  const columns = useMemo(() => [
    columnHelper.accessor('name', {
      header: 'โปรโมชั่น',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('percentDiscount', {
      header: 'หรือรอลด',
      cell: info => `${info.getValue()}%`,
    }),
    columnHelper.accessor('maxDiscount', {
      header: 'หรือรอลดสูงสุด',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('usageLimit', {
      header: 'รับได้สูงสุด',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('specificTime', {
      header: 'เฉพาะเวลา',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('paymentMethod', {
      header: 'ชำระเงิน',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('minSpend', {
      header: 'ยอดเงินขั้นต่ำ',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('maxSpend', {
      header: 'ถอนได้สูงสุด',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('termsAndConditions', {
      header: 'เงื่อนไขการถอนยอดโบนัส',
      cell: info => info.getValue(),
    }),
  ], [])

  const table = useReactTable({
    data: promotions,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  const handleAddPromotion = () => {
    router.push('/promotions/add')
  }

  const handleEditPromotion = (id: string) => {
    router.push(`/promotions/edit/${id}`)
  }

  return (
    <PromotionLayout>
      <div className="w-full">
        <div className="flex items-center py-4">
          <Button onClick={handleAddPromotion}>
            <PlusIcon className="mr-2 h-4 w-4" /> เพิ่มโปรโมชั่น
          </Button>
        </div>
        <div className="rounded-md border">
          <table className="w-full">
            <thead>
              {table.getHeaderGroups().map(headerGroup => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map(header => (
                    <th key={header.id} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody>
              {table.getRowModel().rows.map(row => (
                <tr 
                  key={row.id} 
                  onClick={() => handleEditPromotion(row.original.id)}
                  className="cursor-pointer hover:bg-gray-100"
                >
                  {row.getVisibleCells().map(cell => (
                    <td key={cell.id} className="px-6 py-4 whitespace-nowrap">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </PromotionLayout>
  )
}

export default PromotionList
