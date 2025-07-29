'use client';
import './globals.css';
import styles from './page.module.css';
import Panel from '../components/Panel';
import { useStatus } from '../hooks/useStatus';

export default function Page() {
  const status = useStatus();

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.title}>n8n-ops MONITORING DASHBOARD</div>
        <div className={styles.subtitle}>REAL-TIME DATA</div>
      </div>
      <div className={styles.grid}>
        <Panel title="SYSTEM STATUS">
          <div className={styles.ok}>CLI: ONLINE</div>
          <div className={styles.ok}>MOCK API: ACTIVE</div>
        </Panel>
        <Panel title="LIVE METRICS">
          <div className={styles.info}>WORKFLOWS: {status?.workflows ?? '...'}</div>
          <div className={styles.info}>TIMESTAMP: {status?.timestamp}</div>
        </Panel>
        <Panel title="ENVIRONMENT">
          <div className={styles.ok}>{status?.environment ?? '...'}</div>
        </Panel>
        <Panel title="COMMAND HISTORY">
          <div className={styles.logLine}>n8n-ops sync --demo</div>
          <div className={styles.logLine}>n8n-ops monitor --demo</div>
        </Panel>
        <Panel title="LIVE EVENTS">
          <div className={`${styles.logLine} ${styles.error}`}>[14:23:45] MOCK: Workflow 1002 failed</div>
          <div className={`${styles.logLine} ${styles.warn}`}>[14:23:45] MONITOR: Failure count 2</div>
          <div className={`${styles.logLine} ${styles.info}`}>[14:23:44] DAEMON: File watcher detected change</div>
        </Panel>
      </div>
    </div>
  );
}
