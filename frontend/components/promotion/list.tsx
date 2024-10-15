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
}: {prefix:string, data:DataTableProps<GroupedDatabase>}) {
  const [promotions, setPromotions] = useState<Promotion[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = useState({})
  const [editingPromotion, setEditingPromotion] = useState<string | null>(null);
  const [isAddingPromotion, setIsAddingPromotion] = useState(false);
  const router = useRouter()

  useEffect(() => {
    const fetchPromotions = async () => {
      if (!prefix) {
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const fetchedPromotions = await GetPromotion(prefix);
        const serializablePromotions = fetchedPromotions.map((promo: any) => ({
          id: promo.id,
          name: promo.name,
          percentDiscount: promo.percentDiscount,
          maxDiscount: promo.maxDiscount,
          usageLimit: promo.usageLimit,
          specificTime: promo.specificTime,
          paymentMethod: promo.paymentMethod,
          minSpend: promo.minSpend,
          maxSpend: promo.maxSpend,
          termsAndConditions: promo.termsAndConditions,
        }));
        setPromotions(serializablePromotions);
      } catch (error) {
        console.error('Error fetching promotions:', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchPromotions();
  }, [prefix])

  const columnHelper = createColumnHelper<Promotion>()

  const columns = useMemo(() => [
    columnHelper.accessor('id', {
      header: 'ID',
      cell: info => info.getValue(),
    }),
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
      {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const payment = row.original

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <DotsHorizontalIcon className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(payment.id)}
            >
              Copy Promotion ID
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => openEditPanel(payment.id)}>แก้ไขโปรโมชั่น</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
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
    setIsAddingPromotion(true);
  };

  const handleCloseAddPromotion = () => {
    setIsAddingPromotion(false);
  };

  const openEditPanel = (id: string) => {
    setEditingPromotion(id);
  };

  const closeEditPanel = () => {
    setEditingPromotion(null);
  };

  if (isLoading) {
    return <div>Loading promotions...</div>;
  }

  return (
    <div className="w-full">
      <div className="flex items-center justify-between mt-4 mb-4">
       
        <Button onClick={handleAddPromotion}>เพิ่มโปรโมชั่น</Button>
      </div>
      <div className="flex items-center py-4">
        <Input
          placeholder="Filter databases..."
          value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("name")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="ml-auto">
              Columns <ChevronDownIcon className="ml-2 h-4 w-4" />
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
                    onCheckedChange={(value) =>
                      column.toggleVisibility(!!value)
                    }
                  >
                    {column.id}
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
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-end space-x-2 py-4">
        <div className="flex-1 text-sm text-muted-foreground">
          {table.getFilteredSelectedRowModel().rows.length} of{" "}
          {table.getFilteredRowModel().rows.length} row(s) selected.
        </div>
        <div className="space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
          </Button>
        </div>
      </div>
      {editingPromotion && (
        <EditPromotionPanel 
          promotionId={editingPromotion} 
          onClose={closeEditPanel}
        />
      )}
      <EditPromotionPanel
        isOpen={isAddingPromotion}
        onClose={handleCloseAddPromotion}
        promotionId={null}
      />
    </div>
  )
}