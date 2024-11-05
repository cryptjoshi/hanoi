import React, { useState, useRef, useCallback, useMemo, useEffect, memo } from 'react';

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
const DataDisplay = memo(({ 
    currentPrice, 
    priceChange, 
    volume 
}: { 
    currentPrice: number;
    priceChange: number;
    volume: number;
}) => {
    const priceChangeColor = priceChange >= 0 ? 'text-green-500' : 'text-red-500';
    const formattedVolume = new Intl.NumberFormat('en-US', {
        maximumFractionDigits: 2
    }).format(volume);

    return (
        <div className="absolute top-4 right-4 bg-black/50 p-4 rounded-lg" style={{ backdropFilter: 'blur(4px)' }}>
            <div className="grid grid-cols-1 gap-2">
                <div>
                    <span className="text-gray-400 text-sm">Price:</span>
                    <span className="text-white ml-2">${currentPrice.toFixed(2)}</span>
                </div>
                <div>
                    <span className="text-gray-400 text-sm">Change:</span>
                    <span className={`ml-2 ${priceChangeColor}`}>
                        {priceChange >= 0 ? '+' : ''}{priceChange.toFixed(2)}%
                    </span>
                </div>
                <div>
                    <span className="text-gray-400 text-sm">Volume:</span>
                    <span className="text-white ml-2">{formattedVolume}</span>
                </div>
            </div>
        </div>
    );
});

const ChartComponent: React.FC<ChartProps> = (props) => {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const chartRef = useRef<IChartApi | null>(null);
    const candlestickSeriesRef = useRef<ISeriesApi<"Candlestick">>();
    const volumeSeriesRef = useRef<ISeriesApi<"Histogram">>();
    const websocketRef = useRef<WebSocket | null>(null);
    const currentCandleRef = useRef<any>(null);
    
    // ใช้ useRef สำหรับค่าที่ต้องการอัพเดทบ่อยๆ
    const countdownRef = useRef<number>(60);
    const [displayCountdown, setDisplayCountdown] = useState<number>(60);
    const [currentPrice, setCurrentPrice] = useState<number>(0);
    const [priceChange, setPriceChange] = useState<number>(0);
    const [volume, setVolume] = useState<number>(0);

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
    useEffect(() => {
        let lastTick = Date.now();
        const intervalId = setInterval(() => {
            const now = Date.now();
            const delta = Math.floor((now - lastTick) / 1000);
            if (delta >= 1) {
                countdownRef.current = Math.max(0, countdownRef.current - delta);
                setDisplayCountdown(countdownRef.current);
                lastTick = now;
            }
        }, 100); // ปรับให้เช็คถี่ขึ้นเป็นทุก 100ms

        return () => clearInterval(intervalId);
    }, []);

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

                props.onPriceUpdate(price);
                props.onBettingStateChange(remainingTime > 45);

                // Check if we are in the betting period
                const isBettingPeriod = remainingTime > 45; // Adjust this condition as needed
             
                if (isNewCandle) {
                    if (currentCandle) {
                      //  console.log(klineTimeInSeconds,currentCandle.time," is ",klineTimeInSeconds !== currentCandle.time)
                        props.onCheckPrediction(currentCandle.open, currentCandle.close, isNewCandle);
                    }

                    currentCandle = {
                        time: Math.floor(kline.t / 1000),
                        open: parseFloat(kline.o),
                        high: parseFloat(kline.h),
                        low: parseFloat(kline.l),
                        close: parseFloat(kline.c)
                    };

                    candlestickSeriesRef.current?.update(currentCandle);
                    props.onOpenPrice(currentCandle.open);
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

    // เริ่ม WebSocket เมื่อ component mount
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
            if (websocketRef.current) {
                websocketRef.current.close();
                websocketRef.current = null;
                isMounted.current = false;
                chartRef.current.remove();
                window.removeEventListener('resize', handleResize);
            }
        };
    }
    }, [props]);

    
    // แสดงผลการนับถอยหลัง
    

    return (
        <div className="relative w-full h-full">
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

export default ChartComponent; 