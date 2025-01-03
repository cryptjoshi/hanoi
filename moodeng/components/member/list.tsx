"use client"
 

import React, { useState, useEffect, useMemo } from 'react'
import { useRouter } from 'next/navigation'
import { toast } from "@/hooks/use-toast"
import { cn } from "@/lib/utils"
import {
  CaretSortIcon,
  ChevronDownIcon,
  DotsHorizontalIcon,
} from "@radix-ui/react-icons"

import { PlusIcon } from "@radix-ui/react-icons"

import { TZDate } from "@date-fns/tz";
import { addDays } from "date-fns";
import { formatNumber,decompressGzip,Decrypt } from '@/lib/utils'
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

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
import { CalendarDateRangePicker } from '@/components/date-range-picker';
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
import { DateRange } from "react-day-picker"
 
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
  Commission:number,  
	Win:number,    
	Lose:number,    
	Turnover:number,    
	ProID:string,    
	PartnersKey:string,    
	ProStatus:string,    
  ProActive:string,
  Prefix:string,
  TotalTurnover:number,
  TotalEarnings:number,
  TDeposit:number,
  Deposit:number,
  Withdraw:number,
  TWithdraw:number,
  Crdb:number,
  Winlose:number,
  SumProamount:number


  // Add other properties as needed
}
export interface DataTableProps<TData> {
  columns: ColumnDef<TData, any>[]
  data: TData[]
  rows: []
}

import { ArrowLeftIcon } from "@radix-ui/react-icons"
import EditMember from './EditMember'
import { number } from 'zod'
 
import useAuthStore from '@/store/auth'
import { Form, FormField, FormItem } from '../ui/form'


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

