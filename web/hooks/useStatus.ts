'use client';
import { useEffect, useState } from 'react';

export interface StatusData {
  workflows: number;
  status: string;
  timestamp: string;
  environment: string;
}

export function useStatus() {
  const [data, setData] = useState<StatusData | null>(null);

  useEffect(() => {
    async function fetchStatus() {
      try {
        const res = await fetch('/api/status');
        if (res.ok) {
          const json = await res.json();
          setData(json);
        }
      } catch {
        // ignore network errors
      }
    }

    fetchStatus();
    const interval = setInterval(fetchStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  return data;
}
