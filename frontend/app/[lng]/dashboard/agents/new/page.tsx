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
 

export default function NewPostPage({ params }: { params: { lng: string } }) {
  const { lng } = params;
  return (
    <ContentLayout title="New Agent" >
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
            <BreadcrumbPage>New</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
     <PlaceholderContent>
      {/* <SettingsLayout>  */}
        <SettingsProfilePage params={{
          lng: lng
        }} />
      {/* </SettingsLayout> */}
     </PlaceholderContent>
    </ContentLayout>
  );
}
