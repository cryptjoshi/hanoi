'use client'
import Link from "next/link";

import { useEffect, useState } from "react";

import { useTranslation } from "@/app/i18n/client";
import useGameStore from "@/store/gameStore";
import { useRouter } from "next/navigation";
import { GetGameByProvide, getGameUrl } from "@/actions";
import useAuthStore from "@/store/auth";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"


export default function GamesList({ id,lng }: { id:string,lng:string }) {
   //const { gameStatus, fetchGameStatus } = useGameStore()
  const {t} = useTranslation(lng,'translation',undefined);
  const [gameid,setGameId] =useState(id)
  const router = useRouter()
  const [gamelist,setGameList] = useState<any[]>([{}])
 
  const {accessToken,user,customerCurrency} = useAuthStore()

  if(accessToken){
    useEffect(()=>{
        const fetchgame = async (id:string) =>{
            const response = await GetGameByProvide(accessToken,id)
           //console.log(response)
             if(response.Status){
                 setGameList(response.Data?.games)
             } else {

             }
        }
        fetchgame(id)
    },[])
  } else {
    router.push(`/${lng}/login`)
  }
//   const playgame = (ID:string) =>{
//       router.push(`/${lng}/games/${ID}`)
//   }
const openInNewTab = (url:string) => {
  const newWindow = window.open(url, '_blank', 'noopener,noreferrer')
  if (newWindow) newWindow.opener = null
}
   const playgame = (code:string) =>{
    getGameUrl(accessToken || "",code,user.username,customerCurrency || "USD").then((gameurl)=>{
    if(gameurl.Status){
        //console.log(gameurl.Data.url)
     //   window.open(gameurl.Data.url,'_blank','noopener,noreferrer')
       openInNewTab(gameurl.Data.url)
    }   
   })
  }

  //useEffect(() => {
   //fetchGameStatus(prefix)
    //console.log('fetchGameStatus',gameStatus)
  //}, [prefix, fetchGameStatus])
  //href={`/${lng}/games/list/${id}/play/${item.code}`}

  if (!gamelist) {
    return <div>Loading game status...</div>
  }


    return (
        <div className="grid grid-cols-4 gap-2 sm:gap-4 p-4 sm:p-6">
        {gamelist.map((item: any, index: any) => (
           <Link key={index} onClick={()=>playgame(item.code)} href="">
         <div  className="flex flex-col items-center">
          
           <Avatar>
            <AvatarImage src={item.img} />
            <AvatarFallback>{item.name}</AvatarFallback>
          </Avatar>
        
          
           <span className="text-[10px] sm:text-xs text-center" >{`${item.name|| item.productCode}`}</span>
          
         </div>
         </Link>
       ))}
     </div>
    )
}