import { useEffect, useState } from "react";

interface MovingBarProps {
  duration?: number; // seconds from edge to center
}

const MovingBars: React.FC<MovingBarProps> = ({ duration = 2 }) => {
  const [progress, setProgress] = useState(0); // at component mount, progress = 0. can go up to 1

  useEffect(() => {
    /**
     * - `requestAnimationFrame`: schedules a call to callback function
     *  - `step`calculates how much the animation will progress in one frame
     *     - without it, the animation will progress faster for screens with higher refresh rates.
     */
    let startTime: number | null = null;

    function step(timestamp: number) {
      if (!startTime) startTime = timestamp;
      const elapsed = (timestamp - startTime) / 1000; // secs
      const newProgress = Math.min(elapsed / duration, 1);
      setProgress(newProgress);
      if (newProgress < 1) {
        requestAnimationFrame(step);
      }
    }

    requestAnimationFrame(step); // requestAnimationFrame passes a decimal number in ms since the end of the previous frame's rendering. "similar to calling `performance.now()` at the start of the callback function." https://developer.mozilla.org/en-US/docs/Web/API/DedicatedWorkerGlobalScope/requestAnimationFrame
  }, [duration]);

  const leftBarLeftPercent = progress * 50;
  const rightBarLeftPercent = 100 - progress * 50;

  return (
    <div className="relative w-full h-6 m-1 bg-black mx-auto">
      <div
        className="absolute h-full w-2 bg-yellow-400"
        style={{
          left: `${leftBarLeftPercent}%`,
          transform: "translateX(-50%)",
        }}
      />
      <div
        className="absolute h-full w-2 bg-yellow-400"
        style={{
          left: `${rightBarLeftPercent}%`,
          transform: "translateX(-50%)",
        }}
      />
    </div>
  );
};

export default MovingBars;
