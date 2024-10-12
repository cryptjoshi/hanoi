//import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import { Toaster } from "@/components/ui/toaster";

export default function LoginLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (<>{children}   <Toaster /></>);
}
