'use client'
import React, { useState } from 'react';
import { redirect, useRouter } from 'next/navigation';
import { Avatar} from "@/components/ui/avatar";
import { Card } from '@/components/ui/card';
import { Bell, Search, Wallet, MessageCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import  type { JSX } from 'react';
import { useTranslation } from '@/app/i18n/client';
import { GetUserInfo,GetPromotion, UpdateUser, UpdateUserPromotion } from '@/actions';
import { formatNumber } from '@/lib/utils';
import useGameStore from '@/store/gameStore';
import useAuthStore from '@/store/auth';
import GameList from './homegamelist';
import {Promotion,PromotionList} from './promotionlist';
import { useToast } from '@/hooks/use-toast';

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
  const [currency, setCurrency] = React.useState<string>('USD');

  const {prefix,Logout,setPrefix} = useAuthStore();
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
  const userLoginStatus = JSON.parse(localStorage.getItem('userLoginStatus') || '{}');
  const [token, setToken] = useState<string>(userLoginStatus.state?.accessToken);
  const accpetedPromotion = (promotion:Promotion) =>{

    const accepted = async (promotion:Promotion) => {
     
      //console.log(token,prefix)
      //console.log(promotion)
      
      if(token && prefix!=""){
     
        const res = await UpdateUserPromotion(token,{"prefix":prefix,"pro_status":promotion.ID})
     
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
  } else {
    toast({
      title: t('common.unsuccess'),
      description: t('common.loginFirst'),
      variant: "destructive",
    })
    router.push(`/${lng}/login`);
  }
  }
  
  accepted(promotion)
 // acceptPromotion()
  }

  React.useEffect(() => {
    const fetchBalance = async () => {

      try {
      setLoading(true);
      const userLoginStatus = JSON.parse(localStorage.getItem('userLoginStatus') || '{}');
    
      
 
    
                if(userLoginStatus.state.isLoggedIn && userLoginStatus.state.accessToken) {
        const user:any = await GetUserInfo(userLoginStatus.state.accessToken);
        //console.log(user)
        if(user.Status){
          setBalance(user.Data.balance);
       
          setUser(user.Data);
          setCurrency(userLoginStatus.state.customerCurrency);
          setPrefix(user.Data.prefix);
           
        }  
     
      } else {
        router.push(`/${lng}/login`);
        return;
        }
     
      } catch (error) {
       // router.push(`/${lng}/login`);
       console.log(error)
      }
     
    };


    fetchBalance();
    setLoading(false);
   
  }, [lng, router]);

  React.useEffect(() => {
    const fetchPromotion = async (prefix: string) => {
      setIsLoading(true);
      if(token){
      const promotion = await GetPromotion(token);
      //console.log(promotion)
      if (promotion.Status) {
        // กรองโปรโมชั่นที่มี ID ไม่ตรงกับ user.pro_status
         
        //console.log(user)
        const filtered = promotion.Data.Promotions.filter((promo:any) => 
          //{
            //console.log(promo.ID.toString(), user?.pro_status?.toString())
            (1*promo.ID)-(1*user?.pro_status) != 0 
          //}
        );
        
        setPromotions(promotion.Data.Promotions);
        setFilteredPromotions(promotion.Data.Promotions);
        // ถ้า filtered เป็น array ว่าง ให้สร้าง promotion เริ่มต้น
         
        // if (filtered.length === 0 ) {
        //   setFilteredPromotions([{
        //     ID: 'default',
        //     name: t('defaultPromotion'),
        //     description: t('noAvailablePromotions'),
        //     image: '/path/to/default/image.jpg',
        //     disableAccept: true, // เพิ่มคุณสมบัตินี้เพื่อ disable ปุ่ม Accept
        //     // เพิ่ม properties อื่นๆ ตามที่จำเป็นสำหรับ PromotionList component
        //   }]);
        // } else {
        //     console.log(user?.pro_status)
        //     if(user?.pro_status=="" || user?.pro_status==null || user?.pro_status=="0")
        //     setFilteredPromotions(filtered);
        //     else
        //     setFilteredPromotions([{
        //       ID: 'default',
        //       name: t('defaultPromotion'),
        //       description: t('noAvailablePromotions'),
        //       image: '/path/to/default/image.jpg',
        //       disableAccept: true, // เพิ่มคุณสมบัตินี้เพื่อ disable ปุ่ม Accept
        //       // เพิ่ม properties อื่นๆ ตามที่จำเป็นสำหรับ PromotionList component
        //     }]);
        
        // }
      } else {
        
        //router.push(`/${lng}/login`)toast({
      toast({title: t('unsuccess'),
        description: promotion.Message,
        variant: "destructive",
      });
        return;
      }
    } else {
      router.push(`/${lng}/login`);
      return;
    }
  }
    fetchPromotion(prefix);
    setIsLoading(false);
  }, [prefix, user?.pro_status, t,filteredPromotions])

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
            <p className="text-xs sm:text-sm text-muted-foreground">
    
              {selectedPromotion 
                ? selectedPromotion.name // Display selected promotion name if available
                :  promotions.find(promo => promo.ID.toString() == user?.pro_status && promo.status==1)?.name || t('noPromotion')  // Changed ID to id and added fallback text
              }   
            </p>
          </div>
        </div>
      </div>
      <div className="flex space-x-2 sm:space-x-4 mt-4">
      <Button className="flex-1 bg-yellow-400 text-black hover:bg-yellow-500 text-sm sm:text-base py-2 sm:py-3" onClick={() => router.push(`/${lng}/transaction/deposit`)}>{t('deposit')}</Button>
      <Button className="flex-1 text-sm sm:text-base py-2 sm:py-3" variant="outline" onClick={() => router.push(`/${lng}/transaction/withdraw`)}>{t('withdraw')}</Button>         </div>
     </div>
    
 
     
      <GameList prefix={prefix} includegames={user?.includegames} excludegames={user?.excludegames} lng={lng} />
 
      {isLoading ? <div>Loading...</div> : (
      <PromotionList 
        prefix={prefix} 
        lng={lng} 
        promotions={filteredPromotions} 
        onSelectPromotion={accpetedPromotion} 
      />
      )}

    
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
    {/* <Footer lng={lng} /> */}
   
     {/* <div className="mt-auto fixed bottom-0 left-0 right-0 border-t flex justify-between p-2 sm:p-3 bg-background/95 shadow backdrop-blur supports-[backdrop-filter]:bg-background/60 dark:shadow-secondary">
      
        {['Home', 'Deposit', 'Withdraw', 'History', 'sign_out'].map((item, index) => (
          <Button 
            key={index} 
            variant="ghost" 
            className="flex-col py-1 px-2 sm:py-2 sm:px-3"
            onClick={item === 'sign_out' ? handleSignOut : undefined}
          >
            <span className="text-[10px] sm:text-xs mt-1">{t(`menu.${item.toLowerCase()}`)}</span>
          </Button>
        ))}
      </div> */}
   </div>
  );
};

 
