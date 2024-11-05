'use client'
import React, { useEffect, useRef, useState,useMemo, memo } from 'react';
import { createChart, ColorType, CrosshairMode } from 'lightweight-charts';
import { GetUserInfo } from '@/actions/index';
import { createTransaction } from '@/actions/index';
import useAuthStore from '@/store/auth';
import { useRouter } from 'next/navigation';
import  ChartComponent  from './ChartComponent';
 
export function Options({lng,data}:{lng:string,data:any}) {
    const initialData = data;
    const [isNewCandle, setIsNewCandle] = useState(false);
    const [users, setUsers] = useState<any>(null);
    const [balance, setBalance] = useState<number>(0);
    const [countdown, setCountdown] = useState(60);
    const [isPredictionDisabled, setIsPredictionDisabled] = useState(false);
    const [betAmount, setBetAmount] = useState<number>(0);
    const [isProcessingBet, setIsProcessingBet] = useState(false);
    const [lastBetResult, setLastBetResult] = useState<'win' | 'lose' | null>(null);
    const [selectedLeverage, setSelectedLeverage] = useState(1); // Default leverage
    const [lastPrediction, setLastPrediction] = useState<'up' | 'down' | null>(null);
    const [leverageAmount, setLeverageAmount] = useState(0); // Add this state
    const [isWaitingResult, setIsWaitingResult] = useState(false);
    const [currentPrediction, setCurrentPrediction] = useState<'up' | 'down' | null>(null);
    const {accessToken} = useAuthStore();
    const router = useRouter();
    const [currentPrice, setCurrentPrice] = useState<number | 0>(0);
    const [predictionStartPrice, setPredictionStartPrice] = useState<number | 0>(0);
    const [betPrice, setBetPrice] = useState<number | null>(null);
    const [closePrice, setClosePrice] = useState<number | null>(null);
    const [priceDirection, setPriceDirection] = useState<'up' | 'down' | null>(null);
    
    const isBettingPeriod = countdown > 45;
    // เพม useEffect สำหรับดึงข้อมูล balance
    useEffect(() => {
        const fetchUserInfo = async () => {
            if (accessToken) {
                try {
                    const userInfo = await GetUserInfo(accessToken);
                    //console.log('User Info received:', userInfo); // Debug log
                    
                    setUsers(userInfo.Data);
                    if (userInfo?.Data?.balance) {
                        setBalance(Number(userInfo.Data.balance));
                    }
                } catch (error) {
                    console.error('Error fetching user info:', error);
                }
            } else {
                router.push(`/${lng}/login`);
            }
        };
        
        fetchUserInfo();
    }, [accessToken]);

    useEffect(() => {
        if (isWaitingResult) {

            setClosePrice(currentPrice);
            if (betPrice) {
                setPriceDirection(currentPrice > betPrice ? 'up' : 'down');
            }
        }
    }, [currentPrice, isWaitingResult, betPrice]);

    const handlePrediction = async (prediction: 'up' | 'down') => {
        if (!isPredictionDisabled && !isProcessingBet && !isWaitingResult) {
            const calculatedBetAmount = 1 * selectedLeverage;
           
            if (calculatedBetAmount <= balance) {
                setIsProcessingBet(true);
                
                try {
                    if (accessToken) {
                        await createTransaction(accessToken, {
                            Status: 100,
                            GameProvide: 'options',
                            MemberName: users.username,
                            TransactionAmount: calculatedBetAmount.toString(),
                            ProductID: 9000,
                            BeforeBalance: balance.toString(),
                            Balance: (balance - calculatedBetAmount).toString(),
                            AfterBalance: (balance - calculatedBetAmount).toString()
                        });
                        
                        if (predictionStartPrice !== null) {
                            setPredictionStartPrice(predictionStartPrice);
                        }
                        console.log('Setting bet price:', currentPrice); // เพิ่ม debug log
                        setIsPredictionDisabled(true);
                        setBetPrice(currentPrice);
                        setCurrentPrediction(prediction);
                        setIsWaitingResult(true);
                        setBalance(prev => prev - calculatedBetAmount);
                        setBetAmount(calculatedBetAmount);
                    } else {
                        router.push(`/${lng}/login`);
                    }

                } catch (error) {
                    console.error('Betting error:', error);
                } finally {
                    setIsProcessingBet(false);
                }
            }
        }
    };

    const handleLeverageClick = (leverage: number) => {
        if (!isBettingPeriod || isProcessingBet || isWaitingResult) return;

        setSelectedLeverage(leverage);
        setLeverageAmount(1 * (leverage));
    };

    const handleClearLeverage = () => {
        if (!isBettingPeriod || isProcessingBet || isWaitingResult) return;
        setSelectedLeverage(1);
        setLeverageAmount(0);
    };

    const checkPredictionResult = async (startPrice: number, endPrice: number,isNewCandle:boolean) => {
      
        
     


        const finalClosePrice = endPrice;
        setClosePrice(finalClosePrice);
        setPriceDirection(finalClosePrice > startPrice ? 'up' : 'down');
        // ลงทะเบียน event listener
     


       
        

        if (currentPrediction && betAmount > 0 && isWaitingResult && accessToken) {
            const isCorrect = 
                (currentPrediction === 'up' && finalClosePrice > startPrice) ||
                (currentPrediction === 'down' && finalClosePrice < startPrice);
            
            const winAmount = isCorrect ? betAmount * 2 : 0;

            console.log('Processing result:', {
                isCorrect,
                winAmount,
                prediction: currentPrediction,
                startPrice,
                endPrice: finalClosePrice
            });

            try {
                await createTransaction(accessToken, {
                    Status: 101,
                    GameProvide: 'options',
                    MemberName: users.username,
                    TransactionAmount: winAmount.toString(),
                    ProductID: 9000,
                    BeforeBalance: balance.toString(),
                    Balance: (balance + winAmount).toString(),
                    AfterBalance: (balance + winAmount).toString()
                });

                // อัพเดทผลลัพธ์
                setBalance(prev => prev + winAmount);
                setLastBetResult(isCorrect ? 'win' : 'lose');
                
                // Reset states
                setBetAmount(0);
                setCurrentPrediction(null);
                setIsWaitingResult(false);
                setSelectedLeverage(1);
                setLeverageAmount(0);
            } catch (error) {
                console.error('Result processing error:', error);
                if (!accessToken) {
                    router.push(`/${lng}/login`);
                }
            }
       }
    };

    const memoizedChart = useMemo(() => (
        <ChartComponent 
            data={initialData}
            onPriceUpdate={setCurrentPrice}
            onOpenPrice={setPredictionStartPrice}
            onCheckPrediction={checkPredictionResult}
            onBettingStateChange={(canBet) => setIsPredictionDisabled(!canBet)}
        />
    ), [initialData]);

    // Add useEffect for countdown and candle management
    // useEffect(() => {
    //     if (countdown <= 45) {
    //         setIsPredictionDisabled(true);
    //     }
        
    //     if (countdown === 60) {
    //         if (currentPrediction && isWaitingResult && betPrice && closePrice) {
    //             console.log('Checking final result:', {
    //                 betPrice,
    //                 closePrice,
    //                 currentPrediction,
    //                 isWaitingResult,
    //                 countdown
    //             });
    //             checkPredictionResult(betPrice, closePrice);
    //         }

    //         // Reset states
    //         setIsNewCandle(true);
    //         setIsPredictionDisabled(false);
    //         setLastPrediction(null);
    //         setIsWaitingResult(false);
    //         setCurrentPrediction(null);
    //         setLastBetResult(null);
    //         setBetPrice(null);
    //         setClosePrice(null);
    //         setPriceDirection(null);
    //     }
    // }, [
    //     countdown,
    //     betPrice,
    //     closePrice,
    //     currentPrediction,
    //     isWaitingResult,
    //     checkPredictionResult,
    //     accessToken,
    //     betAmount,
    //     balance
    // ]);

    return (
        <div className="flex flex-col h-[500px] max-h-screen bg-[#1A1C24] max-w-[1024px] mx-auto w-full">
            
            <div className="h-14 bg-[#12141C] flex items-center justify-between px-4 border-b border-gray-800 w-full">
                <div className="flex items-center space-x-4">
                    {/* Menu and Grid Icons */}
                    <div className="flex space-x-2">
                        <button className="p-2 text-gray-400 hover:text-white">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                            </svg>
                        </button>
                        <button className="p-2 text-gray-400 hover:text-white">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                            </svg>
                        </button>
                    </div>
                    
                    {/* User Info */}
                    <div className="flex items-center space-x-4">
                        <div className="flex flex-col">
                            <span className="text-white font-semibold">
                                {users?.username || 'Loading...'}
                            </span>
                            <span className="text-gray-400 text-sm">
                                Balance: ${balance || '0.00'}
                            </span>
                        </div>
                    </div>
              
                    <div className="flex flex-col">
                            <span className="text-white font-semibold">
                               {"Bet Price"}
                            </span>
                            <span className="text-gray-400 text-sm">
                                {betPrice && `${betPrice.toFixed(2)}`}
                            </span>
                    </div>
                    <div className="flex flex-col">
                            <span className="text-white font-semibold">
                               {"Close Price"}
                            </span>
                            <span className="text-gray-400 text-sm">
                                {isWaitingResult && `${closePrice?.toFixed(2)}`}
                            </span>
                    </div>
              
                </div>
                <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                {[5, 10, 15, 20, 25].map((leverage) => (
                             <button
                             key={leverage}
                             onClick={() => handleLeverageClick(leverage)}
                             disabled={!isBettingPeriod || isProcessingBet || isWaitingResult}
                             className={`px-3 py-1 rounded ${
                                 selectedLeverage === leverage
                                     ? 'bg-blue-500 text-white'
                                     : 'bg-gray-200 text-gray-700'
                             } ${
                                 !isBettingPeriod || isProcessingBet || isWaitingResult 
                                     ? 'opacity-50 cursor-not-allowed' 
                                     : 'hover:bg-blue-400'
                             }`}
                         >
                             {leverage}x
                         </button>
                        ))}
                    </div>
                    <span className="text-white">${leverageAmount.toFixed(2)}</span>
                    <button
                    onClick={handleClearLeverage}
                    disabled={!isBettingPeriod || isProcessingBet || isWaitingResult}
                    className={`px-3 py-1 rounded bg-red-500 text-white
                        ${!isBettingPeriod || isProcessingBet || isWaitingResult 
                            ? 'opacity-50 cursor-not-allowed' 
                            : 'hover:bg-red-400'
                        }`}
                >
                    Clear
                </button>
                </div>
            </div>
            
            <div className="flex flex-1">
                {/* Left Sidebar */}
                <div className="w-14 bg-[#12141C] flex flex-col items-center py-4 space-y-4">
                    <button className="p-2 text-gray-400 hover:text-white">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                        </svg>
                    </button>
                    {/* Add more sidebar icons as needed */}
                </div>
                <div className="flex-1 flex">
                    {/* Chart */}
                    <div className="flex-1 relative">
                    {/* <ChartComponent {...props} data={initialData} countdown={countdown} /> */}
                    {memoizedChart}
                    </div>
                    <div className="w-16 flex flex-col justify-center items-center space-y-2 px-2">
                    <button 
                onClick={() => handlePrediction('up')}
                disabled={!isBettingPeriod || isProcessingBet || isWaitingResult}
                className={`w-full py-3 rounded font-bold text-xs text-white relative
                    ${!isBettingPeriod || isProcessingBet || isWaitingResult
                        ? 'bg-green-500/50 cursor-not-allowed' 
                        : 'bg-green-500 hover:bg-green-600 cursor-pointer'} 
                    transition-colors`}
            >
                {isProcessingBet ? (
                    'Processing...'
                ) : isWaitingResult && currentPrediction === 'up' ? (
                    'Processing...'
                ) : (
                    'UP'
                )}
                {isWaitingResult && currentPrediction === 'up' && (
                    <div className="absolute top-0 left-0 w-full h-full flex items-center justify-center">
                        <div className="animate-pulse text-xs">Waiting Result</div>
                    </div>
                )}
            </button>
            
            <button 
                onClick={() => handlePrediction('down')}
                disabled={!isBettingPeriod || isProcessingBet || isWaitingResult}
                className={`w-full py-3 rounded font-bold text-xs text-white relative
                    ${!isBettingPeriod || isProcessingBet || isWaitingResult
                        ? 'bg-red-500/50 cursor-not-allowed' 
                        : 'bg-red-500 hover:bg-red-600 cursor-pointer'} 
                    transition-colors`}
            >
                {isProcessingBet ? (
                    'Processing...'
                ) : isWaitingResult && currentPrediction === 'down' ? (
                    'Processing...'
                ) : (
                    'DOWN'
                )}
                {isWaitingResult && currentPrediction === 'down' && (
                    <div className="absolute top-0 left-0 w-full h-full flex items-center justify-center">
                        <div className="animate-pulse text-xs">Waiting Result</div>
                    </div>
                )}
            </button>

          
                    </div>
                </div>
            </div>
        </div>
    );
}
