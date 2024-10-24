import { useEffect } from 'react'
import useGameStore from '@/store/gameStore'
import { useTranslation } from '@/app/i18n/client';

const GameList = ({ prefix,lng }: { prefix: string,lng:string }) => {
  const { gameStatus, fetchGameStatus } = useGameStore()
  const {t} = useTranslation(lng,'translation',undefined);

  useEffect(() => {
    
    fetchGameStatus(prefix)
    console.log('fetchGameStatus',gameStatus)
  }, [prefix, fetchGameStatus])

  if (!gameStatus) {
    return <div>Loading game status...</div>
  }

  return (
    <>
   
      <div className="grid grid-cols-4 gap-2 sm:gap-4 p-4 sm:p-6">
        {gameStatus.map((item: any, index: any) => (
         <div key={index} className="flex flex-col items-center">
           <div className="w-10 h-10 sm:w-12 sm:h-12 bg-secondary rounded-lg mb-1"></div>
           <span className="text-[10px] sm:text-xs text-center">{t(`games.${item.name}`)}</span>
         </div>
       ))}
     </div>
    </>
  )
}

export default GameList