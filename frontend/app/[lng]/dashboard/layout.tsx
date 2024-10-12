import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
//import { languages } from '@/app/i18n/settings'


// export async function generateStaticParams() {
//   return languages.map((lng) => ({ lng }))
// }

export default function DemoLayout({
  children,
  // params: {
  //   lng
  // }
}: {
  children: React.ReactNode;
}) {
  return <AdminPanelLayout>{children}</AdminPanelLayout>;
}