export default function MemberListDataTable({
  id,
  data,
  lng,
}: { id:string,data:DataTableProps<iMember>, lng:string}) {
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
  const [prefix,setPrefix] = useState(null)
  const router = useRouter()
  //const { customerCurrency,accessToken } = useAuthStore();
  const {t} = useTranslation(lng,'translation',undefined)
  const tzDate = new TZDate(new Date(), "Asia/Bangkok");
  const FormSchema = z.object({
    startdate: z.object({
      from: z.date({
      required_error: "A date of Start is required.",
    }),
    to: z.date({
      required_error: "A date of Start is required.",
    }),
  })})

  
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
        startdate: {
          from: tzDate,
          to: addDays(tzDate, 20),
        } // ตั้งค่าเริ่มต้นให้กับ startdate
    },
  })

  const [date, setDate] = React.useState<DateRange | undefined>
  ({
    from: tzDate,
    to: addDays(tzDate, 20),
  })

  useEffect(() => {
    const fetchGames = async () => {
      
      setIsLoading(true);
      try {
      
      //  if(accessToken){
        const fetchedGames = await GetMemberList(date);
        //console.log(fetchedGames)
        setGames(fetchedGames.Data);
       // } else {
       //   router.replace(`/${lng}/login`)
       // }
      } catch (error) {
        console.error('Error fetching games:', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchGames();
  }, [ refreshTrigger])

  const columnHelper = createColumnHelper<iMember>()

  const columns = useMemo(() => [
    // columnHelper.accessor('ID', {
    //   header: t('member.columns.id'),
    //   cell: info => info.getValue(),
    //   enableHiding: false,
    // }),
    // columnHelper.accessor('ReferralCode', {
    //   header: t('member.columns.referralcode'),
    //   cell: info => info.getValue(),
    // }),
    columnHelper.accessor('Username', {
      header: t('member.columns.username'),
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('Fullname', {
      header: t('member.columns.fullname'),
      cell: info => info.getValue(),
    }),
    // columnHelper.accessor('Bankname', {
    //   header: t('member.columns.bankname'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('Banknumber', {
    //   header: t('member.columns.banknumber'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('Password', {
    //   header: t('columns.username'),
    //   cell: info => info.getValue(),
    // }),

    columnHelper.accessor('Deposit', {
      header: t('member.columns.deposit'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span className={cn(parsedValue <= 0 ? "text-grey-500" : "text-green-500","text-right inline-block") }>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
         columnHelper.accessor('TDeposit', {
      header: t('member.columns.trxdeposit'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span style={{ color: parsedValue <= 0 ? 'grey' : 'blue' }}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
     columnHelper.accessor('Withdraw', {
      header: t('member.columns.withdraw'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span style={{ color: parsedValue < 0 ? 'red' : 'grey' }}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
     columnHelper.accessor('TWithdraw', {
      header: t('member.columns.trxwithdraw'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span style={{ color: parsedValue <= 0 ? 'grey' : 'blue' }}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
      columnHelper.accessor('Crdb', {
      header: t('member.columns.drdb'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span style={{ color: parsedValue < 0 ? 'red' : 'green' }}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
    columnHelper.accessor('SumProamount', {
      header: t('member.columns.promotion'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span style={{ color: parsedValue < 0 ? 'blue' : parsedValue==0?'grey':'blue' }}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
      columnHelper.accessor('Win', {
      header: t('member.columns.win'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
         return (
          <span style={{ color: parsedValue < 0 ? 'red' : parsedValue==0?'grey':'green' }}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
    columnHelper.accessor('Lose', {
      header: t('member.columns.lose'),
      cell: info => {
        const value = info.getValue();
        const parsedValue = parseFloat(value?.toString());
        return (
          <span style={{ color: parsedValue < 0 ? 'red' : parsedValue==0?'grey':'green'}}>
          {formatNumber( parsedValue, 2)}
          </span>
          
         )
      }
    }),
    // columnHelper.accessor('TotalEarnings', {
    //   header: t('member.columns.totalcommission'),
    //   cell: info => {
    //     const value = info.getValue();
    //     return formatNumber(parseFloat(value?.toString()), 2);
    //   }
    // }),
      // columnHelper.accessor('Status', {
      //   header: t('member.columns.status'),
      //   cell: info => {
      //     const value = info.getValue();
      //     return value === 1 ? t('common.active') : value === 0 ? t('common.inactive') : t('common.maintenance');
      //   //  console.log('Raw specificTime value:', value); // For debugging
         
      //   }
      // }),
 
    // columnHelper.accessor('Active', {
    //   header: t('columns.active'),
    //   cell: info => {
    //     const value = info.getValue();
    //     return value === 1 ? t('active') : value === 0 ? t('inactive') : t('maintenance');
    //   }
    // }),
    // columnHelper.accessor('ProStatus', {
    //   header: t('member.columns.prostatus'),
    //   cell: info => info.getValue(),
    // }),
    // columnHelper.accessor('ProActive', {
    //   header: t('member.columns.proactive'),
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
    // {
    //   id: "actions",
    //   enableHiding: false,
    //   cell: ({ row }) => {
    //     const member = row.original as iMember;
    //     return (
    //       <div>
    //         <Button 
    //           variant="ghost" 
    //           onClick={() => openEditPanel(member)}
    //         >
    //           {t('member.edit.title')}
    //         </Button>
    //       </div>
    //     );
    //   },
    // },
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
    setPrefix(member.Prefix)
  };

  const closeEditPanel = () => {
    setEditingGame(null);
    setIsAddingGame(false);
    setShowTable(true);
    setPrefix("")
    setRefreshTrigger((prev:any) => prev + 1);
  };


  function onSubmit(data: z.infer<typeof FormSchema>) {
    toast({
      title: "You submitted the following values:",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(date, null, 2)}</code>
        </pre>
      ),
    })
   
    setIsLoading(true);
   
     GetMemberList(date).then(response=>{
        if(response.Status){
          setGames(response.Data);
        } else {
            // toast({
            //     title: t("common.fetch.error"),
            //     description: t("common.fetch.error_description"),
            //     variant: "destructive",
            //   })
        }
    setIsLoading(false);
    })
   
  }


  if (isLoading) {
    return <div>Loading {t('member.title')}...</div>;
  }

  return (
    <div className="w-full">
     
      {showTable ? (
        <>
         
          <div className="flex items-center justify-between gap-4 py-4">
          <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} >
          <div className="gap-4 flex flex-row">
          <FormField
                control={form.control}
                name="search"
                render={({ field }) => (
                  <FormItem className="">   
            <Input
              placeholder={t('member.columns.search')}
              value={(table.getColumn("Username")?.getFilterValue() as string) ?? ""}
              onChange={(event) =>
                table.getColumn("Username")?.setFilterValue(event.target.value)
              }
              className="max-w"
            />
            </FormItem>
                )} />
            <FormField
                control={form.control}
                name="startdate"
                render={({ field }:any) => (
                  <FormItem className="">   
                      
                      <CalendarDateRangePicker lng={lng}  
                      onChange={(value:any) => { console.log(value); field.onChange(value)} }
                      date={date}
                      setDate={setDate}/>
                   
               
                </FormItem>
                )} />
                 <Button type="submit">{t('button.refresh')}</Button>
        </div>
        </form>
        </Form>
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
          <EditMember 
            memberId={editingGame || 0}
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
