
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
  import { AccountForm } from "@/app/forms/account/account-form";
  import { ProfileForm } from "@/components/agents/new/profile-register";
  import SettingsProfilePage from "@/components/agents/new/page";
  import SettingsLayout from "@/components/agents/new/layout";
   
  export const viewAgent = ( {params }: { params: { lng: string,id: string } }) => {
    const { lng,id } = params;
    return (
      <ContentLayout title="View Agent" >
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href="/">Home</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href="/dashboard">Dashboard</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href="/agens">Agents</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>New</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
       <PlaceholderContent>
        {/* <SettingsLayout>  */}
          <SettingsProfilePage params={{
            lng: lng,id:id
          }} />
        {/* </SettingsLayout> */}
       </PlaceholderContent>
      </ContentLayout>
    );
  }
  