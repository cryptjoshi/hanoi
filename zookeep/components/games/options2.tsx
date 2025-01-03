'use client'
import React, { useEffect, useState, useMemo, useRef, memo } from 'react';
import { GetUserInfo } from '@/actions/index';
import useAuthStore from '@/store/auth';
import { useRouter } from 'next/navigation';
import useBetStore, { BetStore } from '@/store/betStore';
import { createChart, ColorType, CrosshairMode,IChartApi,ISeriesApi } from 'lightweight-charts';

interface ChartProps {
    data: any;
    onPriceUpdate: (price: number) => void;
    onOpenPrice: (price: number) => void;
    onCheckPrediction: (startPrice: number, endPrice: number,isNewCandle:boolean) => void;
    onBettingStateChange: (canBet: boolean) => void;
}

const CountdownDisplay = memo(({ countdown }: { countdown: number }) => {
    const isBettingPeriod = countdown > 45;
    const message = isBettingPeriod ? "เวลาเดิมพัน" : "รอผลเดิมพัน";

    // เพิ่ม log เพื่อตรวจสอบการ render
    //console.log('CountdownDisplay rendered:', countdown);

    return (
        <div className="absolute top-4 left-1/2 transform -translate-x-1/2 z-10 bg-black/50 px-4 sm:px-6 py-1 sm:py-2 rounded-lg flex flex-col items-center">
            <span className="text-white font-bold text-sm mb-1">{message}</span>
            <span className="text-white font-bold text-base sm:text-xl">
                {countdown}s
            </span>
        </div>
    );
});

// ChartComponent definition (you can replace this with your actual ChartComponent implementation)
const ChartComponent: React.FC<ChartProps> = ({ data, onPriceUpdate, onOpenPrice, onCheckPrediction, onBettingStateChange }) => {
    // Simulate chart rendering and price updates
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const chartRef = useRef<IChartApi | null>(null);
    const websocketRef = useRef<WebSocket | null>(null);
    const countdownRef = useRef<number>(60);
    const isMounted = useRef<boolean>(false);
    const [displayCountdown, setDisplayCountdown] = useState<number>(60);
    const candlestickSeriesRef = useRef<ISeriesApi<"Candlestick">>();
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

    const setupWebSocket = () => {
        if (websocketRef.current) {
            websocketRef.current.close();
        }

        const ws = new WebSocket('wss://stream.binance.com:9443/ws/btcusdt@kline_1m');
        let currentCandle: any = null;
       
        ws.onopen = () => {
            console.log('WebSocket Connected');
            // เริ่มต้นโดยการซิงค์เวลากับ server
            syncServerTime();
        };

        // เพิ่มฟังก์ชันซิงค์เวลา
        const syncServerTime = async () => {
            try {
                const response = await fetch('https://api.binance.com/api/v3/time');
                const data = await response.json();
                const serverTime = Math.floor(data.serverTime / 1000);
                const secondsInCurrentMinute = serverTime % 60;
                const remainingSeconds = 60 - secondsInCurrentMinute;
                countdownRef.current = remainingSeconds;
                
                console.log('Time Sync:', {
                    serverTime: new Date(serverTime * 1000).toLocaleTimeString(),
                    remainingSeconds
                });
            } catch (error) {
                console.error('Error syncing time:', error);
            }
        };
        ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            const kline = message.k;

            if (kline) {
                const price = parseFloat(kline.c);
                const klineTimeInSeconds = kline.t / 1000;
                const isNewCandle = !currentCandle || klineTimeInSeconds !== currentCandle.time;
                const currentTime = Math.floor(Date.now() / 1000);
                const candleStartTime = Math.floor(kline.t / 1000);
                const candleEndTime = candleStartTime + 60;
                const remainingTime = candleEndTime - currentTime;
                
                // ปรับปรุงการอัพเดท countdown
                if (Math.abs(remainingTime - countdownRef.current) > 1) {
                    countdownRef.current = Math.max(0, remainingTime);
                    setDisplayCountdown(countdownRef.current);
                }
                
                onPriceUpdate(price);
              // console.log(remainingTime)
                onBettingStateChange(remainingTime > 45);
                
                // Check if we are in the betting period
                const isBettingPeriod = remainingTime > 45; // Adjust this condition as needed
             
                if (isNewCandle) {
                    if (currentCandle) {
                      //  console.log(klineTimeInSeconds,currentCandle.time," is ",klineTimeInSeconds !== currentCandle.time)
                        onCheckPrediction(currentCandle.open, currentCandle.close, isNewCandle);
                    }

                    currentCandle = {
                        time: Math.floor(kline.t / 1000),
                        open: parseFloat(kline.o),
                        high: parseFloat(kline.h),
                        low: parseFloat(kline.l),
                        close: parseFloat(kline.c)
                    };

                    candlestickSeriesRef.current?.update(currentCandle);
                    onOpenPrice(currentCandle.open);
                    // countdownRef.current = 60;

                    // console.log('New Candle:', {
                    //     time: new Date(currentCandle.time * 1000).toLocaleTimeString(),
                    //     remainingTime: remainingTime
                    // });
                } else {
                    currentCandle.high = Math.max(currentCandle.high, parseFloat(kline.h));
                    currentCandle.low = Math.min(currentCandle.low, parseFloat(kline.l));
                    currentCandle.close = price;
                    candlestickSeriesRef.current?.update(currentCandle);
                    countdownRef.current = remainingTime;
                  //  console.log("IsNewCandle:", isNewCandle);
                }
            }
        };


        ws.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        ws.onclose = () => {
            console.log('WebSocket Closed');
            setTimeout(() => {
                if (websocketRef.current === ws) {
                    setupWebSocket();
                }
            }, 3000);
        };

        return ws;
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
                    timeVisible: true,
                    secondsVisible: true,
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

        const ws = setupWebSocket();
        websocketRef.current = ws;
        const interval = setInterval(() => {
            const newPrice = Math.random() * 100; // Simulate price update
            onPriceUpdate(newPrice);
        }, 1000); // Update every second

        return () => clearInterval(interval);
    }
    }, [onPriceUpdate]);

    return (
        <div className=" w-full h-full">
            <div ref={chartContainerRef} className="w-full h-full" 
              style={{ 
                position: 'relative', 
                width: '100%', 
                height: '100%',
                backgroundColor: '#1A1C24',
                zIndex: 0,
            }} 
           
            />
            <CountdownDisplay countdown={displayCountdown} />
            {/* <DataDisplay 
                currentPrice={currentPrice}
                priceChange={priceChange}
                volume={volume}
            /> */}
           
        </div>
    );
};

