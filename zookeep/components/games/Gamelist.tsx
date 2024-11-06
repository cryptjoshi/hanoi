'use client'
import Link from "next/link";

import { useEffect, useState } from "react";

import { useTranslation } from "@/app/i18n/client";
import useGameStore from "@/store/gameStore";
import { useRouter } from "next/navigation";
import { GetGameByType } from "@/actions";
import useAuthStore from "@/store/auth";

export default function GameList({ id,lng }: { id:string,lng:string }) {
   //const { gameStatus, fetchGameStatus } = useGameStore()
  const {t} = useTranslation(lng,'translation',undefined);

  const router = useRouter()
  const [gamelist,setGameList] = useState<any[]>([{}])
  const {accessToken} = useAuthStore() 

  if(accessToken){
    useEffect(()=>{
        const fetchgame = async (id:string) =>{
            const response = await GetGameByType(accessToken,id)
            console.log(response)
            if(response.Status){
                setGameList(response.Data)
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


  //useEffect(() => {
   //fetchGameStatus(prefix)
    //console.log('fetchGameStatus',gameStatus)
  //}, [prefix, fetchGameStatus])

  if (!gamelist) {
    return <div>Loading game status...</div>
  }


    return (
        <div className="grid grid-cols-4 gap-2 sm:gap-4 p-4 sm:p-6">
        {gamelist.map((item: any, index: any) => (
           <Link key={index} href={`/${lng}/games/all`}>
         <div  className="flex flex-col items-center">
           <div className="w-10 h-10 sm:w-12 sm:h-12 bg-secondary rounded-lg mb-1"></div>
          
           <span className="text-[10px] sm:text-xs text-center" >{`${item.name|| item.productCode}`}</span>
          
         </div>
         </Link>
       ))}
     </div>
    )
}