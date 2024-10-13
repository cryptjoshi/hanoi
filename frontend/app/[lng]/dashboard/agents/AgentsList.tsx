'use client';
import { DataTableList } from "./data-table"
import { ColumnDef } from "@tanstack/react-table";

// Define DatabaseEntry interface to match the data structure
interface DatabaseEntry {
  name: string;
  type: 'dev' | 'prod' | 'other';
}

interface AgentsListProps {
  lng: string;
  data: string[];
  columns: ColumnDef<DatabaseEntry, string>[];
}

export default function AgentsList({ lng, data, columns }: AgentsListProps) {
  // แปลงรายการฐานข้อมูลให้เป็น DatabaseEntry[]
  console.log(data)
  const databaseEntries: DatabaseEntry[] = data.map(name => ({
    name,
    type: name.includes('_dev') || name.includes('development') ? 'dev' :
          name.includes('_prod') || name.includes('production') ? 'prod' : 'other'
  }));

  return <DataTableList columns={columns} data={databaseEntries} />;
}
