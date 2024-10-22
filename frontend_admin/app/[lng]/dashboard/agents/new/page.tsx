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
import { useTranslation } from "@/app/i18n";

export default async function NewPostPage({ params }: { params: { lng: string } }) {
  const { lng } = params;
  const { t } = await useTranslation(lng, "translation")
  return (
    <ContentLayout title={t(`menu.new_agent`)} >
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
              <Link href={`/${lng}/dashboard/agents`}>{t(`menu.agent`)}</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{t(`menu.new_agent`)}</BreadcrumbPage>
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
