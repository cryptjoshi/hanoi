'use client'
import { useTranslation } from "@/app/i18n/client"
import { Button } from "../ui/button"
import { Card } from "../ui/card"
import { Avatar } from "../ui/avatar"
import useAuthStore from "@/store/auth"
import { useRouter } from "next/navigation"

function Footer() {

   
    
    //const {t} = useTranslation(lng,"home",undefined)
   // const { t } = useTranslation(lng,'home',undefined);
    const {Logout,lng} = useAuthStore();
    const { t } = useTranslation(lng,'home',undefined);
    const router = useRouter();
    const handleSignOut = () => {
        Logout();
        router.push(`/${lng}/login`);
      };
return (
    <>
    <div className="mt-auto fixed bottom-0 left-0 right-0 border-t flex justify-between p-2 sm:p-3 bg-background/95 shadow backdrop-blur supports-[backdrop-filter]:bg-background/60 dark:shadow-secondary">
      
      {['Home', 'Deposit', 'Withdraw', 'History', 'sign_out'].map((item, index) => (
        <Button 
          key={index} 
          variant="ghost" 
          className="flex-col py-1 px-2 sm:py-2 sm:px-3"
          onClick={item === 'sign_out' ? handleSignOut : () => router.push(`/${lng}/${item.toLowerCase()}`)}> 
            
           <span className="text-[10px] sm:text-xs mt-1">{t(`menu.${item.toLowerCase()}`)}</span> 
        </Button>
      ))}
    </div>
</>

)
}

export default Footer