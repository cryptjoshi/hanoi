import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import useAuthStore from "@/store/auth";

export default function DashboardLayout({
  children,
  params: { lng }
}: {
  children: React.ReactNode;
  params: { lng: string };
}) {
  return <AdminPanelLayout lng={lng}>{children}</AdminPanelLayout>;
}

//import { languages } from '@/app/i18n/settings'


// export async function generateStaticParams() {
//   return languages.map((lng) => ({ lng }))
// }

// export default function DemoLayout({
//   children,
//   params: { lng }
// }: {
//   children: React.ReactNode;
//   params: { lng:any };
// }) {
//   console.log(JSON.stringify(lng))
//   return <AdminPanelLayout lng={lng}>{children}</AdminPanelLayout>;
// }
