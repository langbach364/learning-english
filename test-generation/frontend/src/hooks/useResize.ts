import { useState, useEffect, useCallback } from 'react';

interface UseResizeProps {
  minWidth: number;
  minHeight: number;
}

export const useResize = ({ minWidth, minHeight }: UseResizeProps) => {
  const [isResizing, setIsResizing] = useState(false);
  const [resizeDirection, setResizeDirection] = useState<'bottom-right' | 'bottom-left' | null>(null);
  const [size, setSize] = useState({ 
    width: 225,
    height: 400
  });
  const [startPos, setStartPos] = useState({ x: 0, y: 0 });
  const [startSize, setStartSize] = useState({ width: 0, height: 0 });
  
  const handleMouseMove = useCallback((e: MouseEvent) => {
    if (!isResizing) return;
    const deltaX = e.clientX - startPos.x;
    const deltaY = e.clientY - startPos.y;

    if (resizeDirection === 'bottom-right') {
      setSize({
        width: Math.max(minWidth, startSize.width + deltaX),
        height: Math.max(minHeight, startSize.height + deltaY)
      });
    } else if (resizeDirection === 'bottom-left') {
      setSize({
        width: Math.max(minWidth, startSize.width - deltaX),
        height: Math.max(minHeight, startSize.height + deltaY)
      });
    }
  }, [isResizing, startPos, startSize, resizeDirection, minWidth, minHeight]);

  const startResize = useCallback((e: React.MouseEvent, direction: 'bottom-right' | 'bottom-left') => {
    setIsResizing(true);
    setResizeDirection(direction);
    setStartPos({ x: e.clientX, y: e.clientY });
    setStartSize({ width: size.width, height: size.height });
  }, [size]);

  useEffect(() => {
    if (isResizing) {
      window.addEventListener('mousemove', handleMouseMove);
      window.addEventListener('mouseup', () => {
        setIsResizing(false);
        setResizeDirection(null);
      });
    }
    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseup', () => {
        setIsResizing(false);
        setResizeDirection(null);
      });
    };
  }, [isResizing, handleMouseMove]);

  return { isResizing, startResize, size, setSize };
};