"use client"
 

import React, { useState, useEffect, useMemo } from 'react'
import { useRouter } from 'next/navigation'
 
import {
  CaretSortIcon,
  ChevronDownIcon,
  DotsHorizontalIcon,
} from "@radix-ui/react-icons"

import { PlusIcon } from "@radix-ui/react-icons"

 

import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import EditPromotionPanel from './EditPromotionPanel'
import { GetPromotion } from '@/actions'
import { useTranslation } from '@/app/i18n/client'; 

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
  include: string
  exclude: string
  startDate: string
  endDate: string
  example: string
}
export interface GroupedDatabase {
  // Define the properties of GroupedDatabase here
  id: string;
  name: string;
  prefix: string;
  // Add other properties as needed
}
interface DataTableProps<TData> {
  columns: ColumnDef<TData, any>[]
  data: TData[]
}

export default function PromotionListDataTable({
  prefix,
  data,
  lng,

}: {prefix:string, data:DataTableProps<GroupedDatabase>, lng:string}) {
  const [promotions, setPromotions] = useState<Promotion[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = useState({})
  const [editingPromotion, setEditingPromotion] = useState<number | null>(null);
  const [isAddingPromotion, setIsAddingPromotion] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const router = useRouter()

  const {t} = useTranslation(lng,'translation',{keyPrefix:'promotion'})

  useEffect(() => {
    const fetchPromotions = async () => {
      if (!prefix) {
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const fetchedPromotions = await GetPromotion(prefix);
        setPromotions(fetchedPromotions.Data);
      } catch (error) {
        console.error('Error fetching promotions:', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchPromotions();
  }, [prefix, refreshTrigger])

  const columnHelper = createColumnHelper<Promotion>()

  const columns = useMemo(() => [
    columnHelper.accessor('ID', {
      header: t('ID'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('name', {
      header: t('name'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('percentDiscount', {
      header: t('percentDiscount'),
      cell: info => `${info.getValue()}%`,
    }),
    columnHelper.accessor('maxDiscount', {
      header: t('maxDiscount'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('usageLimit', {
      header: t('usageLimit'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('specificTime', {
      header: t('specificTime'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('include', {
      header: t('include'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('exclude', {
      header: t('exclude'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('startDate', {
      header: t('startDate'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('endDate', {
      header: t('endDate'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('paymentMethod', {
    //   header: t('paymentMethod'),
    //   cell: info => info.getValue(),
    // }),
    columnHelper.accessor('minSpend', {
      header: t('minSpend'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('maxSpend', {
      header: t('maxSpend'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('example', {
      header: t('example'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('termsAndConditions', {
    //   header: t('termsAndConditions'),
    //   cell: info => info.getValue(),
    // }),
      {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const payment = row.original

      return (
 
        <Button variant={"ghost"} onClick={() => openEditPanel(row.original.ID)}>{t('editPromotion')}</Button>
      )
    },
  },
  ], [])

  const table = useReactTable({
    data: promotions,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
    filterFns: {
      customArrayFilter: (row, columnId, filterValue) => {
        const names = row.getValue(columnId) as string[];
        return names.some(name => 
          name.toLowerCase().includes(filterValue.toLowerCase())
        );
      },
    },
  })

  const handleAddPromotion = () => {
    setEditingPromotion(null);
    setIsAddingPromotion(true);
  };

  const handleCloseAddPromotion = () => {
    setIsAddingPromotion(false);
    setRefreshTrigger(prev => prev + 1);
  };

  const openEditPanel = (id: number) => {
    setEditingPromotion(id);
    setIsAddingPromotion(true);
  };

  const closeEditPanel = () => {
    setEditingPromotion(null);
    setRefreshTrigger(prev => prev + 1);
  };

  if (isLoading) {
    return <div>Loading promotions...</div>;
  }

  return (
    <div className="w-full">
      <div className="flex items-center justify-between mt-4 mb-4">
       
        <Button onClick={handleAddPromotion}>{t('addPromotion')}</Button>
      </div>
      <div className="flex items-center py-4">
        <Input
          placeholder={t('filterPromotion')}
          value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("name")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="ml-auto">
              {t('columns')} <ChevronDownIcon className="ml-2 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {table
              .getAllColumns()
              .filter((column) => column.getCanHide())
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className="capitalize"
                    checked={column.getIsVisible()}
                    onCheckedChange={(value) => {
                      if (value !== column.getIsVisible()) {
                        column.toggleVisibility(!!value)
                      }
                    }}
                  >
                    {t(column.id as string)}
                  </DropdownMenuCheckboxItem>
                )
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
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
                  className="h-24 text-center"
                >
                  {t('noResults')}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-end space-x-2 py-4">
        <div className="flex-1 text-sm text-muted-foreground">
          {table.getFilteredSelectedRowModel().rows.length} {t('of')}{" "}
          {table.getFilteredRowModel().rows.length} {t('rowSelected')}.
        </div>
        <div className="space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            {t('previous')}
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            {t('next')}
          </Button>
        </div>
      </div>
      {editingPromotion && (
        <EditPromotionPanel 
          promotionId={editingPromotion} 
          onClose={closeEditPanel}
          lng={lng}
          prefix={prefix}
        />
      )}
      <EditPromotionPanel
        promotionId={editingPromotion} 
        isOpen={isAddingPromotion}
        lng={lng}
        prefix={prefix}
        onClose={handleCloseAddPromotion}
      
      />
    </div>
  )
}
