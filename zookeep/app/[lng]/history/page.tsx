"use client";
 
import { ContentLayout } from "@/components/admin-panel/content-layout";
 
import { useSidebar } from "@/hooks/use-sidebar";
import { useStore } from "@/hooks/use-store";
import { useTranslation } from "@/app/i18n/client";
import { History } from "@/components/history";
 
//import { useTranslation } from '@/app/i18n'

export default   function HistoryPage({ params: { lng } }: { params: { lng: string } }) {
  //const { t } = await useTranslation(lng)
  const { t } =  useTranslation(lng,'translation' ,'menu');
 
  const sidebar = useStore(useSidebar, (x) => x);
  if (!sidebar) return null;
  const { settings, setSettings } = sidebar;
  return (
    <ContentLayout title={t('history')}>
        <History lng={lng}/>
    </ContentLayout>
  );
}
