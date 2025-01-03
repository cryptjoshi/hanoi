'use client'
import React, { useState } from 'react';
import { redirect, useRouter } from 'next/navigation';
import { Avatar} from "@/components/ui/avatar";
import { Card } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

import { Bell, Search, Wallet, MessageCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import  type { JSX } from 'react';
import { useTranslation } from '@/app/i18n/client';
import { GetUserInfo,GetPromotion, UpdateUser, UpdateUserPromotion, GetUserPromotion, CancelPromotion } from '@/actions';
import { formatNumber } from '@/lib/utils';
import useGameStore from '@/store/gameStore';
import useAuthStore from '@/store/auth';
import GameList from './homegamelist';
import {Promotion,PromotionList} from './promotionlist';
import { useToast } from '@/hooks/use-toast';
import { getSession } from '@/actions';
import { cn } from "@/lib/utils"
type PromotionFilter = {
  ID: string
  name: string
  description:string
  image:string
  disableAccept: boolean
}

export default function HomePage({lng}:{lng:string}): JSX.Element {
  const router = useRouter();
  const [loading, setLoading] = React.useState<boolean>(true);
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [balance, setBalance] = React.useState<number>(0);
  const [user, setUser] = React.useState<any>(null);
  const [promotions, setPromotions] = React.useState<any[]>([]);
  const [userPro,setUserPro] = React.useState<any>()
  const [currency, setCurrency] = React.useState<string>('USD');
  const [isBlinking, setIsBlinking] = useState(false);
  const [open,setOpen] = useState(false)
  //const session = getSession()
  //const {prefix,Logout,setPrefix} = useAuthStore();
  // const handleSignOut = () => {
  //   Logout();
  //   router.push(`/${lng}/login`);
  // };
  const [selectedPromotion, setSelectedPromotion] = React.useState<Promotion | null>(null);

  // เพิ่ม state สำหรับ filteredPromotions
  const [filteredPromotions, setFilteredPromotions] = React.useState<PromotionFilter[]>([]);
  const { t } = useTranslation(lng,'home',undefined);

  const { toast } = useToast();
 // const { accessToken } = useAuthStore()
 // const userLoginStatus = JSON.parse(localStorage.getItem('userLoginStatus') || '{}');
  //const [token, setToken] = useState<string>(userLoginStatus.state?.accessToken);




  const accpetedPromotion = (promotion:Promotion) =>{

    const accepted = async (promotion:Promotion) => {
     
      const session = await getSession()
      
    
     
      const res = await UpdateUserPromotion({"prefix":session.prefix,"pro_status":promotion.ID.toString()})
     
      if(res.Status){
      toast({
        title: t('common.success'),
        description: t('common.promotionAccept'),
        variant: "default",
      })

      //onSelectPromotion(item);
      setSelectedPromotion(promotion)
    } else {
      toast({
        title: t('common.unsuccess'),
        description: res.Message,
        variant: "destructive",
      })
     // router.push(`/${lng}/login`);
    }
 
  }
  
  accepted(promotion)
 // acceptPromotion()
  }

  React.useEffect(() => {
    const fetchBalance = async () => {
      const session = await getSession();
      try {
        const user: any = await GetUserInfo();
        console.log(user)
        if (user.Status) {
          setBalance(user.Data.balance);
          setUser(user.Data);
          session.customerCurrency = user.Data.currency;
          setCurrency(session.customerCurrency);
        } else {
          toast({
            title: t('unsuccess'),
            description: user.Message,
            variant: "destructive",
          });
        }
      } catch (error) {
        console.log(error);
      }
    };

    const fetchPromotion = async () => {
      try {
        const user_pro = await GetUserPromotion();
        console.log(user_pro)
        if (user_pro.Status) {
          setUserPro(user_pro.Data);
          setIsBlinking(parseInt(user_pro.Data.status) > 0);
        }

        const promotion = await GetPromotion();
    
        if (promotion.Status) {
          if (promotion.Data.length > 0) {
            setPromotions(promotion.Data[0]);
            setFilteredPromotions(promotion.Data[0]);
          }
        } else {
          toast({
            title: t('unsuccess'),
            description: promotion.Message,
            variant: "destructive",
          });
        }
      } catch (error) {
        console.log(error);
      }
    };
    fetchBalance().then(()=>{
      fetchPromotion();
    })
   
   
   
    setLoading(false);
  }, [lng, t]);

  const handleCancelPromotion = () => {
    const cancelPro = async ()=>{
      const response = await CancelPromotion()
      if(response.Status){
        toast({
          title: t('common.success'),
          description: t('common.promotionAccept'),
          variant: "default",
        })
        setOpen(false);
        setUserPro(null);
        window.location.reload();
      } else {
        toast({
          title: t('common.unsuccess'),
          description: response.Message,
          variant: "destructive",
        })
      }
    }
    // เพิ่มโค้ดสำหรับการยกเลิกโปรโมชั่นที่นี่
    //console.log("โปรโมชั่นถูกยกเลิก"); // ตัวอย่างการแสดงข้อความใน console
    // ปิด Dialog หลังจากยกเลิก
    cancelPro()
    
    
};
  
  return loading ? <div>Loading...</div> : (
    <div className="max-w-md mx-auto bg-background text-foreground min-h-screen flex flex-col">
      <div className="p-4 sm:p-6">
       <div className="grid grid-cols-2 gap-4">
        <div>
          <p className="text-xs sm:text-sm text-muted-foreground">{t('balance')}</p>
          <h2 className="text-xl sm:text-2xl font-bold mt-1">{formatNumber(balance)}</h2>
          <p className="text-xs sm:text-sm text-muted-foreground mt-1">≈${formatNumber(balance)} {currency}</p>
        </div>
        <div>
        <p className="text-xs sm:text-sm text-muted-foreground"> ref: {user?.referredby}</p>
          <p className="text-xs sm:text-sm text-muted-foreground">{user?.fullname}</p>
          <p className="text-xs sm:text-sm text-muted-foreground">{user?.username}</p>
          <p className="text-xs sm:text-sm text-muted-foreground">{user?.bankname}</p>
          <div className="mt-2">
            <p className="text-xs sm:text-sm font-semibold">{t('promotionStatus')}:</p>
            <div className="flex items-center p-2 gap-2 bg-background/95 shadow backdrop-blur supports-[backdrop-filter]:bg-background/60 dark:shadow-secondary">

             <span className={cn(
                "transition-opacity",
                parseInt(userPro?.status) === 1 ? "text-green-500" : "text-red-500",
                isBlinking && "animate-pulse"
              )}>
              
              {selectedPromotion 
                ? selectedPromotion.name // Display selected promotion name if available
                : userPro?.status != "2" ? userPro?.Name : t('noPromotion')  // Changed ID to id and added fallback text
              }   
           </span>
          { selectedPromotion  ||  userPro?.status == "0" ? (
           <Button 
              className="bg-red-500 text-white hover:bg-red-600 transition duration-200 text-xs" 
              onClick={() => setOpen(true)} // เปิด Dialog
           >
              <span className="text-lg">x</span> {/* สัญลักษณ์กากบาท */}
           </Button>
          ):<></>
          }
          </div>
          </div>
        </div>
      </div>
      <div className="flex space-x-2 sm:space-x-4 mt-4">
      <Button className="flex-1 bg-yellow-400 text-black hover:bg-yellow-500 text-sm sm:text-base py-2 sm:py-3" onClick={() => router.push(`/${lng}/transaction/deposit`)}>{t('deposit')}</Button>
      <Button className="flex-1 text-sm sm:text-base py-2 sm:py-3" variant="outline" onClick={() => router.push(`/${lng}/transaction/withdraw`)}>{t('withdraw')}</Button>         </div>
     </div>
    
      <GameList includegames={user?.includegames} excludegames={user?.excludegames} lng={lng} />
 
      <PromotionList 
        lng={lng} 
        promotions={user?.promotionlog}
        onSelectPromotion={accpetedPromotion} 
      />
     
     <div className="p-4 sm:p-6">
       <div className="flex justify-between items-center mb-2">
         <h3 className="font-bold text-sm sm:text-base">{t('lastgameplay')}</h3>
         <Button variant="ghost" size="sm" className="text-xs sm:text-sm">{t('viewMore')}</Button>
       </div>
       <div className="flex space-x-2 sm:space-x-4 overflow-x-auto pb-2">
         {['LTGao', 'XTUSER-', '稳稳的'].map((trader, index) => (
           <Card key={index} className="flex-shrink-0 w-28 sm:w-32 p-2 sm:p-3">
             <div className="text-center mb-2">
               <Avatar className="mx-auto w-8 h-8 sm:w-10 sm:h-10" />
               <p className="text-xs sm:text-sm font-bold mt-1">{trader}</p>
             </div>
             <p className="text-[10px] sm:text-xs text-center">{t('7DROI')}</p>
             <p className="text-xs sm:text-sm text-green-500 text-center font-bold">+{Math.floor(Math.random() * 100)}%</p>
           </Card>
         ))}
       </div>
     </div>
     <div>
           <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
            
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>ยกเลิกโปรโมชั่น</DialogTitle>
                <DialogDescription>
                  คุณแน่ใจหรือไม่ว่าต้องการยกเลิกโปรโมชั่นนี้?
                </DialogDescription>
              </DialogHeader>
              <div className="flex justify-end">
                <Button onClick={handleCancelPromotion} className="mr-2">ยืนยัน</Button>
                <Button onClick={() => setOpen(false)}>ยกเลิก</Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
   </div>
   
  );
};

 
