import React, { useState, useEffect } from 'react';

interface CountdownProps {
  initialSeconds: number;
  onComplete?: () => void;
}

const Countdown: React.FC<CountdownProps> = ({ initialSeconds, onComplete }) => {
  const [seconds, setSeconds] = useState(initialSeconds);

  useEffect(() => {
    const interval = setInterval(() => {
      setSeconds((prevSeconds) => {
        if (prevSeconds <= 1) {
          if (onComplete) onComplete();
          return initialSeconds; // รีเซ็ตกลับไปที่ 20 วินาที
        }
        return prevSeconds - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [initialSeconds, onComplete]);

  return (
    <div className="text-white text-xl font-bold">
      {seconds}s
    </div>
  );
};

export default Countdown;