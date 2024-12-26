import { ModeToggle } from "@/components/mode-toggle";
import { UserNav } from "@/components/admin-panel/user-nav";
import { SheetMenu } from "@/components/admin-panel/sheet-menu";
import LanguageSwitcher from "../LanguageSwitcher";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {Search,Wallet,MessageCircle,Bell} from 'lucide-react'

interface NavbarProps {
  title: string;
}

export function Navbar({ title }: NavbarProps) {
  return (
    <header className="sticky top-0 z-10 w-full bg-background/95 shadow backdrop-blur supports-[backdrop-filter]:bg-background/60 dark:shadow-secondary">
      {/* <div className="mx-4 sm:mx-8 flex h-14 items-center">
        <div className="flex items-center space-x-4 lg:space-x-0">
          <SheetMenu />
          <h1 className="font-bold">{title}</h1>
        </div>
        <div className="flex flex-1 items-center justify-end">
          <LanguageSwitcher />
          <ModeToggle />
          <UserNav />
        </div>
      </div> */}
       <div className="flex justify-between items-center p-4 sm:p-6">
        <Avatar className="w-8 h-8 sm:w-10 sm:h-10" />
        <div className="flex space-x-2 sm:space-x-4">
          <Button variant="ghost" size="icon"><Search className="h-4 w-4 sm:h-5 sm:w-5" /></Button>
          <Button variant="ghost" size="icon"><Wallet className="h-4 w-4 sm:h-5 sm:w-5" /></Button>
          {/* <Button variant="ghost" size="icon"><MessageCircle className="h-4 w-4 sm:h-5 sm:w-5" /></Button>
          <Button variant="ghost" size="icon"><Bell className="h-4 w-4 sm:h-5 sm:w-5" /></Button> */}
          <ModeToggle />
          {/* <LanguageSwitcher /> */}
        </div>
      </div>
    </header>
  );
}
