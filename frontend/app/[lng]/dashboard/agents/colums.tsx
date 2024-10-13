"use client"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { CaretSortIcon, DotsHorizontalIcon } from "@radix-ui/react-icons"
import { ColumnDef,FilterFn } from "@tanstack/react-table"
 

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
export const columns: ColumnDef<DatabaseEntry>[] = [
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
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Names
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
                onClick={() => navigator.clipboard.writeText(database.id)}
              >
                Copy database ID
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>View database</DropdownMenuItem>
              <DropdownMenuItem>Delete database</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
  

 
