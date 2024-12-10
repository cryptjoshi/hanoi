 
import Link from "next/link";
import { ContentLayout } from "@/components/admin-panel/content-layout";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from "@/components/ui/breadcrumb";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from "@/components/ui/tooltip";
import { useSidebar } from "@/hooks/use-sidebar";
import { useStore } from "@/hooks/use-store";
import { useTranslation } from "@/app/i18n/client";
//import { useTranslation } from '@/app/i18n'
//import { redirect } from 'next/navigation';

import OverViewPage from './_components/overview';

export const metadata = {
  title: 'Dashboard : Overview'
};

export default  function DashboardPage({ params: { lng } }) {
  //const { t } = await useTranslation(lng)
  const { t } =  useTranslation(lng,'translation' ,'menu');
  const sidebar = useStore(useSidebar, (x) => x);
  if (!sidebar) return null;
  const { settings, setSettings } = sidebar;
  
 
    return (
    <ContentLayout title="Dashboard">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href={`/${lng}`}>{t(`menu.home`)}</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{t(`menu.dashboard`)}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <TooltipProvider>
        <div className="flex gap-6 mt-6">
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="flex items-center space-x-2">
                <Switch
                  id="is-hover-open"
                  onCheckedChange={(x) => setSettings({ isHoverOpen: x })}
                  checked={settings.isHoverOpen}
                />
                <Label htmlFor="is-hover-open">{t(`menu.hover_open`)}</Label>
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>{t(`menu.hover_open`)}</p>
            </TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="flex items-center space-x-2">
                <Switch
                  id="disable-sidebar"
                  onCheckedChange={(x) => setSettings({ disabled: x })}
                  checked={settings.disabled}
                />
                <Label htmlFor="disable-sidebar">{t(`menu.disable_sidebar`)}</Label>
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>{t(`menu.disable_sidebar`)}</p>
            </TooltipContent>
          </Tooltip>
        </div> 
      </TooltipProvider>
      <OverViewPage />;

    </ContentLayout>
   );
}
