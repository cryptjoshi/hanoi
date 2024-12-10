//import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import LanguageSwitcher from "@/components/LanguageSwitcher";
import { Toaster } from "@/components/ui/toaster";

export default function RegisterLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="relative min-h-screen">
      <div className="absolute top-4 right-4 z-10">
        <LanguageSwitcher />
      </div>
      <main>{children}</main>
    </div>
  );
}
