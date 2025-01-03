import { Navbar } from "@/components/admin-panel/navbar";
import { cn } from "@/lib/utils";

interface ContentLayoutProps {
  title: string;
  children: React.ReactNode;
}

export function ContentLayout({ title, children }: ContentLayoutProps) {
  return (
    <div>
      <Navbar title={title} />
      <div   
       className={cn("container pt-8 pb-8 px-4 sm:px-8",
        "min-h-[calc(100vh)]",
      )}
      >{children}</div>
    </div>
  );
}
