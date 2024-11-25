import React, { useState, useEffect } from 'react';

interface PredictionButtonsProps {
  onPrediction: (prediction: 'up' | 'down') => void;
  isNewCandle: boolean;
}

const PredictionButtons: React.FC<PredictionButtonsProps> = ({ onPrediction, isNewCandle }) => {
  const [timeLeft, setTimeLeft] = useState(30);
  const [isDisabled, setIsDisabled] = useState(false);

  useEffect(() => {
    if (isNewCandle) {
      // รีเซ็ตเมื่อเริ่มแท่งเทียนใหม่
      setTimeLeft(30);
      setIsDisabled(false);
    }

    const timer = setInterval(() => {
      setTimeLeft((prevTime) => {
        if (prevTime <= 1) {
          setIsDisabled(true);
          clearInterval(timer);
          return 0;
        }
        return prevTime - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [isNewCandle]);

  return (
    <div className="flex flex-col gap-4 items-center">
      <div className="text-white text-lg">
        {timeLeft > 0 ? `${timeLeft}s to predict` : 'Prediction closed'}
      </div>
      <div className="flex gap-4">
        <button
          onClick={() => onPrediction('up')}
          disabled={isDisabled}
          className={`
            px-6 py-3 rounded-lg font-bold text-white
            ${isDisabled 
              ? 'bg-gray-500 cursor-not-allowed' 
              : 'bg-green-500 hover:bg-green-600 active:bg-green-700'}
          `}
        >
          UP ▲
        </button>
        <button
          onClick={() => onPrediction('down')}
          disabled={isDisabled}
          className={`
            px-6 py-3 rounded-lg font-bold text-white
            ${isDisabled 
              ? 'bg-gray-500 cursor-not-allowed' 
              : 'bg-red-500 hover:bg-red-600 active:bg-red-700'}
          `}
        >
          DOWN ▼
        </button>
      </div>
    </div>
  );
};

export default PredictionButtons;