
import { ContentLayout } from "@/components/admin-panel/content-layout";
import { useSidebar } from "@/hooks/use-sidebar";
import { useStore } from "@/hooks/use-store";
import { useTranslation } from "@/app/i18n";

import HomePage from "@/components/home/home";
import { getSession } from "@/actions";
import { redirect } from "next/navigation";

//import { useTranslation } from '@/app/i18n'

export default async function DashboardPage({ params: { lng } }: { params: { lng: string } }) {
  //const { t } = await useTranslation(lng)
  const { t } = await useTranslation(lng,'translation' ,'menu');
  const session = await getSession();
  //const sidebar = useStore(useSidebar, (x) => x);
  //if (!sidebar) return null;
  //const { settings, setSettings } = sidebar;
//console.log(session)
  if(!session.isLoggedIn)
    redirect(`/${lng}/login`)
  return (
    <ContentLayout title="Dashboard">
      
        <HomePage lng={lng}/>
    </ContentLayout>
  );
}
