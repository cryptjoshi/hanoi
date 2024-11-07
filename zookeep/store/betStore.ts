// zookeep/store/betStore.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware';

export type BetStore = {
    isWaitingResult: boolean;
    setWaitingResult: (value: boolean) => void;
    betAmount: number;
    setBetAmount: (amount: number) => void;
    betPrice: number;
    setBetPrice: (value: number) => void;
    betPredict: string;
    setBetPredict: (value: string) => void;
    
};
// สร้าง store สำหรับการจัดการสถานะการเดิมพัน
const useBetStore = create(persist((set) => ({
    isWaitingResult: false,
    betAmount: 0,
    betPredict:"",
    betPrice:0,
    isProcessingBet:false,
    // ... other states ...
    setWaitingResult: (value: boolean) => {
        set({ isWaitingResult: value });
        localStorage.setItem('isWaitingResult', JSON.stringify(value)); // เก็บค่าใน localStorage
    },
    setBetAmount: (amount: number) => {
        set({ betAmount: amount });
        localStorage.setItem('betAmount', JSON.stringify(amount)); // เก็บค่าใน localStorage
    },
    setBetPredict: (prediction:string)=>{
        set({prediction:prediction});
        localStorage.setItem('betPredict',prediction)
    },
    setIsProcessingBet:(value: boolean)=>{
        set({isProcessingBet:value})
        localStorage.setItem('isProcessingBet',JSON.stringify(value))
    },
    setBetPrice:(value:number)=>{
        set({betPrice: value})
        localStorage.setItem('betPrice',JSON.stringify(value))
    }
}), {
    name: 'bet-storage', // ชื่อสำหรับ localStorage
}));

export default useBetStore;