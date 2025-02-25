import { useEffect, useRef, type RefObject } from 'react';

export function useScrollToBottom<T extends HTMLElement>(): [
  RefObject<T>,
  RefObject<T>,
] {
  const containerRef = useRef<T>(null);
  const endRef = useRef<T>(null);

  useEffect(() => {
    const container = containerRef.current;
    const end = endRef.current;

    if (container && end) {
      // Scroll after short delay to ensure content is rendered before we scroll
      const initialScroll = () => {
        setTimeout(() => {
          end.scrollIntoView({ behavior: 'instant', block: 'end' });
        }, 50);
      };

      const observer = new MutationObserver(() => {
        end.scrollIntoView({ behavior: 'instant', block: 'end' });
      });

      observer.observe(container, {
        childList: true,
        attributes: true,
        characterData: true,
      });

      initialScroll();

      return () => {
        observer.disconnect();
      };
    }
  }, []);

  return [containerRef, endRef];
}
