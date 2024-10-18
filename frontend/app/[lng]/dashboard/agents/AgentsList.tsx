'use client';
 
import { AgentListDataTable, GroupedDatabase } from "./data-table"
import { ColumnDef } from "@tanstack/react-table";
 
 

interface AgentsListProps {
  lng: string;
  data: GroupedDatabase[];
  columns: ColumnDef<GroupedDatabase, any>[];
}

export default function AgentsList({ lng, data, columns }: AgentsListProps) {
  // Sort by prefix
  const sortedData = [...data].sort((a, b) => a.prefix.localeCompare(b.prefix));
 

  return <AgentListDataTable lng={lng} columns={columns} data={sortedData} />;
}
