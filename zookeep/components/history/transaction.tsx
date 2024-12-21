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
    TableCell,  // เพิ่ม TableCell
    TableHead,
    TableHeader,
    TableRow,   // เพิ่ม TableRow
    TableFooter,
  } from "@/components/ui/table"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { useTranslation } from "react-i18next"
import { useMemo, useState } from "react"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface TransactionItem {
  ID: string
  Status: string
  CreatedAt: string
  GameProvide: string
  BeforeBalance: number
  BetAmount: number
  TransactionAmount: number
  Balance: number
}

interface TransactionTableProps {
  statement: TransactionItem[]
}

export function TransactionTable({ statement }: TransactionTableProps) {
  const { t } = useTranslation()
  const [selectedGameProvider, setSelectedGameProvider] = useState<string>('all')

  const gameProviders = useMemo(() => {
    const providers = new Set(statement.map(item => item.GameProvide))
    return ['all', ...Array.from(providers)]
  }, [statement])

  const filteredData = useMemo(() => {
    if (selectedGameProvider === 'all') return statement
    return statement.filter(item => item.GameProvide === selectedGameProvider)
  }, [statement, selectedGameProvider])

  const columns = useMemo<ColumnDef<TransactionItem>[]>(() => [
    {
      accessorKey: "CreatedAt",
      header: t('transaction.transactionDate'),
      cell: ({ row }) => {
        const date = new Date(row.getValue("CreatedAt"))
        return (
          <div className="text-gray-600">
            {date.toLocaleDateString()} {date.toLocaleTimeString()}
          </div>
        )
      }
    },
    {
      accessorKey: "GameProvide",
      header: t('transaction.gameprovide'),
      cell: ({ row }) => (
        <div className="text-gray-700">
          {row.getValue("GameProvide")}
        </div>
      )
    },
    {
      accessorKey: "BeforeBalance",
      header: t('transaction.beforeBalance'),
      cell: ({ row }) => (
        <div className="text-gray-700">
          {(row.getValue("BeforeBalance") as number).toLocaleString()}
        </div>
      )
    },
    {
      accessorKey: "BetAmount",
      header: t('transaction.betamount'),
      cell: ({ row }) => (
        <div className="text-gray-700 font-medium">
          {(row.getValue("BetAmount") as number).toLocaleString()}
        </div>
      )
    },
    {
      accessorKey: "TransactionAmount",
      header: t('transaction.transactionAmount'),
      cell: ({ row }) => {
        const amount = row.getValue("TransactionAmount") as number
        return (
          <div className={cn(
            "font-medium",
            amount >= 0 ? "text-green-600" : "text-red-600"
          )}>
            {amount >= 0 ? `+${amount.toLocaleString()}` : amount.toLocaleString()}
          </div>
        )
      }
    },
    {
      accessorKey: "Balance",
      header: t('transaction.balance'),
      cell: ({ row }) => (
        <div className="text-gray-700 font-medium">
          {(row.getValue("Balance") as number).toLocaleString()}
        </div>
      )
    },
    {
      accessorKey: "Turnover",
      header: t('transaction.turover'),
      cell: ({ row }) => (
        <div className="text-gray-700">
          {(row.getValue("BetAmount") as number).toLocaleString()}
        </div>
      )
    },
    {
      accessorKey: "Status",
      header: t('transaction.Status'),
      cell: ({ row }) => {
        const status = row.getValue("Status") as string
        const result = row.getValue("TransactionAmount") as number
        return (
          <div className={cn(
            "px-3 py-1 rounded-full text-xs font-medium w-fit",
            status === '100' 
              ? 'bg-yellow-100 text-yellow-700' // เปลี่ยนสีสถานะ Bet เป็นสีเหลือง
              : result >= 0 
                ? 'bg-green-100 text-green-700' // ผลลัพธ์เป็นบ��ก = สีเขียว
                : 'bg-red-100 text-red-700' // ผลลัพธ์เป็นลบ = สีแดง
          )}>
            {status === "100" ? t('transaction.bet') : t('transaction.result')}
          </div>
        )
      }
    },
  ], [t])

  const table = useReactTable({
    data: filteredData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  })

  const currentPageSums = useMemo(() => {
    const pageRows = table.getRowModel().rows
    return {
        betAmount: pageRows.reduce((sum, row) => sum + Number(row.getValue("BetAmount")), 0),
        transactionAmount: pageRows.reduce((sum, row) => sum + Number(row.getValue("TransactionAmount")), 0),
        turnover: pageRows.reduce((sum, row) => sum + Number(row.getValue("BetAmount")), 0),
       }
  }, [table.getRowModel().rows])

  const totalSums = useMemo(() => {
    return {
      betAmount: filteredData.reduce((sum, row) => sum + Number(row.BetAmount), 0),
      transactionAmount: filteredData.reduce((sum, row) => sum + Number(row.TransactionAmount), 0),
      turnover: filteredData.reduce((sum, row) => sum + Number(row.BetAmount), 0),
    }
  }, [filteredData])

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-gray-700">
            {t('transaction.filterByProvider')}:
          </span>
          <Select value={selectedGameProvider} onValueChange={setSelectedGameProvider}>
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder={t('transaction.selectProvider')} />
            </SelectTrigger>
            <SelectContent>
              {gameProviders.map((provider) => (
                <SelectItem key={provider} value={provider}>
                  {provider === 'all' ? t('transaction.allProviders') : provider}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
      <div className="rounded-md border border-gray-200">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id} className="bg-gray-50">
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} className="font-semibold text-gray-700">
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
                className="hover:bg-gray-50 transition-colors"
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
          <TableFooter>
            {/* Current Page Total Row */}
            <TableRow className="border-t-2 border-gray-200">
              <TableCell colSpan={3} className="font-medium text-gray-700">
                {t('transaction.currentPageTotal')}
              </TableCell>
              <TableCell className="font-medium text-gray-900">
                {currentPageSums.betAmount.toLocaleString()}
              </TableCell>
              <TableCell className={cn(
                "font-medium",
                currentPageSums.transactionAmount >= 0 
                  ? "text-green-700" 
                  : "text-red-700"
              )}>
                {currentPageSums.transactionAmount >= 0 ? "+" : ""}
                {currentPageSums.transactionAmount.toLocaleString()}
              </TableCell>
              <TableCell></TableCell>
              <TableCell className="font-medium text-gray-900">
                {currentPageSums.turnover.toLocaleString()}
              </TableCell>
              <TableCell></TableCell>
            </TableRow>  

            {/* Grand Total Row */}
            <TableRow className="bg-gray-50 border-t-2 border-gray-300">
              <TableCell colSpan={3} className="font-semibold text-gray-800">
                {t('transaction.grandTotal')}
              </TableCell>
              <TableCell className="font-semibold text-gray-900">
                {totalSums.betAmount.toLocaleString()}
              </TableCell>
              <TableCell className={cn(
                "font-semibold",
                totalSums.transactionAmount >= 0 
                  ? "text-green-700" 
                  : "text-red-700"
              )}>
                {totalSums.transactionAmount >= 0 ? "+" : ""}
                {totalSums.transactionAmount.toLocaleString()}
              </TableCell>
              <TableCell></TableCell>
              <TableCell className="font-semibold text-gray-900">
                {totalSums.turnover.toLocaleString()}
              </TableCell>
              <TableCell></TableCell>
            </TableRow>
          </TableFooter>
        </Table>
      </div>
      
      {/* Pagination Controls */}
      <div className="flex items-center justify-end space-x-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
          className="hover:bg-gray-100"
        >
          {t('common.previous')}
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
          className="hover:bg-gray-100"
        >
          {t('common.next')}
        </Button>
      </div>
    </div>
  )
}