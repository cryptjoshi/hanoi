'use client'
import { useTranslation } from "@/app/i18n/client"
import { Button } from "../ui/button"
import { Card } from "../ui/card"
import { Avatar } from "../ui/avatar"
import useAuthStore from "@/store/auth"
import { useRouter } from "next/navigation"
import { getSession, Logout } from "@/actions"
import {useState,useEffect} from 'react'
interface routeFoot {
  label:string
  route:string
}
//{['Home', 'Deposit', 'Withdraw', 'History', 'sign_out']
const footRouter:routeFoot[] = [{label:"home",route:"home"},{label:"deposit",route:"transaction\\deposit"},{label:"withdraw",route:"transaction\\withdraw"},{label:"history",route:"history"},{label:"sign_out",route:"signout"}]


function Footer() {

   
    const [lng,setLng] = useState("")
    //const {t} = useTranslation(lng,"home",undefined)
   // const { t } = useTranslation(lng,'home',undefined);
    
   useEffect(()=>{
     const fetchSession = async ()=>{
      const session = await getSession()
      setLng(session.lng)
    }
    fetchSession()
   },[lng])

    //const {Logout,lng} = useAuthStore();

    const { t } = useTranslation(lng,'home',undefined);
 
    const router = useRouter();
    const handleSignOut = () => {
        //console.log(lng)
        Logout();
       // router.push(`/${lng}/login`);
      };
return (
    <>
    <div className="mt-auto fixed bottom-0 left-0 right-0 border-t flex justify-between p-2 sm:p-3 bg-background/95 shadow backdrop-blur supports-[backdrop-filter]:bg-background/60 dark:shadow-secondary">
      
      {footRouter.map((item, index) => (
        <Button 
          key={index} 
          variant="ghost" 
          className="flex-col py-1 px-2 sm:py-2 sm:px-3"
          onClick={item.label === 'sign_out' ? handleSignOut : () => router.push(`/${lng}/${item.route}`)}> 
            
           <span className="text-[10px] sm:text-xs mt-1">{t(`menu.${item.label}`)}</span> 
        </Button>
      ))}
    </div>
</>

)
}

export default Footer