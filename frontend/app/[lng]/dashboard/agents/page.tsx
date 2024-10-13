import Link from "next/link";

import PlaceholderContent from "@/components/demo/placeholder-content";
import { ContentLayout } from "@/components/admin-panel/content-layout";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from "@/components/ui/breadcrumb";
import { DatabaseEntry, columns } from "./colums"
import { DataTableList } from "./data-table"
import { GetDatabaseList } from "@/actions";
import { ColumnDef } from "@tanstack/react-table";
//import { DataTableDemo } from "@/components/agents/lists";
 
 
 
import { useParams } from 'next/navigation';

// เพิ่ม import สำหรับ Client Component
import AgentsList from './AgentsList';

interface DatabaseEntry {
  Databases: string[];
  Status: boolean;
  Message: string;
}

export default async function PostsPage({ params }: { params: { lng: string } }) {
  const { lng } = params;
  const data:DatabaseEntry = await GetDatabaseList();

  return (
    <ContentLayout title="All Agents">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href={`/${lng}`}>Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href={`/${lng}/dashboard`}>Dashboard</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Agents</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <AgentsList lng={lng} data={data.Databases} columns={columns} />
    </ContentLayout>
  );
}
