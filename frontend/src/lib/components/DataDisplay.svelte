<script lang="ts">
  import { onMount } from 'svelte';
  import type { ScraperData, ScraperState, ScraperPayload, ScraperErrorPayload, LogEntry, URLStatus } from '$lib/types/scraper';
  import { EventsOn } from '../../../wailsjs/runtime';
  import DataGrid from './display/DataGrid.svelte';
  import ErrorCard from './display/ErrorCard.svelte';
  import StatusIndicator from './display/StatusIndicator.svelte';

  const MAX_INFO_LOGS = 100;
  const MAX_LOG_ENTRIES = 500;

  let {
    onNewLinesAdded,
    displayData = $bindable(),
    rawScrapedData = $bindable([]),
    filterConfig,
    urlStatuses = $bindable({}),
    urlStatusList = $bindable([]),
    currentState = $bindable('stopped'),
    logEntries = $bindable([]),
    lastError = $bindable(null)
  }: {
    onNewLinesAdded?: (indices: number[]) => void;
    displayData?: ScraperData[];
    rawScrapedData?: ScraperData[];
    filterConfig?: any;
    urlStatuses?: Record<string, boolean>;
    urlStatusList?: URLStatus[];
    currentState?: ScraperState;
    logEntries?: LogEntry[];
    lastError?: string | null;
  } = $props();
  let scraperState = $state<ScraperState>('stopped');
  let lastUpdated = $state<Date | null>(null);
  let errorMessage = $state<string | null>(null);

  onMount(async () => {
    if (typeof EventsOn !== 'function') {
      console.warn('Wails EventsOn not available yet');
      return;
    }

    // Listen for scraper data updates (processed + raw)
    EventsOn('polled:data', (payload: ScraperPayload) => {
      try {
        if (Array.isArray(payload.data)) {
          const dataChanged = JSON.stringify(displayData) !== JSON.stringify(payload.data);
          if (dataChanged) {
            displayData = payload.data;
          }

          if (Array.isArray(payload.rawData)) {
            rawScrapedData = payload.rawData;
          }

          lastUpdated = new Date(payload.timestamp);
          errorMessage = null;
        } else {
          throw new Error('Invalid data format from backend');
        }
      } catch (e) {
        console.error('Failed to parse scraper data:', e);
        errorMessage = `Failed to parse data: ${e instanceof Error ? e.message : 'Unknown error'}`;
      }
    });

    // Listen for scraper state changes
    EventsOn('polled:state', (state: string) => {
      if (['idle', 'scraping', 'error', 'stopped'].includes(state)) {
        scraperState = state as ScraperState;
        currentState = scraperState;
      }
    });

    // Listen for per-URL status updates
    EventsOn('polled:url-status', (statuses: URLStatus[]) => {
      urlStatusList = statuses;
      const map: Record<string, boolean> = {};
      for (const s of statuses) {
        map[s.url] = s.hasData;
      }
      urlStatuses = map;
    });

    // Listen for scraper errors
    EventsOn('polled:error', (payload: ScraperErrorPayload) => {
      scraperState = 'error';
      currentState = 'error';
      errorMessage = payload.message;
      lastError = payload.message;
    });

    // Listen for backend log entries
    EventsOn('polled:log', (entry: LogEntry) => {
      let entries = [...logEntries, entry];
      const infoDebug = entries.filter(e => e.level === 'INFO' || e.level === 'DEBUG');
      if (infoDebug.length > MAX_INFO_LOGS) {
        const excess = infoDebug.length - MAX_INFO_LOGS;
        let removed = 0;
        entries = entries.filter(e => {
          if ((e.level === 'INFO' || e.level === 'DEBUG') && removed < excess) {
            removed++;
            return false;
          }
          return true;
        });
      }
      if (entries.length > MAX_LOG_ENTRIES) {
        entries = entries.slice(-MAX_LOG_ENTRIES);
      }
      logEntries = entries;
    });
  });

  // Display data directly — no re-filtering (backend sends processed data)
  const safeDisplayData = $derived(displayData ?? []);
  const hasData = $derived(safeDisplayData.length > 0);
  const isEmpty = $derived(safeDisplayData.length === 0);
  const formattedTime = $derived(
    lastUpdated ? lastUpdated.toLocaleTimeString() : 'Never'
  );

  const hasUrlFailure = $derived(
    Object.keys(urlStatuses).length > 0 &&
    Object.values(urlStatuses).some(v => v === false)
  );

  const effectiveStatus = $derived<'ok' | 'stopped' | 'error'>(
    scraperState === 'stopped' ? 'stopped' :
    scraperState === 'error' ? 'error' :
    hasUrlFailure ? 'error' :
    'ok'
  );

  function handleReset() {
    scraperState = 'idle';
    errorMessage = null;
  }
</script>

<div class="flex-1 min-h-0 flex flex-col overflow-hidden">
  <div class="shrink-0 flex items-center justify-between gap-2 px-3 py-2.5 border-b border-gray-700/50">
    <div class="text-xs text-gray-400">Updated: {formattedTime}</div>
    <StatusIndicator {effectiveStatus} />
  </div>

  <div class="flex-1 overflow-y-auto p-3">
    {#if scraperState === 'error'}
      <ErrorCard message={errorMessage || 'Unknown error occurred'} reset={handleReset} />
    {:else if isEmpty && (scraperState === 'idle' || scraperState === 'scraping') && !lastUpdated}
      <div class="flex flex-col items-center justify-center py-12 gap-3">
        <svg class="w-6 h-6 text-gray-500 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
        <p class="text-sm text-gray-400">Waiting for data...</p>
      </div>
    {:else if isEmpty && scraperState === 'stopped'}
      <p class="text-sm text-gray-400 text-center py-12">Scraper is stopped. Use Preview or Start to see data.</p>
    {:else if isEmpty}
      <p class="text-sm text-gray-400 text-center py-12">No data available</p>
    {:else}
      <DataGrid data={safeDisplayData} />
    {/if}
  </div>
</div>
