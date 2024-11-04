'use client'
import React, { useEffect, useRef, useState,useMemo, memo } from 'react';
import { createChart, ColorType, CrosshairMode } from 'lightweight-charts';
import { GetUserInfo } from '@/actions/index';
import { createTransaction } from '@/actions/index';
import useAuthStore from '@/store/auth';
import { useRouter } from 'next/navigation';
interface ChartComponentProps {
    countdown: number;
    onUpdateCountdown: (newCountdown: number) => void;  // เพิ่ม setter function
    onPriceUpdate?: (price: number) => void;  // Add this line
    data?: any;
    onOpenPrice?:(open:number)=> void;
    onCheckPrediction?: (startPrice: number, endPrice: number) => void;  // Add this line
}
export const ChartComponent = (props: ChartComponentProps) => {
    const chartContainerRef = useRef<HTMLDivElement | null>(null);
    const chartRef = useRef<any>(null);
    const candlestickSeriesRef = useRef<any>(null);
    const isMounted = useRef<boolean>(false);
    
    let currentCandle: any = null;

    const fetchBTCData = async () => {
        const response = await fetch('https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1m&limit=1440');
        const data = await response.json();

        return data.map((item: any) => ({
            open: parseFloat(item[1]),
            high: parseFloat(item[2]),
            low: parseFloat(item[3]),
            close: parseFloat(item[4]),
            time: Math.floor(item[0] / 1000),
        }));
    };
  //  const [predictionStartPrice, setPredictionStartPrice] = useState<number | null>(null);

    const setupWebSocket = () => {
        const socket = new WebSocket('wss://stream.binance.com:9443/ws/btcusdt@trade');

        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            const price = parseFloat(message.p);
            const timestamp = message.T;
            const currentTime = Math.floor(timestamp / 1000);
           // console.log('currentCandlet:',currentCandle)
            if (currentCandle) {
                // currentCandle.high = Math.max(currentCandle.high, price);
                // currentCandle.low = Math.min(currentCandle.low, price);
                // currentCandle.close = price;
                // candlestickSeriesRef.current.update(currentCandle);

                // const currentTime = Math.floor(timestamp / 1000);
                // const cancelEndTime = currentCandle.time + 60;
                // const remainingTime = cancelEndTime - currentTime;

                // if(remainingTime >= 0 && remainingTime <= 60){
                //     props.onUpdateCountdown( remainingTime); // ไม่ให้ค่าต่ำกว่า 1

                // }
                if (currentTime >= currentCandle.time + 60) {
                    // สร้างแท่งเทียนใหม่
                    const newCandle = {
                        time: Math.floor(currentTime / 60) * 60,
                        open: currentCandle.close,
                        high: price,
                        low: price,
                        close: price
                    };
                
                    currentCandle = newCandle;
                    
                    candlestickSeriesRef.current.update(currentCandle);
                    // รีเซ็ต countdown เป็น 60 วินาที
                    
                    props.onOpenPrice?.(newCandle.open)
                    props.onUpdateCountdown(60);
                } else {
                    // อัพเดตแท่งเทียนปัจจุบัน
                    currentCandle.high = Math.max(currentCandle.high, price);
                    currentCandle.low = Math.min(currentCandle.low, price);
                    currentCandle.close = price;
                    candlestickSeriesRef.current.update(currentCandle);
                    //props.onPriceUpdate?.(price);
                    // คำนวณเวลาที่เหลือ
                    const remainingTime = Math.min(
                        60,
                        (currentCandle.time + 60) - currentTime
                    );
                    if (remainingTime >= 0) {
                        props.onUpdateCountdown(remainingTime);
                    }
                  //  console.log("remainingTime:",remainingTime)
                  //  console.log("predictionStartPrice:",predictionStartPrice)

                     if (remainingTime === 0) {
                    //    props.onOpenPrice?.(currentCandle.open)
                       props.onCheckPrediction?.(currentCandle.open, price);
                      //  setPredictionStartPrice(null);
                    }
                }
            }
           // if (currentTime >= currentCandle.time + 60) 
            //    props.onOpenPrice?.(price)
            //    else    
            //    props.onPriceUpdate?.(price);

        };

        return { socket };
    };

    useEffect(() => {
        isMounted.current = true;

        if (chartContainerRef.current) {
            chartRef.current = createChart(chartContainerRef.current, {
                layout: {
                    background: { type: ColorType.Solid, color: '#1A1C24' },
                    textColor: '#ffffff',
                },
                width: chartContainerRef.current.clientWidth,
                height: chartContainerRef.current.clientHeight,
                timeScale: {
                    timeVisible: false,
                    secondsVisible: false,
                    borderColor: '#2a2d3e',
                },
                crossHair: {
                    mode: CrosshairMode.Magnet, // Magnet mode for smoother crosshair
                },
                grid: {
                    vertLines: {
                        visible: true,
                        color: 'rgba(255, 255, 255, 0.1)' // ความโปร่งแสง 10%
                    },
                    horzLines: {
                        visible: true,
                        color: 'rgba(255, 255, 255, 0.1)' // ความโปร่งแสง 10%
                    }
                },
                rightPriceScale: {
                    borderColor: '#2a2d3e',
                },
            });

            candlestickSeriesRef.current = chartRef.current.addCandlestickSeries({
                upColor: '#00C853',       // สีเขียวเข้ม
                downColor: '#D50000',     // สีแดงเข้ม
                borderUpColor: '#00E676', // สีเขียวอ่อน
                borderDownColor: '#FF1744', // สีแดงอ่อน
                wickUpColor: '#69F0AE',   // สีเขียวสำหรับไส้เทียน
                wickDownColor: '#FF5252',  // สีแดงสำหรับไส้เทียนเอาเส้น
            });

            const loadData = async () => {
                const btcData = await fetchBTCData();
                candlestickSeriesRef.current.setData(btcData);
                // ปรับ time scale หลังจากโหลดข้อมูล
                const timeScale = chartRef.current.timeScale();
                            
                // คำนวณเวลาสำหรับ 2 ชั่วโมงย้อนหลัง
                const currentTime = Date.now() / 1000;
                const twoHoursAgo = currentTime - (2 * 60 * 60);
                
                // ตั้งค่าช่วงเวลาที่ต้องการแสดง
               

                // ปรับการแสดงผลให้พอดีกับหน้าจอ
                timeScale.fitContent();
                if (btcData.length > 0) {
                    currentCandle = { ...btcData[btcData.length - 1] };
                    timeScale.setVisibleRange({
                        from: twoHoursAgo,
                        to: currentTime
                    });
                }
            };

            loadData();
            const { socket } = setupWebSocket();

            const intervalId = setInterval(() => {
                if (currentCandle) {
                    const newCandle = {
                        open: currentCandle.close,
                        high: currentCandle.high,
                        low: currentCandle.low,
                        close: currentCandle.close,
                        time: currentCandle.time + 60,
                    };
                    candlestickSeriesRef.current.update(newCandle);
                    currentCandle = { ...newCandle };
                }
            }, 60000);
            const handleResize = () => {
                if (chartContainerRef.current && chartRef.current) {
                    chartRef.current.applyOptions({
                        width: chartContainerRef.current.clientWidth,
                        height: chartContainerRef.current.clientHeight
                    });
                }
            };
        
            window.addEventListener('resize', handleResize);
        
 
            return () => {
                isMounted.current = false;
                socket.close();
                clearInterval(intervalId);
                chartRef.current.remove();
                window.removeEventListener('resize', handleResize);
            };
        }
    }, [props]);
    return (
        <div className="relative w-full h-full"> {/* Removed fixed min-heights */}
        {/* <div 
            className="absolute top-4 left-1/2 transform -translate-x-1/2 z-10 bg-black/50 px-4 sm:px-6 py-1 sm:py-2 rounded-lg"
            style={{ backdropFilter: 'blur(4px)' }}
        >
            <span className="text-white font-bold text-base sm:text-xl">{props.countdown}s</span>
        </div> */}
        <div 
            ref={chartContainerRef} 
            style={{ 
                position: 'relative', 
                width: '100%', 
                height: '100%',
                backgroundColor: '#1A1C24',
                zIndex: 0,
            }} 
            className="w-full h-full" // Changed from flex-1 to w-full h-full
        />
    </div>
    );
    // return <div ref={chartContainerRef} style={{ position: 'relative', width: '100%', height: '300px' ,backgroundColor: '#1A1C24'}} />;
};

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
    
    const isBettingPeriod = countdown > 30;
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

    const handlePrediction = async (prediction: 'up' | 'down') => {
        if (!isPredictionDisabled && !isProcessingBet && !isWaitingResult) {
            const baseBetAmount = 1; // เดิมพันพื้นฐาน 1
            const calculatedBetAmount = baseBetAmount * selectedLeverage; // คูณด้ย leverage
            
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
    const CountdownDisplay = memo(({ countdown }: { countdown: number }) => {
        // คำนวณเวลาแยกเป็น 2 ช่วง
        const displayTime = countdown > 30 
            ? countdown - 30  // 30 วินาทีแรก (30-1)
            : countdown;      // 30 วินาทีหลัง (30-1)
    
        const message = countdown > 30 
            ? "กรุณาลงเดิมพัน"
            : "รอผลเดิมพัน";
    
        return (
            <div 
                className="absolute top-4 left-1/2 transform -translate-x-1/2 z-10 bg-black/50 px-4 sm:px-6 py-1 sm:py-2 rounded-lg flex flex-col items-center"
                style={{ backdropFilter: 'blur(4px)' }}
            >
                <span className="text-white font-bold text-sm mb-1">{message}</span>
                <span className="text-white font-bold text-base sm:text-xl">
                    {displayTime}s
                </span>
            </div>
        );
    });

    const checkPredictionResult = async (currentPrice: number, newPrice: number) => {
        //console.log(currentPrediction, betAmount > 0 ,isWaitingResult , accessToken)
        if (currentPrediction && betAmount > 0 && isWaitingResult && accessToken) {
            const isCorrect = (currentPrice < newPrice && currentPrediction === 'up') ||
                              (currentPrice > newPrice && currentPrediction === 'down');
            
            const winAmount = isCorrect ? betAmount : 0; // บวกเพิ่ม 1 เท่าของยอดเดิมพันถ้าถูก

            try {
                await createTransaction(accessToken, {
                    Status: 101,
                    GameProvide: 'options',
                    MemberName: users.username,
                    TransactionAmount: winAmount.toString(),
                    ProductID: 9000,
                    BeforeBalance: balance.toString(),
                    Balance: (balance + winAmount).toString(),
                    AfterBalance:  (balance + winAmount).toString()
                });

                if (isCorrect) {
                    setBalance(prev => prev + winAmount);
                    setLastBetResult('win');
                } else {
                    setLastBetResult('lose');
                }
                
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
            countdown={countdown}
            onUpdateCountdown={setCountdown}
            onPriceUpdate={setCurrentPrice}
            onCheckPrediction={checkPredictionResult}
            onOpenPrice={setPredictionStartPrice}
        />
    ), [initialData]);

    // Add useEffect for countdown and candle management
    useEffect(() => {
        if (countdown <= 5) {
            setIsPredictionDisabled(true);
        }
        if (countdown === 60) {
            setIsNewCandle(true);
            setIsPredictionDisabled(false);
            setLastPrediction(null);
            setIsWaitingResult(false);
            setCurrentPrediction(null);
            setLastBetResult(null);
        }
        
         if (countdown === 0 && currentPrediction && isWaitingResult) {
            checkPredictionResult(predictionStartPrice,currentPrice);
         }
    }, [countdown]);

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
                    <CountdownDisplay countdown={countdown} />
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
                    'UP...'
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
                    'DOWN...'
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
