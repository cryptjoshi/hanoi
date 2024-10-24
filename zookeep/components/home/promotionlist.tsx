'use client'
import { useEffect,useState } from 'react'
import { useTranslation } from '@/app/i18n/client';
import { GetPromotion } from '@/actions';
import { Card } from '@/components/ui/card';

const GameList = ({ prefix,lng }: { prefix: string,lng:string }) => {
 
const {t} = useTranslation(lng,'translation',undefined);
const [promotion, setPromotion] = useState(null);

  useEffect(() => {
    const fetchPromotion = async (prefix:string) => {
    const promotion = await GetPromotion(prefix);
        if(promotion.Status){
         //   console.log('promotion',promotion.Data)
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
             <h4 className="font-bold text-yellow-400 text-sm sm:text-base">{item.name}</h4>
             <p className="text-green-400 text-xs sm:text-sm">{item.description}</p>
           </div>
           <div className="text-right">
             <span className="text-xs sm:text-sm">{item.end_date}</span>
           </div>
         </div>
       </Card>
       ))}

     </div>
    </>
  )
}

export default GameList