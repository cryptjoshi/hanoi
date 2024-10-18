"use client"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { CaretSortIcon, DotsHorizontalIcon } from "@radix-ui/react-icons"
import { ColumnDef,FilterFn } from "@tanstack/react-table"
import { viewAgent } from './agentActions';
import { useRouter } from 'next/navigation';
import { usePathname } from 'next/navigation';
import { Link } from "lucide-react"
import { useTranslation } from '@/app/i18n/client'
import useAuthStore from "@/store/auth"

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type DatabaseEntry = {
  id: string
  prefix: string
  names: string[]
}
const customArrayFilter: FilterFn<DatabaseEntry> = (row, columnId, filterValue) => {
  const names = row.getValue(columnId) as string[];
  return names.some(name => name.toLowerCase().includes(filterValue.toLowerCase()));
};
export const  columns: ColumnDef<DatabaseEntry>[] =[
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "prefix",
      header: "Prefix",
      cell: ({ row }) => (
        <div className="capitalize">{row.getValue("prefix")}</div>
      ),
    },
    {
      accessorKey: "names",
      header: ({ column }) => {
        const pathname = usePathname();
        const pathParts = pathname.split('/')
        const lng = pathParts[1] // Extract lng from the current path
        const {t} = useTranslation(lng,'translation',{keyPrefix:'promotion'})
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            {t("names")}
            <CaretSortIcon className="ml-2 h-4 w-4" />
          </Button>
        )
      },
      cell: ({ row }) => <div className="lowercase">{(row.getValue("names") as string[]).join(", ")}</div>,
      filterFn: customArrayFilter
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        const database = row.original
        const router = useRouter();
        const pathname = usePathname();
        const pathParts = pathname.split('/')
        const lng = pathParts[1] // Extract lng from the current path
        const {t} = useTranslation(lng,'translation',{keyPrefix:'promotion'})
        const handleViewAgent = () => {
          
          router.push(`/${lng}/dashboard/agents/${database.prefix}/`)
        }

        return (
          <Button  variant="ghost" onClick={handleViewAgent}>{t("viewagent")}</Button>
          // <DropdownMenu>
          //   <DropdownMenuTrigger>Actions</DropdownMenuTrigger>
          //   <DropdownMenuContent align="end">
          //     <DropdownMenuItem onClick={handleViewAgent}>View Agent</DropdownMenuItem>
          //     <DropdownMenuItem>Drop Agent</DropdownMenuItem>
          //   </DropdownMenuContent>
          // </DropdownMenu>
        )
      },
    },
  ]

 
