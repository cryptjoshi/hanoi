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
 
import { GetMemberList } from '@/actions'
import { useTranslation } from '@/app/i18n/client';
 
 
export interface iMember {
  // Define the properties of GroupedDatabase here
  ID:number,
  Walletid:number,       
	Username:string,    
	Password:string,    
	ProviderPassword:string,    
	Fullname:string,    
	Bankname:string,    
	Banknumber:string,    
	Balance:number,    
	Beforebalance:number,    
	Token:string,    
	Role:string,    
	Salt:string,    
	Status:number,    
	Betamount:number,    
	Win:number,    
	Lose:number,    
	Turnover:number,    
	ProID:string,    
	PartnersKey:string,    
	ProStatus:string,    
  ProActive:string

  // Add other properties as needed
}
interface DataTableProps<TData> {
  columns: ColumnDef<TData, any>[]
  data: TData[]
}

import { ArrowLeftIcon } from "@radix-ui/react-icons"
import EditMember from './EditPartner'
import { number } from 'zod'
import { formatNumber } from '@/lib/utils'
import EditPartner from './EditPartner'


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

export default function PartnerList({
  prefix,
  data,
  lng,
}: {prefix:string, data:DataTableProps<iMember>, lng:string}) {
  const [games, setGames] = useState<iMember[]>([])
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

  const {t} = useTranslation(lng,'translation',undefined)

  useEffect(() => {
    const fetchGames = async () => {
      if (!prefix) {
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const fetchedGames = await GetMemberList(prefix);
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

  const columnHelper = createColumnHelper<iMember>()

  const columns = useMemo(() => [
    columnHelper.accessor('ID', {
      header: t('member.columns.id'),
      cell: info => info.getValue(),
      enableHiding: false,
    }),
    columnHelper.accessor('Username', {
      header: t('member.columns.username'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('Fullname', {
      header: t('member.columns.fullname'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('Bankname', {
      header: t('member.columns.bankname'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('Banknumber', {
      header: t('member.columns.banknumber'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('Password', {
    //   header: t('columns.username'),
    //   cell: info => info.getValue(),
    // }),
    columnHelper.accessor('Balance', {
      header: t('member.columns.balance'),
      cell: info => {
        const value = info.getValue();
        return formatNumber(parseFloat(value?.toString()), 2);
      }
    }),
      columnHelper.accessor('Status', {
        header: t('member.columns.status'),
        cell: info => {
          const value = info.getValue();
          return value === 1 ? t('common.active') : value === 0 ? t('common.inactive') : t('common.maintenance');
        //  console.log('Raw specificTime value:', value); // For debugging
         
        }
      }),
 
    // columnHelper.accessor('Active', {
    //   header: t('columns.active'),
    //   cell: info => {
    //     const value = info.getValue();
    //     return value === 1 ? t('active') : value === 0 ? t('inactive') : t('maintenance');
    //   }
    // }),
    columnHelper.accessor('ProStatus', {
      header: t('member.columns.prostatus'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('ProActive', {
      header: t('member.columns.proactive'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('position', {
    //   header: t('columns.position'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('urlimage', {
    //   header: t('columns.urlimage'),
    //   cell: info => info.getValue(),
    // }),
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
        const member = row.original as iMember;
        return (
          <div>
            <Button 
              variant="ghost" 
              onClick={() => openEditPanel(member)}
            >
              {t('member.edit.title')}
            </Button>
          </div>
        );
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

  const openEditPanel = (member: iMember) => {
  
    setEditingGame(member.ID);
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
    return <div>Loading {t('member.title')}...</div>;
  }

  return (
    <div className="w-full">
      {showTable ? (
        <>
          <div className="flex items-center justify-between mt-4 mb-4">
            <Button onClick={handleAddGame}>{t('member.add.title')}</Button>
          </div>
          <div className="flex items-center py-4">
            <Input
              placeholder={t('member.columns.search')}
              value={(table.getColumn("Username")?.getFilterValue() as string) ?? ""}
              onChange={(event) =>
                table.getColumn("Username")?.setFilterValue(event.target.value)
              }
              className="max-w-sm"
            />
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" className="ml-auto">
                  {t('common.columnsfilter')} <ChevronDownIcon className="ml-2 h-4 w-4" />
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
                        // className="capitalize"
                        checked={column.getIsVisible()}
                        onCheckedChange={(value) => {
                          if (value !== column.getIsVisible()) {
                            column.toggleVisibility(!!value)
                          }
                        }}
                      >
                        {t(`member.columns.${(column.id as string).toLowerCase()}`)}
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
              {table.getFilteredSelectedRowModel().rows.length} {t('common.of')}{" "}
              {table.getFilteredRowModel().rows.length} {t('common.rowSelected')}.
            </div>
            <div className="space-x-2">
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
        </>
      ) : (
        <div className="mt-4">
          <Button
            variant="outline"
            onClick={closeEditPanel}
            className="mb-4"
          >
            <ArrowLeftIcon className="mr-2 h-4 w-4" />
            {t('member.columns.backToList')}
          </Button>
          <EditPartner
            partnerId={editingGame || 0}
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
