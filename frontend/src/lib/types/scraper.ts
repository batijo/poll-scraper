export interface ScraperData {
  name: string;
  value: string;
}

export type ScraperState = 'idle' | 'scraping' | 'error';

export interface ScraperPayload {
  data: ScraperData[];
  timestamp: string;
}

export interface ScraperErrorPayload {
  message: string;
  timestamp: string;
}
