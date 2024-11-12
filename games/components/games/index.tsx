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
                    background: { type: ColorType.Solid, color: 'white' },
                    textColor: 'black',
                },
                width: chartContainerRef.current.clientWidth,
                height: 300,
                timeScale: {
                    timeVisible: true,
                    secondsVisible: false,
                    rightOffset: 2,
                    shiftVisibleRangeOnNewBar: true,
                },
                crossHair: {
                    mode: CrosshairMode.Magnet, // Magnet mode for smoother crosshair
                },
                grid: {
                    vertLines: { color: 'rgba(197, 203, 206, 0.7)' },
                    horzLines: { color: 'rgba(197, 203, 206, 0.7)' },
                },
            });

            candlestickSeriesRef.current = chartRef.current.addCandlestickSeries({
                upColor: '#4FFF00',
                downColor: '#FF4976',
                borderUpColor: '#4FFF00',
                borderDownColor: '#FF4976',
                wickUpColor: '#4FFF00',
                wickDownColor: '#FF4976',
            });

            const loadData = async () => {
                const btcData = await fetchBTCData();
                candlestickSeriesRef.current.setData(btcData);

                if (btcData.length > 0) {
                    currentCandle = { ...btcData[btcData.length - 1] };
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

    return <div ref={chartContainerRef} style={{ position: 'relative', width: '100%', height: '300px' }} />;
};

export function Game(props: any) {
    const initialData: any = [];
    return <ChartComponent {...props} data={initialData}></ChartComponent>;
}
