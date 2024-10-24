'use client'
import { useEffect,useState } from 'react'
import { useTranslation } from '@/app/i18n/client';
<<<<<<< HEAD
import { GetPromotion, UpdateUser } from '@/actions';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useToast } from '@/hooks/use-toast';

=======
import { GetPromotion } from '@/actions';
>>>>>>> 6c7dfb82ae96a678b769c3016b6e256e832fc090

const GameList = ({ prefix,lng }: { prefix: string,lng:string }) => {
 
const {t} = useTranslation(lng,'translation',undefined);
const [promotion, setPromotion] = useState(null);
const { toast } = useToast();
  const handleAccept =  (id:number) => {
    console.log('id',id)
     UpdateUser(prefix,{pro_status:id})
    toast({
      title: t('common.success'),
      description: t('common.promotionAccept'),
    })
    
  }


  useEffect(() => {
    const fetchPromotion = async (prefix:string) => {
    const promotion = await GetPromotion(prefix);
        if(promotion.Status){
          setPromotion(promotion.Data);
        }
    }
    fetchPromotion(prefix);
  }, [prefix])

  if (!promotion) {
    return <div>{t(`games.title`)}</div>
  }

  return (
    <>
   
   <div className="p-4 sm:p-6">
       <h3 className="font-bold text-sm sm:text-base mb-2">{t('latestEvents')}</h3>

       {promotion && promotion.map((item, index) => (
       <Card key={index} className="bg-black text-white p-3 sm:p-4">
         <div className="flex justify-between items-center">
           <div>
             <h4 className="font-bold text-yellow-400 text-sm sm:text-base">{item.title}</h4>
             <p className="text-green-400 text-xs sm:text-sm">{item.description}</p>
           </div>
           <div className="text-right">
             <span className="text-xs sm:text-sm">{item.end_date}</span>
           </div>
         </div>
         <CardContent>
         <p className="text-xs sm:text-sm">{item.description}</p>
         </CardContent>
       <CardFooter>
        <Button className="bg-yellow-400 text-black hover:bg-yellow-500" onClick={() => handleAccept(item)}>{t('common.accept')}</Button>
       </CardFooter>
       </Card>
       ))}

     </div>
    </>
  )
}

export default GameList