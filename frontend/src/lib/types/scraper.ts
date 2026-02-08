export interface ScraperData {
  name: string;
  value: string;
}

export type ScraperState = 'idle' | 'scraping' | 'error' | 'stopped';

export interface ScraperPayload {
  data: ScraperData[];
  rawData: ScraperData[];
  timestamp: string;
}

export interface ScraperErrorPayload {
  message: string;
  timestamp: string;
}

export interface LogEntry {
  level: string;
  message: string;
  time: string;
}

export interface PreviewResult {
  rawData: ScraperData[];
  data: ScraperData[];
  statuses: { url: string; hasData: boolean }[];
}
