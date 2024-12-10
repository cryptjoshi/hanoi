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

import { useTranslation } from "@/app/i18n";
import EditMemberSettings from "@/components/member/page";


export default async function MembersPage({ params }: { params: {  prefix: string,lng:string } }){
  const { prefix,lng } = params;
    
    console.log(lng)
    
  const {t} = await useTranslation(lng,'translation');
  return (
    <ContentLayout title="Users">
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
          <Link href={`/${lng}/dashboard/members`}>{t(`member.title`)}</Link>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <PlaceholderContent>
          <EditMemberSettings params={{
            lng: lng,id:prefix
          }} />
       </PlaceholderContent>
    </ContentLayout>
  );
}
