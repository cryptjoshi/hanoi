
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
import { useTranslation } from "@/app/i18n";
import useAuthStore from "@/store/auth";
   
  export default async function AgentAction({ params }: { params: {  prefix: string,lng:string } }) {
    const { prefix,lng } = params;
    
    
    
    const {t} = await useTranslation(lng,'translation');
    return (
      <ContentLayout title={t(`agents.title`)+" "+prefix} >
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href={`/${lng}/`}>{t(`menu.home`)}</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href={`/${lng}/dashboard`}>{t(`menu.dashboard`)}</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href={`/${lng}/dashboard/agents`}>{t(`agents.title`)}</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
          
              <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href={`/${lng}/dashboard/agents`}>{t(`menu.all_agents`)}</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>{t(`menu.edit`)}</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
       <PlaceholderContent>
          <EditAgentSettings params={{
            lng: lng,id:prefix
          }} />
       </PlaceholderContent>
      </ContentLayout>
    );
  }
  