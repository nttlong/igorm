import React, { useRef, useEffect, useCallback, useState } from 'react';

interface UnvsVirtualScrollProps {
  onDemand?: (callback:()=>void) => void; // nhận pageIndex
  threshold?: number;
  children: React.ReactNode;
  style?: React.CSSProperties;
  hasMore?: boolean;
  pageIndex?: number;
}

const UnvsVirtualScroll: React.FC<UnvsVirtualScrollProps> = ({
  onDemand,
  threshold = 100,
  children,
  style = {},
  hasMore = true,
  pageIndex = 0,
}) => {
  const [data, setData] = useState([]);
  const containerRef = useRef<HTMLDivElement>(null);
  const scrollBodyRef = useRef<HTMLDivElement>(null);
  const lockRef = useRef(false);
//   const [pageIndex, setPageIndex] = useState(0); // lưu page index nội bộ

  const handleScroll = useCallback(() => {
    // debugger;
    const container = containerRef.current;
    if (!container) return;
    if (!scrollBodyRef.current) return;
    if (!onDemand) return;
    
    // const { scrollTop, scrollHeight, clientHeight } = container;
    const rec=scrollBodyRef.current.getBoundingClientRect()
    const bottomPosition=rec.y+rec.height
    const rec2=container.getBoundingClientRect()
    const visibleBottom=rec2.y+rec2.height
    // const delta =container.getBoundingClientRect().height-scrollTop
    if (bottomPosition-visibleBottom<=threshold){
      container.removeEventListener('scroll', handleScroll);
      onDemand(()=>{
        debugger
        container.addEventListener('scroll', handleScroll);
      });
      
    }
    // console.log(delta)

    // if (delta <= 0 && !lockRef.current && onDemand) {
    //   lockRef.current = true;
      
    //   onDemand();
    // }
  }, [onDemand, pageIndex, threshold]);

  useEffect(() => {
    if (hasMore) {
        lockRef.current = false; // cho phép gọi tiếp
      }
    const container = containerRef.current;
    if (!container) return;

    container.addEventListener('scroll', handleScroll);
    return () => container.removeEventListener('scroll', handleScroll);
  }, [handleScroll]);

  // Reset lock khi người dùng cuộn lên để cho phép gọi lại
  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const checkIfNotAtBottom = () => {
      if (scrollBodyRef.current){
        const bottomPosition=scrollBodyRef.current.offsetTop+B.offsetHeight
        const { scrollTop, scrollHeight, clientHeight } = container;
        const distanceToBottom = scrollHeight - scrollTop - clientHeight;
        if (distanceToBottom > threshold) {
          lockRef.current = false;
        }
      }
    };

    //checkIfNotAtBottom();
  }, [children, threshold]);

  return (
    <div
      ref={containerRef}
      className='debug'
      style={{
        overflowY: 'auto',
        flexDirection: 'column',
        display: 'flex',
        height: '100%',
        ...style,
      }}
    >
      <div ref={scrollBodyRef} className='debug' style={{marginBottom:`${threshold}px`}}>
      {children}
      </div>
    </div>
  );
};

export default UnvsVirtualScroll;
