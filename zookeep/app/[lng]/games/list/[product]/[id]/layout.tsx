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