import { useTranslation } from "@/app/i18n";
import { ProfileEdit } from "@/components/agents/edit/profile-edit";
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AccountForm } from "../../../../components/settings/account/account-form";
import { ContentLayout } from "@/components/admin-panel/content-layout";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import Link from "next/link";
import PlaceholderContent from "@/components/demo/placeholder-content";
import EditSettings from "@/components/settings/page";
//import { ProfileForm } from "@/components/agents/new/profile-register"

export default async function SettingsProfilePage({ params }: { params: { lng: string } }) {
  const { lng } = params;
  const { t } = await useTranslation(lng, "translation")
  return (
    <ContentLayout title={t('settings.title')} >
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
          <BreadcrumbLink asChild>
            <Link href={`/${lng}/dashboard/settings`}>{t(`menu.settings`)}</Link>
          </BreadcrumbLink>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>
   <PlaceholderContent>
    
   <EditSettings params={{ lng, id: '' }} />
   </PlaceholderContent>
  </ContentLayout>
   
  )
}
