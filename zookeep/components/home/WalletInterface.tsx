'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { Avatar} from "@/components/ui/avatar";
import { Card } from '@/components/ui/card';
import { Bell, Search, Wallet, MessageCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import  type { JSX } from 'react';
import { useTranslation } from '@/app/i18n/client';
import { GetUserInfo,GetPromotion } from '@/actions';
import { formatNumber } from '@/lib/utils';
import useGameStore from '@/store/gameStore';
import useAuthStore from '@/store/auth';
import GameList from './gamelist';
import PromotionList from './promotionlist';


export default function WalletInterface({lng}:{lng:string}): JSX.Element {
  const router = useRouter();
  const [loading, setLoading] = React.useState(true);
  const [balance, setBalance] = React.useState(0);
  const [user, setUser] = React.useState(null);

  const [currency, setCurrency] = React.useState('USD');

  const {prefix,Logout,setPrefix} = useAuthStore();
  const handleSignOut = () => {
    Logout();
    router.push(`/${lng}/login`);
  };
  React.useEffect(() => {
    const fetchBalance = async () => {
      setLoading(true);
      const userLoginStatus = JSON.parse(localStorage.getItem('userLoginStatus') || '{}');
    
      

      if (userLoginStatus.state) {
                if(userLoginStatus.state.isLoggedIn && userLoginStatus.state.accessToken) {
        const user = await GetUserInfo(userLoginStatus.state.accessToken);
     
        if(user.Status){
          setBalance(user.Data.balance);
          setUser(user.Data);
          setCurrency(userLoginStatus.state.customerCurrency);
          setPrefix(user.Data.prefix);
        } else {
          // Redirect to login page if token is null
        router.push(`/${lng}/login`);
        return;
        }
       
     
      } else {
        router.push(`/${lng}/login`);
        return;
        }
      } else {
        router.push(`/${lng}/login`);
        return;
      }
      setLoading(false);
    };

    fetchBalance();
   
  }, [lng, router]);



  const { t } = useTranslation(lng,'home',undefined);

  return loading ? <div>Loading...</div> : (
    <div className="max-w-md mx-auto bg-background text-foreground min-h-screen flex flex-col">
      <div className="p-4 sm:p-6">
       <div className="grid grid-cols-2 gap-4">
        <div>
          <p className="text-xs sm:text-sm text-muted-foreground">{t('balance')}</p>
          <h2 className="text-xl sm:text-2xl font-bold mt-1">{formatNumber(parseFloat(balance))}</h2>
          <p className="text-xs sm:text-sm text-muted-foreground mt-1">≈${formatNumber(parseFloat(balance))} {currency}</p>
        </div>
        <div>
          <p className="text-xs sm:text-sm text-muted-foreground">{user?.fullname}</p>
          <p className="text-xs sm:text-sm text-muted-foreground">{user?.username}</p>
          <p className="text-xs sm:text-sm text-muted-foreground">{user?.bankname}</p>
        </div>
      </div>
      <div className="flex space-x-2 sm:space-x-4 mt-4">
          <Button className="flex-1 bg-yellow-400 text-black hover:bg-yellow-500 text-sm sm:text-base py-2 sm:py-3">{t('deposit')}</Button>
          <Button className="flex-1 text-sm sm:text-base py-2 sm:py-3" variant="outline">{t('withdraw')}</Button>
         </div>
     </div>
    
 
 
      <GameList prefix={prefix} lng={lng} />
 
      <PromotionList prefix={prefix} lng={lng} />

    
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

   
     <div className="mt-auto fixed bottom-0 left-0 right-0 border-t flex justify-between p-2 sm:p-3 bg-background/95 shadow backdrop-blur supports-[backdrop-filter]:bg-background/60 dark:shadow-secondary">
      
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
      </div>
   </div>
  );
};

 
