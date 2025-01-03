"use client"
 

import React, { useState, useEffect, useMemo,useRef } from 'react'
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
 
import { GetPartnerList } from '@/actions'
import { useTranslation } from '@/app/i18n/client';
 
 
export interface iPartners {
  // Define the properties of GroupedDatabase here
  ID:number,
	username:string,    
	password:string,    
	name:string,    
	bankname:string,    
	banknumber:string,    
	balance:number,    
  affiliatekey:string,
	status:number,    
	 

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
import EditPartner from "./EditPartner"


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

// ตั���อย่างการใช้งาน
// const jsonString = "{\"type\":\"weekly\",\"daysOfWeek\":[\"mon\"],\"hour\":\"11\",\"minute\":\"10\"}";
// console.log(formatSpecificTime(jsonString));

export default function PartnerList({
  prefix,
  data,
  lng,
}: {prefix:string, data:DataTableProps<iPartners>, lng:string}) {
  const [partners, setPartners] = useState<iPartners[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = useState({})
  const [partnerId, setParnerId] = useState<iPartners>(null);
  const [isAddingGame, setIsAddingGame] = useState(false);
  const [isEditingGame, setIsEditingGame] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showTable, setShowTable] = useState(true);
  const router = useRouter()

  const {t} = useTranslation(lng,'translation',undefined)
 // const isSeedFetchedRef = useRef(false);
  useEffect(() => {
    const redirect = ()=>{
      location.replace(`/${lng}/login`)
  }

    const fetchPartners = async () => {
      if (!prefix) {
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const Response = await GetPartnerList(prefix);
        if (Response && Response.Data) {
           console.log(Response.Data)
          setPartners(Response.Data);
        } else {
          console.error('No data received from backend');
        }
      } catch (error) {
        console.error('Error fetching partners:', error);
        redirect();
      } finally {
        setIsLoading(false);
      }
    };
    fetchPartners();
  }, [prefix, refreshTrigger]);

  //console.log('Partners:', partners);

  const columnHelper = createColumnHelper<iPartners>()

  const columns = useMemo(() => [
    columnHelper.accessor('id', {
      header: t('partner.columns.id'),
      cell: info => info.getValue(),
      enableHiding: false,
    }),
    columnHelper.accessor('username', {
      header: t('partner.columns.username'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('name', {
      header: t('partner.columns.fullname'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('bankname', {
      header: t('partner.columns.bankname'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('banknumber', {
      header: t('partner.columns.banknumber'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('Password', {
    //   header: t('columns.username'),
    //   cell: info => info.getValue(),
    // }),
   
    columnHelper.accessor('totalearnings', {
      header: t('partner.columns.totalearnings'),
      cell: info => {
        const value = info.getValue();
        return formatNumber(parseFloat(value?.toString()), 2);
      }
    }),
    columnHelper.accessor('totalcommission', {
      header: t('partner.columns.totalcommission'),
      cell: info => {
        const value = info.getValue();
        return formatNumber(parseFloat(value?.toString()), 2);
      }
    }),
    columnHelper.accessor('balance', {
      header: t('partner.columns.balance'),
      cell: info => {
        const value = info.getValue();
        return formatNumber(parseFloat(value?.toString()), 2);
      }
    }),
      columnHelper.accessor('status', {
        header: t('partner.columns.status'),
        cell: info => {
          const value = info.getValue();
          return value === 1 ? t('common.active') :  t('common.inactive') ;
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
    // columnHelper.accessor('ProStatus', {
    //   header: t('partner.columns.prostatus'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('ProActive', {
    //   header: t('partner.columns.proactive'),
    //   cell: info => info.getValue(),
    // }),
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
        const member = row.original as iPartners;
      //  console.log(member)
        return (
          <div>
            <Button 
              variant="ghost" 
              onClick={() => openEditPanel(member)}
            >
              {t('partner.edit.title')}
            </Button>
          </div>
        );
      },
    },
  ], [])

  const table = useReactTable({
    data: partners,
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
    //setEditingGame(null);
    setIsAddingGame(true);
    setShowTable(false);
  };

  const handleCloseAddGame = () => {
    setIsAddingGame(false);
    setRefreshTrigger((prev:any) => prev + 1);
  };

  const openEditPanel = (member: iPartners) => {
    //console.log(member)
    //console.log(JSON.stringify(member))
    setParnerId(member);
    
    setIsAddingGame(false);
    setShowTable(false);
  };

  const closeEditPanel = () => {
    setParnerId(null);
    setIsAddingGame(false);
    setShowTable(true);
    setRefreshTrigger((prev:any) => prev + 1);
  };

  if (isLoading) {
    return <div>Loading {t('partner.title')}...</div>;
  }

  return (
    <div className="w-full">
      {showTable ? (
        <>
          <div className="flex items-center justify-between mt-4 mb-4">
            <Button onClick={handleAddGame}>{t('partner.add.title')}</Button>
          </div>
          <div className="flex items-center py-4">
            <Input
              placeholder={t('common.search')}
              value={(table.getColumn("username")?.getFilterValue() as string) ?? ""}
              onChange={(event:any) =>
                table.getColumn("username")?.setFilterValue(event.target.value)
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
                  .filter((column:any) => column.getCanHide())
                  .map((column:any) => {
                    return (
                      <DropdownMenuCheckboxItem
                        key={column.id}
                        // className="capitalize"
                        checked={column.getIsVisible()}
                        onCheckedChange={(value:any) => {
                          if (value !== column.getIsVisible()) {
                            column.toggleVisibility(!!value)
                          }
                        }}
                      >
                        {t(`partner.columns.${(column.id as string).toLowerCase()}`)}
                      </DropdownMenuCheckboxItem>
                    )
                  })}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                {table.getHeaderGroups().map((headerGroup:any) => (
                  <TableRow key={headerGroup.id}>
                    {headerGroup.headers.map((header:any) => {
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
                  table.getRowModel().rows.map((row:any) => (
                    <TableRow
                      key={row.id}
                      data-state={row.getIsSelected() && "selected"}
                    >
                      {row.getVisibleCells().map((cell:any) => (
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
            {t('partner.columns.backToList')}
          </Button>
          <EditPartner
            partnerId={partnerId.id}
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
