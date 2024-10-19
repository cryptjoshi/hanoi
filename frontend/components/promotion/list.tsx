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

import { ArrowLeftIcon } from "@radix-ui/react-icons"
import { useQuery } from '@tanstack/react-query'
import { Dialog, DialogContent, DialogFooter, DialogTrigger } from '@/components/ui/dialog'

import EditPromotionPanel from './EditPromotionPanel'
import { GetGameStatus, GetPromotion } from '@/actions'
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
  includegames: string
  excludegames: string
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



function formatSpecificTime(jsonString: string,lng:string): string {
  if (!jsonString) {
    return '';
  }

  try {
    // Remove any leading/trailing whitespace and quotes
    const {t} = useTranslation(lng,'translation',{keyPrefix:'common'})
    let cleanJsonString = jsonString.trim().replace(/^["']|["']$/g, '');
    
    // Replace escaped quotes with regular quotes
    cleanJsonString = cleanJsonString.replace(/\\"/g, '"');
    
    // Parse the cleaned JSON string
    const data = JSON.parse(cleanJsonString);
    
    const daysMap: { [key: string]: string } = {
      'mon': t('mon'),
      'tue': t('tue'),
      'wed': t('wed'),
      'thu': t('thu'),
      'fri': t('fri'),
      'sat': t('sat'),
      'sun': t('sun')
    };

    const typeMap: { [key: string]: string } = {
      'weekly': t('weekly'),
      'once': t('once'),
      'monthly': t('monthly')
    };

    let days = '';
    if (data.daysOfWeek && data.daysOfWeek.length > 0) {
      days = data.daysOfWeek.map(day => daysMap[day] || day).join(', ');
    }

    let frequency = typeMap[data.type] || '';
    let time = '';

    if (data.hour && data.minute) {
      time = ` ${data.hour.padStart(2, '0')}:${data.minute.padStart(2, '0')}`;
    }

    const parts = [days, frequency, time].filter(Boolean);
    return parts.join(' ');
  } catch (error) {
   // console.error('Error parsing specificTime:', error);
    //console.log('Problematic string:', jsonString);
    return '';
  }
}

// ตัวอย่างการใช้งาน
// const jsonString = "{\"type\":\"weekly\",\"daysOfWeek\":[\"mon\"],\"hour\":\"11\",\"minute\":\"10\"}";
// console.log(formatSpecificTime(jsonString));

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
  const [showTable, setShowTable] = useState(true);
  const router = useRouter()

  const {t} = useTranslation(lng,'translation',undefined)
  const { data: gameTypes, isLoading: gameStatusLoading } = useQuery({
    queryKey: ['gameTypes'],
    queryFn: async () => await GetGameStatus(prefix),
  });
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
      header: t('promotion.ID'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('name', {
      header: t('promotion.name'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('percentDiscount', {
      header: t('promotion.percentDiscount'),
      cell: info => `${info.getValue()}%`,
    }),
    columnHelper.accessor('maxDiscount', {
      header: t('promotion.maxDiscount'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('usageLimit', {
      header: t('promotion.usageLimit'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('specificTime', {
      header: t('promotion.specificTime'),
      cell: info => {
        const value = info.getValue();
        return typeof value === 'string' ? formatSpecificTime(value,lng) : '';
      }
    }),
    columnHelper.accessor('includegames', {
      header: t('promotion.include'),
      cell: info => {
        const value = info.getValue();
        if (!gameTypes || !gameTypes.Data || typeof gameTypes.Data !== 'object') return value;
        return value.split(',').map(id => {
          const game = Object.values(gameTypes.Data).find((g: any) => {
            try {
              const status = JSON.parse(g.status);
              return status.id.toString() === id.trim();
            } catch (e) {
              console.error('Error parsing game status:', e);
              return false;
            }
          });
          if (game) {
            try {
              const status = JSON.parse(game.status);
              return t(`games.${status.name}`);
            } catch (e) {
              console.error('Error parsing game status:', e);
              return id;
            }
          }
          return id;
        }).join(', ');
      },
    }),
    columnHelper.accessor('excludegames', {
      header: t('promotion.exclude'),
      cell: info => {
        const value = info.getValue();
        if (!gameTypes || !gameTypes.Data || typeof gameTypes.Data !== 'object') return value;
        return value.split(',').map(id => {
          const game = Object.values(gameTypes.Data).find((g: any) => {
            try {
              const status = JSON.parse(g.status);
              return status.id.toString() === id.trim();
            } catch (e) {
              console.error('Error parsing game status:', e);
              return false;
            }
          });
          if (game) {
            try {
              const status = JSON.parse(game.status);
              return t(`games.${status.name}`);
            } catch (e) {
              console.error('Error parsing game status:', e);
              return id;
            }
          }
          return id;
        }).join(', ');
      },
    }),
    columnHelper.accessor('startDate', {
      header: t('promotion.startDate'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('endDate', {
      header: t('promotion.endDate'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('paymentMethod', {
    //   header: t('paymentMethod'),
    //   cell: info => info.getValue(),
    // }),
    columnHelper.accessor('minSpend', {
      header: t('promotion.minSpend'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('maxSpend', {
      header: t('promotion.maxSpend'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('example', {
      header: t('promotion.example'),
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
            // <DropdownMenu>
            //   <DropdownMenuTrigger asChild>
            //     <Button variant={"ghost"}>{t('promotion.edit.options')}</Button>
            //   </DropdownMenuTrigger>  
            //   <DropdownMenuContent>
            //     <DropdownMenuItem  onClick={() => openEditPanel(row.original.ID)}>{t('promotion.edit.edit')}</DropdownMenuItem>
            //     <DropdownMenuItem onClick={() => openEditPanel(row.original.ID)}>{t('promotion.edit.delete')}</DropdownMenuItem>
            //   </DropdownMenuContent>
            // </DropdownMenu>
            <DropdownMenu>
            <DropdownMenuTrigger> <Button variant={"ghost"}>{t('promotion.edit.options')}</Button></DropdownMenuTrigger>
            <DropdownMenuContent>
            <DropdownMenuItem  onClick={() => openEditPanel(row.original.ID)}>{t('promotion.edit.edit')}</DropdownMenuItem>
            <DropdownMenuSeparator />
              <Dialog>
                <DialogTrigger asChild>
                  <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                  {t('promotion.edit.delete')}
                  </DropdownMenuItem>
                </DialogTrigger>
                <DialogContent>
                  {t('promotion.edit.delete_description')}
                  <DialogFooter className='flex gap-2 justfy-between'>
                    <Button onClick={() => openEditPanel(row.original.ID)}>{t('promotion.edit.delete')}</Button>
                    <Button onClick={() => openEditPanel(row.original.ID)}>{t('promotion.edit.cancel')}</Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </DropdownMenuContent>
          </DropdownMenu>
       
      )
    },
  },
  ], [gameTypes, t])

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
    setShowTable(false);
  };

  const handleCloseAddPromotion = () => {
    setIsAddingPromotion(false);
    setRefreshTrigger(prev => prev + 1);
  };

  const openEditPanel = (id: number) => {
    setEditingPromotion(id);
    setIsAddingPromotion(true);
    setShowTable(false);
  };

  const closeEditPanel = () => {
    setEditingPromotion(null);
    setIsAddingPromotion(false);
    setShowTable(true);
    setRefreshTrigger(prev => prev + 1);
  };

  if (isLoading) {
    return <div>Loading promotions...</div>;
  }

  return (
    <div className="w-full">
      {showTable ? (
        <>
          <div className="flex items-center justify-between mt-4 mb-4">
            <Button onClick={handleAddPromotion}>{t('promotion.addPromotion')}</Button>
          </div>
          <div className="flex items-center py-4">
            <Input
              placeholder={t('promotion.filterPromotion')}
              value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
              onChange={(event) =>
                table.getColumn("name")?.setFilterValue(event.target.value)
              }
              className="max-w-sm"
            />
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" className="ml-auto">
                  {t('promotion.columns.columns')} <ChevronDownIcon className="ml-2 h-4 w-4" />
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
                      {t('promotion.noResults')}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
          <div className="flex items-center justify-end space-x-2 py-4">
            <div className="flex-1 text-sm text-muted-foreground">
              {table.getFilteredSelectedRowModel().rows.length} {t('promotion.of')}{" "}
              {table.getFilteredRowModel().rows.length} {t('promotion.rowSelected')}.
            </div>
            <div className="space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => table.previousPage()}
                disabled={!table.getCanPreviousPage()}
              >
                {t('promotion.previous')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => table.nextPage()}
                disabled={!table.getCanNextPage()}
              >
                {t('promotion.next')}
              </Button>
            </div>
          </div>
        </>
      ) : (
        <div className="mt-4">
          <Button
            variant="outline"
            onClick={closeEditPanel}
            className="mb-4"
          >
            <ArrowLeftIcon className="mr-2 h-4 w-4" />
            {t('promotion.backToList')}
          </Button>
          <EditPromotionPanel
            promotionId={editingPromotion}
            lng={lng}
            prefix={prefix}
            onClose={closeEditPanel}
            onCancel={closeEditPanel}
          />
        </div>
      )}
    </div>
  )
}