export function Options({ lng, data }: { lng: string; data: any }) {
    const initialData = data;
    const [users, setUsers] = useState<any>(null);
    const [balance, setBalance] = useState<number>(0);
    const [selectedLeverage, setSelectedLeverage] = useState(1);
    const [currentPrice, setCurrentPrice] = useState<number | 0>(0);
    const [closePrice, setClosePrice] = useState<number | 0>(0);
    const [isProcessingBet, setIsProcessingBet] = useState(false);
 
    const [isWaitingResult, setWaitingResultState] = useState<boolean>(() => {
        const storedValue = localStorage.getItem('waitingResult');
        return storedValue === 'true';
    });

    const { isLoggedIn, accessToken } = useAuthStore();
    const { betPrice, betAmount,betPredict, setBetAmount, setBetPrice, setBetPredict } = useBetStore() as BetStore;
    const router = useRouter();
    const [leverageAmount,setLeverageAmount] = useState(0)
    // Update localStorage when isWaitingResult changes
    useEffect(() => {
        localStorage.setItem('waitingResult', JSON.stringify(isWaitingResult));
    }, [isWaitingResult]);

    const handlePrediction = async (prediction: 'up' | 'down') => {
        if (prediction) {
            const calculatedBetAmount = 1 * selectedLeverage;
            if (calculatedBetAmount <= balance && calculatedBetAmount > 0) {
                try {
                    setIsProcessingBet(true);
                    setWaitingResultState(true); // Set waiting result to true
                    
                    setBetPredict(prediction);
                    setBetAmount(calculatedBetAmount);
                    setBetPrice(currentPrice);
                    setBalance(prev => prev - calculatedBetAmount);

                    // Simulate a delay for processing (replace with actual processing logic)
                    await new Promise(resolve => setTimeout(resolve, 2000)); // 2 seconds delay

                } catch (error) {
                    console.error('Betting error:', error);
                } finally {
                    setIsProcessingBet(false);
                    setWaitingResultState(false); // Set waiting result back to false
                }
            }
        }
    };

    const memoizedChart = useMemo(() => (
        <ChartComponent 
            data={initialData}
            onPriceUpdate={setCurrentPrice}
            onOpenPrice={() => {}}
            onCheckPrediction={() => {}}
            onBettingStateChange={() => {}}
        />
    ), [initialData]);

    const handleLeverageClick = (leverage: number)=> {
        //throw new Error('Function not implemented.');
        setLeverageAmount(leverage)
    }

    const handleClearLeverag = (event:any) => {
        //throw new Error('Function not implemented.');
        setLeverageAmount(0)
    }

      

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
                             disabled={ isProcessingBet || isWaitingResult}
                             className={`px-3 py-1 rounded ${
                                 selectedLeverage === leverage
                                     ? 'bg-blue-500 text-white'
                                     : 'bg-gray-200 text-gray-700'
                             } ${
                                  isProcessingBet || isWaitingResult 
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
                    onClick={handleClearLeverag}
                    disabled={ isProcessingBet || isWaitingResult}
                    className={`px-3 py-1 rounded bg-red-500 text-white
                        ${ isProcessingBet || isWaitingResult 
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
                        disabled={ isProcessingBet || isWaitingResult}
                        className={`w-full py-3 w-10 h-10 rounded font-bold text-xs text-white relative
                        ${isProcessingBet || isWaitingResult
                        ? 'bg-green-500/50 cursor-not-allowed' 
                        : 'bg-green-500 hover:bg-green-600 cursor-pointer'} 
                        transition-colors`}
                    >
                        {isProcessingBet ? (
                            ''
                        ) : isWaitingResult && betPredict === 'up' ? (
                            ''
                        ) : (
                            'UP'
                        )}
                        {isWaitingResult && betPredict === 'up' && (
                            <div className="absolute top-0 left-0 w-full h-full flex items-center justify-center">
                                <div className="animate-pulse text-xs">Waiting Result</div>
                            </div>
                        )}
                    </button>
            
                    <button 
                        onClick={() => handlePrediction('down')}
                        disabled={isProcessingBet || isWaitingResult}
                        className={`w-full py-3 w-10 h-10 rounded font-bold text-xs text-white relative
                            ${isProcessingBet || isWaitingResult
                                ? 'bg-red-500/50 cursor-not-allowed' 
                                : 'bg-red-500 hover:bg-red-600 cursor-pointer'} 
                            transition-colors`}
                    >
                        {isProcessingBet ? (
                            ''
                        ) : isWaitingResult && betPredict === 'down' ? (
                            ''
                        ) : (
                            'DOWN'
                        )}
                        {isWaitingResult && betPredict === 'down' && (
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