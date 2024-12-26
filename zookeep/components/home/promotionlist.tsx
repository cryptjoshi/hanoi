'use client'
import { useEffect,useState } from 'react'
import { useTranslation } from '@/app/i18n/client';
import { useRouter } from 'next/navigation';
import { GetPromotion, UpdateUser } from '@/actions';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useToast } from '@/hooks/use-toast';
import useAuthStore from '@/store/auth'
 

import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"

import { getSession } from '@/actions';


export interface Promotion {
  ID: string;
  name: string;
  description: string;
  image: string;
  disableAccept?: boolean;
}

interface PromotionListProps {

  lng: string;
  promotions: Promotion[];
  onSelectPromotion: (promotion: Promotion) => void;
}

export const PromotionList = ({ lng, promotions, onSelectPromotion }: PromotionListProps) => {
  const router = useRouter();
  const { t } = useTranslation(lng, 'translation', undefined);
  const { toast } = useToast();
  //const { accessToken } = useAuthStore()

//   const handleAccept =  (item: Promotion) => {
//     const acceptPromotion = async () => {
//     // const session = await getSession()
//     // if(session.token && session.prefix!=""){

//     const res = await UpdateUser({"pro_status":item.ID.toString()})
  
//     if(res.Status){
//     toast({
//       title: t('common.success'),
//       description: t('common.promotionAccept'),
//       variant: "default",
//     })
//     onSelectPromotion(item);
//   } else {
//     toast({
//       title: t('common.unsuccess'),
//       description: res.Message,
//       variant: "destructive",
//     })
//    // router.push(`/${lng}/login`);
//   }
// // } else {
// //   toast({
// //     title: t('common.unsuccess'),
// //     description: t('common.loginFirst'),
// //     variant: "destructive",
// //   })
// //  // router.push(`/${lng}/login`);
// // }
// }
//   acceptPromotion()
// }


  if (!promotions || promotions.length === 0) {
    return <div>{t(`games.title`)}</div>
  }

  return (
    <>
   
   <div className="p-4 sm:p-6">
       <h3 className="font-bold text-sm sm:text-base mb-2">{t('latestEvents')}</h3>
       <Carousel
      opts={{
        align: "start",
      }}
      className="w-full max-w-xl"
    >
      <CarouselContent className={promotions.length === 1 ? "flex justify-center" : ""}>
       {promotions.map((item, index) => (
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
         <CardContent>
         <p className="text-xs sm:text-sm">{item.description}</p>
         </CardContent>
       <CardFooter>
        <Button 
          onClick={() => onSelectPromotion(item)}
          disabled={item.disableAccept}
        >
          {t('accept')}
        </Button>
       </CardFooter>
       </Card>
       ))}
</CarouselContent>
</Carousel>
     </div>
    </>
  )
}

 
