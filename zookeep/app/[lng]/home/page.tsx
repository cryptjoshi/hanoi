"use client";
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
import WalletInterface from "@/components/home/WalletInterface";
import HomePage from "@/components/home/home";
//import { useTranslation } from '@/app/i18n'

export default  function DashboardPage({ params: { lng } }) {
  //const { t } = await useTranslation(lng)
  const { t } =  useTranslation(lng,'translation' ,'menu');
  const sidebar = useStore(useSidebar, (x) => x);
  if (!sidebar) return null;
  const { settings, setSettings } = sidebar;
  return (
    <ContentLayout title="Dashboard">
        <HomePage lng={lng}/>
    </ContentLayout>
  );
}
