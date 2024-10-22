
  // Implement the viewAgent logic here
  import Link from "next/link";
  import { Separator } from "@/components/ui/separator"
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
  import EditAgentSettings from "@/components/agents/edit/page";
   
  export default function AccountPage({ params }: { params: { lng: string, prefix: string } }) {
    const { lng,prefix } = params;
    return (
      <ContentLayout title="View Agent" >
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href={`/${lng}/`}>Home</Link>
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
              <BreadcrumbLink asChild>
                <Link href={`/${lng}/dashboard/agents`}>Agents</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>Edit</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
       <PlaceholderContent>
      
        <div>
          <h3>Account</h3>
          {/* Account content */}
        </div>
          {/* <EditAgentSettings params={{
            lng: lng,id:prefix
          }} /> */}
       </PlaceholderContent>
      </ContentLayout>
    );
  }
  