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
//import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTrigger } from '@/components/ui/dialog'

import EditPromotionPanel from './EditPromotionPanel'
import { GetGameStatus, GetPromotion,DeletePromotion} from '@/actions'
import { useTranslation } from '@/app/i18n/client'; 
import { ResponsiveDialog } from '../responsive-dialog'
import DeletePromotionDialog from './DeletePromotionDialog'
import { format, parse } from 'date-fns';

export interface Promotion {
  id: string
  ID: string
  name: string
  percentDiscount: number
  maxDiscount: number
  usageLimit: string
  specificTime: string
  paymentMethod: string
  minDept:number
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
      days = data.daysOfWeek.map((day: string) => daysMap[day] || day).join(', ');
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
  
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showTable, setShowTable] = useState(true);
 
  const [openDropdown, setOpenDropdown] = useState<string | null>(null)
  const [selectedPromotionId, setSelectedPromotionId] = useState<string | null>(null);

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
    columnHelper.accessor('minDept', {
      header: t('promotion.minDept'),
      cell: info => `${info.getValue()}`,
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
          const game: any = Object.values(gameTypes.Data).find((g: any) => {
            try {
              const status = g.status;
              return status.id.toString() === id.trim();
            } catch (e) {
              console.error('Error parsing game status:', e);
              return false;
            }
          });
          if (game) {
            try {
              const status = game.status;
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
          const game: any = Object.values(gameTypes.Data).find((g: any) => {
            try {
              const status = g.status;
              return status.id.toString() === id.trim();
            } catch (e) {
              console.error('Error parsing game status:', e);
              return false;
            }
          });
          if (game) {
            try {
              const status = game.status;
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
      cell: info => {
        const date = info.getValue();
        if (!date) return '';
        // Parse the date string from backend format and then format it to desired format
        return format(parse(date, 'yyyy-MM-dd', new Date()), 'dd-MM-yyyy');
      },
    }),
    columnHelper.accessor('endDate', {
      header: t('promotion.endDate'),
      cell: info => {
        const date = info.getValue();
        if (!date) return '';
        // Parse the date string from backend format and then format it to desired format
        return format(parse(date, 'yyyy-MM-dd', new Date()), 'dd-MM-yyyy');
      },
    }),
    // columnHelper.accessor('paymentMethod', {
    //   header: t('paymentMethod'),
    //   cell: info => info.getValue(),
    // }),
    columnHelper.accessor('minSpend', {
      header: t('promotion.minSpend') + "%",
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
    cell: ({ row }: { row: any }  ) => {
      const promotion = row.original

      return (
        <DropdownMenu open={openDropdown === promotion.ID} onOpenChange={(open) => {
          if (open) setOpenDropdown(promotion.ID)
          else setOpenDropdown(null)
        }}>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost">{t('promotion.edit.options')}</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => openEditPanel(promotion.ID)}>
              {t('promotion.edit.edit')}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem  
              onClick={() => {
                setSelectedPromotionId(promotion.ID);
                setIsDeleteOpen(true);
              }}
              className="w-full justify-start flex text-red-500 rounded-md p-2 transition-all duration-75 hover:bg-neutral-100"
            >
              {t('promotion.edit.delete')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
  ], [gameTypes, t, openDropdown])

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
   console.log('openEditPanel', id);
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
    return <div>Loading {t('promotion.title')}...</div>;
  }

  return (
    <div className="w-full">
      <ResponsiveDialog
        lng={lng}
        isOpen={isDeleteOpen}
        setIsOpen={setIsDeleteOpen}
        title={t('promotion.edit.delete_description')}
        ns='translation'
        prefixkey='promotion'
      >
        <DeletePromotionDialog 
          prefix={prefix} 
          lng={lng} 
          promotionId={selectedPromotionId ?? ''}
          setIsOpen={setIsDeleteOpen} 
          onDeleteSuccess={() => {
            setRefreshTrigger(prev => prev + 1);
          }}
        />
      </ResponsiveDialog>

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
                {table.getRowModel()?.rows?.length ? (
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
            promotionId={editingPromotion ?? null}
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
