"use client"

import {
    ColumnDef,
    flexRender,
    getCoreRowModel,
    useReactTable,
    getPaginationRowModel,
  } from "@tanstack/react-table"

  import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
    TableFooter,  // แก้จาก TableFoot เป็น TableFooter
  } from "@/components/ui/table"


import { Button } from "@/components/ui/button"
import { BellIcon } from "@radix-ui/react-icons"
import { cn } from "@/lib/utils"
import { useTranslation } from "@/app/i18n/client"
import { useMemo, useCallback } from "react" // เพิ่ม useCallback


interface PromotionLog {
  ID: number
  CreatedOn: string
  CreatedAt: string
  promotioncode: string
  promotionname: string
  Promoamount: number
  UserID: number
  StatementID: number
  WalletID: number
  Transactionamount: number
  Beforebalance: number
  Proamount: number
  AddOnamount: number
  Balance: number
  Status: number
  Example: string
}

interface HistoryTableProps {
  lng:string
  history: PromotionLog[]
 
}

export function HistoryPromotion({ lng,history }: HistoryTableProps) {
  const { t } = useTranslation(lng,"translation",undefined)
//console.log(history)
  // ฟังก์ชันสำหรับจัดรูปแบบตัวเลข
  const formatNumber = (num: number) => {
    return new Intl.NumberFormat(lng=='th'?'th-TH':'en-EN', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(num)
  }

  const columns = useMemo<ColumnDef<PromotionLog>[]>(() => [
    {
      accessorKey: "CreatedAt",
      header: t('promotionlog.createdAt'),
      cell: ({ row }) => {
        const date = new Date(row.original.CreatedAt)
        return (
          <span className="text-gray-600">
            {date.toLocaleDateString()} {date.toLocaleTimeString()}
          </span>
        )
      }
    },
    {
      accessorKey: "promotioncode",
      header: t('promotionlog.promotioncode'),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.promotioncode}</span>
      )
    },
    {
      accessorKey: "promotionname",
      header: t('promotionlog.promotionname'),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.promotionname}</span>
      )
    },
    {
        accessorKey: "Transactionamount",
        header: t('promotionlog.transactionamount'),
        cell: ({ row }) => (
          <span className="text-blue-600 font-medium">
            {formatNumber(row.original.Transactionamount)}
          </span>
        )
      },
    {
      accessorKey: "Proamount",
      header: t('promotionlog.promoamount'),
      cell: ({ row }) => (
        <span className="text-blue-600 font-medium">
          {formatNumber(row.original.Proamount)}
        </span>
      )
    },
    {
        accessorKey: "Balance",
        header: t('promotionlog.balance'),
        cell: ({ row }) => (
          <span className="text-blue-600 font-medium">
            {formatNumber(row.original.Balance)}
          </span>
        )
      },
    // ... เพิ่มคอลัมน์อื่น ๆ ตามต้องการ ...
  ], [t])

  const table = useReactTable({
    data: history,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  })
  const calculateTotals = useCallback(() => {
    return history.reduce((acc, item) => {
      return {
        deposit: acc.deposit + (item.Transactionamount > 0 ? Number(item.Transactionamount) : 0),
        withdraw: acc.withdraw + (item.Transactionamount < 0 ? Math.abs(item.Transactionamount) : 0),
        proAmount: acc.proAmount + Number(item.Proamount),
        balance: (acc.balance) + Number(item.Balance)
      }
    }, { deposit: 0, withdraw: 0, proAmount: 0 ,balance: 0})
  }, [history])
  const calculatePageTotals = useCallback(() => {
    return table.getRowModel().rows.reduce((acc, row) => {
      const item = row.original
      return {
        deposit: (acc.deposit + Number(item.Transactionamount) > 0 ? Number(item.Transactionamount) : 0),
        withdraw: (acc.withdraw + item.Transactionamount < 0 ? Math.abs(item.Transactionamount) : 0),
        proAmount: (acc.proAmount + Number(item.Proamount)),
        balance: (acc.balance) + Number(item.Balance)
      }
    }, { deposit: 0, withdraw: 0, proAmount: 0 ,balance:0 })
  }, [table])

  const pageTotals = useMemo(() => calculatePageTotals(), [calculatePageTotals])
  const showPageTotals = table.getPageCount() > 1

  const totals = useMemo(() => calculateTotals(), [calculateTotals])
  return (
    <div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} className="font-semibold">
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  className="hover:bg-gray-50"
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center text-gray-500"
                >
                  {t('common.noResults')}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
          {showPageTotals && (
            <TableFooter className="bg-gray-100 border-t">
              <TableRow>
                <TableCell colSpan={3} className="text-right text-sm">
                  {t('common.pageTotals')}:
                </TableCell>
                <TableCell className="text-green-600">
                  {formatNumber(pageTotals.deposit)}
                </TableCell>
                {/* <TableCell className="text-red-600">
                  {formatNumber(pageTotals.withdraw)}
                </TableCell> */}
                <TableCell>
                  <span className={cn(
                    pageTotals.proAmount > 0 ? "text-green-600" : 
                    pageTotals.proAmount < 0 ? "text-red-600" : "text-gray-600"
                  )}>
                    {formatNumber(pageTotals.proAmount)}
                  </span>
                </TableCell>
                <TableCell className="text-green-600">
                  {formatNumber(pageTotals.balance)}
                </TableCell>
                <TableCell></TableCell>
              </TableRow>
            </TableFooter>
          )}
          <TableFooter className="bg-gray-50 font-semibold">
          <TableRow>
            <TableCell colSpan={3} className="text-right">
              {t('common.total')}:
            </TableCell>
            <TableCell className="text-green-600">
              {formatNumber(totals.deposit)}
            </TableCell>
            {/* <TableCell className="text-red-600">
              {formatNumber(totals.withdraw)}
            </TableCell> */}
            <TableCell>
              <span className={cn(
                totals.proAmount > 0 ? "text-green-600" : 
                totals.proAmount < 0 ? "text-red-600" : "text-gray-600"
              )}>
                {formatNumber(totals.proAmount)}
              </span>
            </TableCell>
            <TableCell className="text-green-600">
              {formatNumber(totals.balance)}
            </TableCell>
            <TableCell></TableCell>
          </TableRow>
        </TableFooter>
        </Table>
      </div>
      
      <div className="flex items-center justify-end space-x-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          {t('common.previous')}
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          {t('common.next')}
        </Button>
      </div>
    </div>
  )
}