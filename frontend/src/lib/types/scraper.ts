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

export interface URLStatus {
  url: string;
  hasData: boolean;
  lineCount: number;
  error: boolean;
}

export interface PreviewResult {
  rawData: ScraperData[];
  data: ScraperData[];
  statuses: URLStatus[];
}
