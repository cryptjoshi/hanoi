import { useEffect, useState } from 'react'
import useGameStore from '@/store/gameStore'
import { useTranslation } from '@/app/i18n/client';
import  { useRouter } from "next/navigation"
import Link from 'next/link';
import useAuthStore from '@/store/auth';
import { User } from 'lucide-react';
import { GameStatus } from '@/lib/zod/gameStatus';


const GameList = ({ lng,includegames,excludegames }: { lng:string,includegames:string,excludegames:string }) => {
  const { gameStatus, fetchGameStatus } = useGameStore()
  const [games,setGames] = useState<any>()
  const {t} = useTranslation(lng,'translation',undefined);
  //const {accessToken} = useAuthStore()
  const router = useRouter()
  const [isLoading,setIsLoading] = useState<boolean>(false)

  const playgame = (ID:string) =>{
      router.push(`/${lng}/games/${ID}`)
  }


  useEffect(() => {
    
    setIsLoading(true)

    //if(accessToken){
      fetchGameStatus()
      
      
   // } else {
  //     router.push(`/${lng}/login`)
  //  }
    setIsLoading(false)
    //console.log('fetchGameStatus',gameStatus)
  }, [fetchGameStatus])

  if (isLoading) {
    return <div>Loading game status...</div>
  }   
    //gameStatus.filter((game)=>includegames.split(",").includes(game.id))
 

    
 
// 
  return (
    <>
   
      <div className="grid grid-cols-4 gap-2 sm:gap-4 p-4 sm:p-6">
        {Array.isArray(gameStatus) && includegames!="" && gameStatus.filter((game) => includegames?.split(",").includes(game.id)).map((item: any, index: any) => (
           <Link key={index} href={`/${lng}/games/${item.name}`}>
         <div  className="flex flex-col items-center">
           <div className="w-10 h-10 sm:w-12 sm:h-12 bg-secondary rounded-lg mb-1"></div>
          
           <span className="text-[10px] sm:text-xs text-center" >{t(`games.${item.name}`)}</span>
          
         </div>
         </Link>
       ))}
        {Array.isArray(gameStatus) && includegames=="" && gameStatus.map((item: any, index: any) => (
           <Link key={index} href={`/${lng}/games/${item.name}`}>
         <div  className="flex flex-col items-center">
           <div className="w-10 h-10 sm:w-12 sm:h-12 bg-secondary rounded-lg mb-1"></div>
          
           <span className="text-[10px] sm:text-xs text-center" >{t(`games.${item.name}`)}</span>
          
         </div>
         </Link>
       ))}
     </div>
    </>
  )
}

export default GameList