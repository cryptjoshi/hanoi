import { create } from 'zustand'
import { GetGameStatus } from '@/actions'

type GameStatus = {
  // กำหนดโครงสร้างข้อมูลของ GameStatus ตามที่ API ส่งกลับมา
}

interface GameStore {
  gameStatus: GameStatus | null
  fetchGameStatus: (prefix: string,token:string) => Promise<void>
}

const useGameStore = create<GameStore>((set) => ({
  gameStatus: null,
  fetchGameStatus: async (prefix: string,token:string) => {
    set({gameStatus:null})
    try {
      const response = await GetGameStatus(prefix,token)
  
       if(response && response.Status){
        const mappedData = response.Data.map((item: any) => {
            const status = item.status // แปลง JSON string เป็นอ็อบเจ็กต์
            return {
                id: status.id,
                name: status.name
            };
        });
        set({ gameStatus: mappedData })
       }


      //set({ gameStatus: mappedData })
    } catch (error) {
      console.error('Failed to fetch game status:', error)
    }
  },
}))

export default useGameStore