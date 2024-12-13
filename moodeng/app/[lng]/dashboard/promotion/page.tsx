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
import MemberListDataTable, { DataTableProps } from "@/components/member/list";
import PromotionListDataTable, { GroupedDatabase } from "@/components/promotion/list";


export default async function PromotionPage({ params }: { params: { prefix: string, lng: string } }){
  const { prefix, lng } = params;
  let data: DataTableProps<GroupedDatabase> = {
    data: [],
    columns: [],
    rows: []
  };
  console.log(lng);
  
  const { t } = await useTranslation(lng, 'translation');
  return (
    <ContentLayout title="Users" children={undefined}>
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
      <PlaceholderContent children={undefined}>
        <PromotionListDataTable lng={lng} data={data} prefix={prefix} />
          {/* <EditMemberSettings params={{
            lng: lng,id:prefix
          }} /> */}
       </PlaceholderContent>
    </ContentLayout>
  );
}
