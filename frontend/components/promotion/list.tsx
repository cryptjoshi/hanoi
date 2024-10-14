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
  names: string;
  prefix: string;
  // Add other properties as needed
}
interface DataTableProps<TData> {
  columns: ColumnDef<TData, any>[]
  data: TData[]
}


export default function PromotionListDataTable({
 
  data,
}: DataTableProps<GroupedDatabase>) {
  const [promotions, setPromotions] = useState<Promotion[]>([])
  const router = useRouter()
 
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = useState({})
  const [editingPromotion, setEditingPromotion] = useState<string | null>(null);
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

 

  const handleAddPromotion = () => {
    router.push('/promotions/add')
  }

  const openEditPanel = (id: string) => {
    setEditingPromotion(id);
  };

  const closeEditPanel = () => {
    setEditingPromotion(null);
  };

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

  return (
    <div className="w-full">
      <div className="flex items-center py-4">
        <Input
          placeholder="Filter databases..."
          value={(table.getColumn("names")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("names")?.setFilterValue(event.target.value)
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
    </div>
  )
}
// function DataTable<TData, TValue>({
//   columns,
//   data,
// }: DataTableProps<TData, TValue>) {
//   const table = useReactTable({
//     data,
//     columns,
//     getCoreRowModel: getCoreRowModel(),
//   })

//   return (
//     <div className="rounded-md border">
//       <Table>
//         <TableHeader>
//           {table.getHeaderGroups().map((headerGroup) => (
//             <TableRow key={headerGroup.id}>
//               {headerGroup.headers.map((header) => {
//                 return (
//                   <TableHead key={header.id}>
//                     {header.isPlaceholder
//                       ? null
//                       : flexRender(
//                           header.column.columnDef.header,
//                           header.getContext()
//                         )}
//                   </TableHead>
//                 )
//               })}
//             </TableRow>
//           ))}
//         </TableHeader>
//         <TableBody>
//           {table.getRowModel().rows?.length ? (
//             table.getRowModel().rows.map((row) => (
//               <TableRow
//                 key={row.id}
//                 data-state={row.getIsSelected() && "selected"}
//               >
//                 {row.getVisibleCells().map((cell) => (
//                   <TableCell key={cell.id}>
//                     {flexRender(cell.column.columnDef.cell, cell.getContext())}
//                   </TableCell>
//                 ))}
//               </TableRow>
//             ))
//           ) : (
//             <TableRow>
//               <TableCell colSpan={columns.length} className="h-24 text-center">
//                 No results.
//               </TableCell>
//             </TableRow>
//           )}
//         </TableBody>
//       </Table>
//     </div>
//   )
// } function DataTableDemo<TData, TValue>({
//     columns,
//     data,
//   }: DataTableProps<TData, TValue>){
    
//     const [sorting, setSorting] = React.useState<SortingState>([])
//     const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
//       []
//     )
//     const [columnVisibility, setColumnVisibility] =
//       React.useState<VisibilityState>({})
//     const [rowSelection, setRowSelection] = React.useState({})
  
//     const table = useReactTable({
//       data,
//       columns,
//       onSortingChange: setSorting,
//       onColumnFiltersChange: setColumnFilters,
//       getCoreRowModel: getCoreRowModel(),
//       getPaginationRowModel: getPaginationRowModel(),
//       getSortedRowModel: getSortedRowModel(),
//       getFilteredRowModel: getFilteredRowModel(),
//       onColumnVisibilityChange: setColumnVisibility,
//       onRowSelectionChange: setRowSelection,
//       state: {
//         sorting,
//         columnFilters,
//         columnVisibility,
//         rowSelection,
//       },
//     })
  
//     return (
//       <div className="w-full">
//         <div className="flex items-center py-4">
//           <Input
//             placeholder="Filter emails..."
//             value={(table.getColumn("email")?.getFilterValue() as string) ?? ""}
//             onChange={(event) =>
//               table.getColumn("email")?.setFilterValue(event.target.value)
//             }
//             className="max-w-sm"
//           />
//           <DropdownMenu>
//             <DropdownMenuTrigger asChild>
//               <Button variant="outline" className="ml-auto">
//                 Columns <ChevronDownIcon className="ml-2 h-4 w-4" />
//               </Button>
//             </DropdownMenuTrigger>
//             <DropdownMenuContent align="end">
//               {table
//                 .getAllColumns()
//                 .filter((column) => column.getCanHide())
//                 .map((column) => {
//                   return (
//                     <DropdownMenuCheckboxItem
//                       key={column.id}
//                       className="capitalize"
//                       checked={column.getIsVisible()}
//                       onCheckedChange={(value) =>
//                         column.toggleVisibility(!!value)
//                       }
//                     >
//                       {column.id}
//                     </DropdownMenuCheckboxItem>
//                   )
//                 })}
//             </DropdownMenuContent>
//           </DropdownMenu>
//         </div>
//         <div className="rounded-md border">
//           <Table>
//             <TableHeader>
//               {table.getHeaderGroups().map((headerGroup) => (
//                 <TableRow key={headerGroup.id}>
//                   {headerGroup.headers.map((header) => {
//                     return (
//                       <TableHead key={header.id}>
//                         {header.isPlaceholder
//                           ? null
//                           : flexRender(
//                               header.column.columnDef.header,
//                               header.getContext()
//                             )}
//                       </TableHead>
//                     )
//                   })}
//                 </TableRow>
//               ))}
//             </TableHeader>
//             <TableBody>
//               {table.getRowModel().rows?.length ? (
//                 table.getRowModel().rows.map((row) => (
//                   <TableRow
//                     key={row.id}
//                     data-state={row.getIsSelected() && "selected"}
//                   >
//                     {row.getVisibleCells().map((cell) => (
//                       <TableCell key={cell.id}>
//                         {flexRender(
//                           cell.column.columnDef.cell,
//                           cell.getContext()
//                         )}
//                       </TableCell>
//                     ))}
//                   </TableRow>
//                 ))
//               ) : (
//                 <TableRow>
//                   <TableCell
//                     colSpan={columns.length}
//                     className="h-24 text-center"
//                   >
//                     No results.
//                   </TableCell>
//                 </TableRow>
//               )}
//             </TableBody>
//           </Table>
//         </div>
//         <div className="flex items-center justify-end space-x-2 py-4">
//           <div className="flex-1 text-sm text-muted-foreground">
//             {table.getFilteredSelectedRowModel().rows.length} of{" "}
//             {table.getFilteredRowModel().rows.length} row(s) selected.
//           </div>
//           <div className="space-x-2">
//             <Button
//               variant="outline"
//               size="sm"
//               onClick={() => table.previousPage()}
//               disabled={!table.getCanPreviousPage()}
//             >
//               Previous
//             </Button>
//             <Button
//               variant="outline"
//               size="sm"
//               onClick={() => table.nextPage()}
//               disabled={!table.getCanNextPage()}
//             >
//               Next
//             </Button>
//           </div>
//         </div>
//       </div>
//     )
//   }

 
