'use client';
import styles from './Panel.module.css';
import { ReactNode } from 'react';

interface PanelProps {
  title: string;
  children: ReactNode;
}

export default function Panel({ title, children }: PanelProps) {
  return (
    <div className={styles.panel}>
      <div className={styles.title}>██ {title}</div>
      {children}
    </div>
  );
}
