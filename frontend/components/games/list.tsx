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
import EditGamesPanel from './EditGamesPanel'
import { GetGameList } from '@/actions'
import { useTranslation } from '@/app/i18n/client';
import { Games } from '@/types';
 
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

import { ArrowLeftIcon } from "@radix-ui/react-icons"
import EditGame from './EditGame'


function formatSpecificTime(jsonString: string,lng:string): string {
  if (!jsonString) {
    return '';
  }

  try {
    // Remove any leading/trailing whitespace and quotes
    const {t} = useTranslation(lng,'translation',{keyPrefix:'games'})
    let cleanJsonString = jsonString.trim().replace(/^["']|["']$/g, '');
    
    // Replace escaped quotes with regular quotes
    cleanJsonString = cleanJsonString.replace(/\\"/g, '"');
    
    // Parse the cleaned JSON string
    const data = JSON.parse(cleanJsonString);
   
    return  t(data.name)
  } catch (error) {
   // console.error('Error parsing specificTime:', error);
    //console.log('Problematic string:', jsonString);
    return '';
  }
}

// ตัวอย่างการใช้งาน
// const jsonString = "{\"type\":\"weekly\",\"daysOfWeek\":[\"mon\"],\"hour\":\"11\",\"minute\":\"10\"}";
// console.log(formatSpecificTime(jsonString));

export default function GameListDataTable({
  prefix,
  data,
  lng,
}: {prefix:string, data:DataTableProps<GroupedDatabase>, lng:string}) {
  const [games, setGames] = useState<Games[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = useState({})
  const [editingGame, setEditingGame] = useState<number | null>(null);
  const [isAddingGame, setIsAddingGame] = useState(false);
  const [isEditingGame, setIsEditingGame] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showTable, setShowTable] = useState(true);
  const router = useRouter()

  const {t} = useTranslation(lng,'translation',{keyPrefix:'games'})

  useEffect(() => {
    const fetchGames = async () => {
      if (!prefix) {
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const fetchedGames = await GetGameList(prefix);
        //onsole.log(fetchedGames)
        setGames(fetchedGames.Data);
      } catch (error) {
        console.error('Error fetching games:', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchGames();
  }, [prefix, refreshTrigger])

  const columnHelper = createColumnHelper<Games>()

  const columns = useMemo(() => [
    columnHelper.accessor('ID', {
      header: t('columns.id'),
      cell: info => info.getValue(),
      enableHiding: false,
    }),
    columnHelper.accessor('name', {
      header: t('columns.name'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('productCode', {
      header: t('columns.productCode'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('product', {
      header: t('columns.product'),
      cell: info => info.getValue(),
    }),
      columnHelper.accessor('gameType', {
        header: t('columns.gameType'),
        cell: info => {
          const value = info.row.original.status;
        //  console.log('Raw specificTime value:', value); // For debugging
          return typeof value === 'string' ? formatSpecificTime(value,lng) : '';
        }
      }),
    columnHelper.accessor('active', {
      header: t('columns.active'),
      cell: info => {
        const value = info.getValue();
        return value === 1 ? t('active') : value === 0 ? t('inactive') : t('maintenance');
      }
    }),
    columnHelper.accessor('remark', {
      header: t('columns.remark'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('position', {
      header: t('columns.position'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('urlimage', {
      header: t('columns.urlimage'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('status', {
    //   header: t('columns.status'),
    //   cell: info => {
    //     const value = info.getValue();
    //   //  console.log('Raw specificTime value:', value); // For debugging
    //     return typeof value === 'string' ? formatSpecificTime(value,lng) : '';
    //   }
    // }),
    // }),
    // columnHelper.accessor('startDate', {
    //   header: t('startDate'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('endDate', {
    //   header: t('endDate'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('paymentMethod', {
    //   header: t('paymentMethod'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('minSpend', {
    //   header: t('minSpend'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('maxSpend', {
    //   header: t('maxSpend'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('example', {
    //   header: t('example'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('termsAndConditions', {
    //   header: t('termsAndConditions'),
    //   cell: info => info.getValue(),
    // }),
    {
      id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
        

      return (
 
        <Button variant={"ghost"} onClick={() => openEditPanel(row.original.ID)}>{t('editGame')}</Button>
      )
      },
    },
  ], [])

  const table = useReactTable({
    data: games,
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

  const handleAddGame = () => {
    setEditingGame(null);
    setIsAddingGame(true);
    setShowTable(false);
  };

  const handleCloseAddGame = () => {
    setIsAddingGame(false);
    setRefreshTrigger(prev => prev + 1);
  };

  const openEditPanel = (id: number) => {
  
    setEditingGame(id);
    setIsAddingGame(false);
    setShowTable(false);
  };

  const closeEditPanel = () => {
    setEditingGame(null);
    setIsAddingGame(false);
    setShowTable(true);
    setRefreshTrigger(prev => prev + 1);
  };

  if (isLoading) {
    return <div>Loading games...</div>;
  }

  return (
    <div className="w-full">
      {showTable ? (
        <>
          <div className="flex items-center justify-between mt-4 mb-4">
            <Button onClick={handleAddGame}>{t('addGame')}</Button>
          </div>
          <div className="flex items-center py-4">
            <Input
              placeholder={t('filterGame')}
              value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
              onChange={(event) =>
                table.getColumn("name")?.setFilterValue(event.target.value)
              }
              className="max-w-sm"
            />
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" className="ml-auto">
                  {t('columnsfilter')} <ChevronDownIcon className="ml-2 h-4 w-4" />
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
                        {t(`columns.${column.id as string}`)}
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
        </>
      ) : (
        <div className="mt-4">
          <Button
            variant="outline"
            onClick={closeEditPanel}
            className="mb-4"
          >
            <ArrowLeftIcon className="mr-2 h-4 w-4" />
            {t('backToList')}
          </Button>
          <EditGame
            gameId={editingGame}
            isAdd={isAddingGame}
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
