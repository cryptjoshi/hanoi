'use client'
import React, { useEffect, useRef } from 'react';
import { createChart, ColorType, CrosshairMode } from 'lightweight-charts';

export const ChartComponent = (props: any) => {
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

    const setupWebSocket = () => {
        const socket = new WebSocket('wss://stream.binance.com:9443/ws/btcusdt@trade');

        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            const price = parseFloat(message.p);

            if (currentCandle) {
                currentCandle.high = Math.max(currentCandle.high, price);
                currentCandle.low = Math.min(currentCandle.low, price);
                currentCandle.close = price;
                candlestickSeriesRef.current.update(currentCandle);
            }
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

            return () => {
                isMounted.current = false;
                socket.close();
                clearInterval(intervalId);
                chartRef.current.remove();
            };
        }
    }, [props]);

    return <div ref={chartContainerRef} style={{ position: 'relative', width: '100%', height: '300px' ,backgroundColor: '#1A1C24'}} />;
};

export function Options(props: any) {
    const initialData: any = [];
    return (
        <div className="flex flex-col h-[500px] max-h-screen bg-[#1A1C24]">
            
            <div className="h-14 bg-[#12141C] flex items-center justify-between px-4 border-b border-gray-800">
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
                    
                    {/* Currency Pair */}
                    <div className="flex items-center space-x-2">
                        <span className="text-white font-semibold">BTC/USDT</span>
                        <span className="text-gray-400">Forex</span>
                    </div>
                </div>
                <div className="flex items-center space-x-4">
                    <span className="text-white">$1.0876</span>
                    <button className="px-4 py-1.5 bg-green-600 text-white text-sm rounded">
                       Balance
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
                    <ChartComponent {...props} data={initialData} />
                    </div>
                    <div className="w-16 flex flex-col justify-center items-center space-y-2 px-2">
                        <button className="w-full py-3 bg-green-500 text-white rounded hover:bg-green-600 transition-colors font-bold text-xs">
                            BUY
                        </button>
                        <button className="w-full py-3 bg-red-500 text-white rounded hover:bg-red-600 transition-colors font-bold text-xs">
                            SELL
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
