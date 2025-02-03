import React from 'react';
import { motion } from 'framer-motion';

interface CollapsibleCardProps {
  title: string;
  children: React.ReactNode;
  isCollapsed: boolean;
  onCollapse: () => void;
  size: {width: number; height: number};
  x: any;
  y: any;
  isDraggable: boolean;
  onResizeStart: (e: React.MouseEvent, direction: string) => void;
}

const CollapsibleCard: React.FC<CollapsibleCardProps> = ({
  title,
  children,
  isCollapsed,
  onCollapse,
  size,
  x,
  y,
  isDraggable,
  onResizeStart
}) => {
  return (
    <motion.div 
      className="collapsible-content"
      style={{ 
        x, 
        y,
        width: size.width,
        height: size.height,
        fontSize: `${Math.max(10, size.width * 0.02)}px`,
        overflow: isCollapsed ? 'hidden' : 'visible',
      }}
      drag={isDraggable}
      dragMomentum={false}
      dragElastic={0}
    >
      <button 
        className="collapse-button"
        onClick={onCollapse}
      >
        {isCollapsed ? '▼' : '▲'}
      </button>

      <div 
        className="resize-handle resize-handle-bottom-right"
        onMouseDown={(e) => onResizeStart(e, 'bottom-right')}
      />
      <div 
        className="resize-handle resize-handle-bottom-left"
        onMouseDown={(e) => onResizeStart(e, 'bottom-left')}
      />
      
      <div className="card-header">
        <span>{title}</span>
      </div>

      <div className={`content-wrapper ${isCollapsed ? 'collapsed' : ''}`}>
        <div className="scroll-container">
          {children}
        </div>
      </div>
    </motion.div>
  );
};
