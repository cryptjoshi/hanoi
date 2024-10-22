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
 
import { GetDatabaseList } from "@/actions";
 
 
 
 
import { useParams } from 'next/navigation';

// เพิ่ม import สำหรับ Client Component

import { useTranslation } from "@/app/i18n";
import PromotionList from "@/components/promotion/list";

// ปรับ interface ให้ตรงกับข้อมูลที่ได้รับ
interface DatabaseResponse {
  Databases: {
    [prefix: string]: string[];
  };
  Status: boolean;
  Message: string;
}

interface DatabaseEntry {
  name: string;
  prefix: string;
  type: string;
}

 

export default async function PostsPage({ params }: { params: { lng: string } }) {
  const { lng } = params;
  const data: DatabaseResponse = await GetDatabaseList();
  const { t } = await useTranslation(lng, "agents")
  

  // แปลงข้อมูลจาก object เป็น DatabaseEntry[]
  const flattenedData: DatabaseEntry[] = Object.entries(data.Databases).reduce((acc, [prefix, names]) => {
    const entries = names.map(name => {
      const [, type = 'other'] = name.split('_').reverse();
      return { name, prefix, type };
    });
    return [...acc, ...entries];
  }, [] as DatabaseEntry[]);
 
  return (
    <ContentLayout title={t(`menu.agent`)}>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href={`/${lng}`}>{t(`menu.home`)}</Link>
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
            <BreadcrumbPage>{t(`menu.agent`)}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <PlaceholderContent>
        <PromotionList  />
      </PlaceholderContent>
    </ContentLayout>
  );
}
