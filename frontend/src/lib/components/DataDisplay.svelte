<script lang="ts">
  import { onMount } from 'svelte';
  import type { ScraperData, ScraperState, ScraperPayload, ScraperErrorPayload, LogEntry } from '$lib/types/scraper';
  import { EventsOn } from '../../../wailsjs/runtime';
  import { GetConfig } from '../../../wailsjs/go/main/App';
  import StatusIndicator from './display/StatusIndicator.svelte';
  import DataGrid from './display/DataGrid.svelte';
  import ErrorCard from './display/ErrorCard.svelte';
  import NewLinesWarning from './NewLinesWarning.svelte';

  const MAX_LOG_ENTRIES = 200;

  let {
    onNewLinesAdded,
    displayData = $bindable(),
    filterConfig,
    urlStatuses = $bindable({}),
    currentState = $bindable('idle'),
    logEntries = $bindable([]),
    lastError = $bindable(null)
  }: {
    onNewLinesAdded?: (indices: number[]) => void;
    displayData?: ScraperData[];
    filterConfig?: any;
    urlStatuses?: Record<string, boolean>;
    currentState?: ScraperState;
    logEntries?: LogEntry[];
    lastError?: string | null;
  } = $props();
  let scraperState = $state<ScraperState>('idle');
  let lastUpdated = $state<Date | null>(null);
  let errorMessage = $state<string | null>(null);
  let expectedLineCount = $state<number>(0);
  let newLines = $state<ScraperData[]>([]);
  let showNewLinesWarning = $state(false);

  onMount(async () => {
    // Ensure Wails is ready before registering listeners
    if (typeof EventsOn !== 'function') {
      console.warn('Wails EventsOn not available yet');
      return;
    }

    // Set initial expected count when component mounts
    expectedLineCount = 0;

    // Listen for scraper data updates
    EventsOn('polled:data', (payload: ScraperPayload) => {
      console.log('[DataDisplay] Received polled:data', payload);
      try {
        if (Array.isArray(payload.data)) {
          // Update display data
          const dataChanged = JSON.stringify(displayData) !== JSON.stringify(payload.data);
          if (dataChanged) {
            displayData = payload.data;
          }

          // Initialize expected count on first data load
          if (expectedLineCount === 0 && payload.data.length > 0) {
            expectedLineCount = payload.data.length;
            console.log('[DataDisplay] Initialized expectedLineCount to', payload.data.length);
          } else if (expectedLineCount > 0 && payload.data.length > expectedLineCount) {
            // Detect new lines only if we've previously set an expected count
            const newLineCount = payload.data.length - expectedLineCount;
            if (newLineCount > 0) {
              newLines = payload.data.slice(expectedLineCount);
              console.log('[DataDisplay] New lines detected:', newLineCount);
              showNewLinesWarning = true;
            }
          }

          lastUpdated = new Date(payload.timestamp);
          errorMessage = null;
        } else {
          throw new Error('Invalid data format from backend');
        }
      } catch (e) {
        scraperState = 'error';
        errorMessage = `Failed to parse data: ${e instanceof Error ? e.message : 'Unknown error'}`;
      }
    });

    // Listen for scraper state changes
    EventsOn('polled:state', (state: string) => {
      console.log('[DataDisplay] Received polled:state', state);
      if (['idle', 'scraping', 'error'].includes(state)) {
        scraperState = state as ScraperState;
        currentState = scraperState;
      }
    });

    // Listen for per-URL status updates
    EventsOn('polled:url-status', (statuses: Array<{ url: string; hasData: boolean }>) => {
      const map: Record<string, boolean> = {};
      for (const s of statuses) {
        map[s.url] = s.hasData;
      }
      urlStatuses = map;
    });

    // Listen for scraper errors
    EventsOn('polled:error', (payload: ScraperErrorPayload) => {
      console.log('[DataDisplay] Received polled:error', payload);
      scraperState = 'error';
      currentState = 'error';
      errorMessage = payload.message;
      lastError = payload.message;
    });

    // Listen for backend log entries
    EventsOn('polled:log', (entry: LogEntry) => {
      logEntries = [...logEntries.slice(-(MAX_LOG_ENTRIES - 1)), entry];
    });
  });

  function filterByIndices(data: ScraperData[] | undefined, filterLines: number[] | undefined): ScraperData[] {
    // Handle undefined or empty data
    if (!data || !Array.isArray(data)) {
      return [];
    }
    // undefined/null filter_lines means NO filtering (show all)
    // Empty array means user unchecked all (show nothing)
    if (filterLines === undefined || filterLines === null) {
      return data;
    }
    if (filterLines.length === 0) {
      return []; // Nothing selected = hide all
    }
    // filter_lines uses 1-based indexing - show only selected lines
    return data.filter((_, idx) => filterLines.includes(idx + 1));
  }

  const safeDisplayData = $derived(displayData ?? []);
  const filteredData = $derived(filterByIndices(safeDisplayData, filterConfig?.filter_lines ?? []));
  const hasRawData = $derived(safeDisplayData.length > 0);
  const allFilteredOut = $derived(hasRawData && filteredData.length === 0);
  const isEmpty = $derived(filteredData.length === 0 && !hasRawData);
  const totalLineCount = $derived(safeDisplayData.length);
  const filteredLineCount = $derived(
    !filterConfig?.filter_lines || filterConfig.filter_lines === undefined || filterConfig.filter_lines === null
      ? safeDisplayData.length  // No filter config = show all
      : filterConfig.filter_lines.length  // Show count of selected lines
  );
  const formattedTime = $derived(
    lastUpdated ? lastUpdated.toLocaleTimeString() : 'Never'
  );

  function handleReset() {
    scraperState = 'idle';
    errorMessage = null;
  }

  function handleAddNewLines() {
    // Calculate 1-based indices for new lines
    const newIndices = newLines.map((_, idx) => expectedLineCount + idx + 1);
    if (onNewLinesAdded) {
      onNewLinesAdded(newIndices);
    }
    // Update expected count to include new lines
    expectedLineCount = displayData.length;
    showNewLinesWarning = false;
  }

  function handleDismissWarning() {
    showNewLinesWarning = false;
  }
</script>

<div class="bg-gray-900 rounded-lg p-4 space-y-3 w-full max-w-full overflow-hidden">
  <div class="flex items-center justify-between gap-2">
    <div class="text-xs text-gray-400 flex-1">Updated: {formattedTime}</div>
    {#if totalLineCount > 0}
      <div class="text-xs text-gray-400">
        {filteredLineCount} of {totalLineCount} lines
      </div>
    {/if}
    <StatusIndicator status={scraperState} hasData={hasRawData} />
  </div>

  <div class="space-y-2 min-h-[120px]">
    {#if scraperState === 'error'}
      <ErrorCard message={errorMessage || 'Unknown error occurred'} reset={handleReset} />
    {:else if allFilteredOut}
      <p class="text-sm text-gray-400 text-center py-12">All lines are hidden by filters</p>
    {:else if isEmpty && scraperState === 'idle' && !lastUpdated}
      <div class="flex flex-col items-center justify-center py-12 gap-3">
        <svg class="w-6 h-6 text-gray-500 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
        <p class="text-sm text-gray-400">Waiting for data...</p>
      </div>
    {:else if isEmpty}
      <p class="text-sm text-gray-400 text-center py-12">No data available</p>
    {:else}
      <DataGrid data={filteredData} />
    {/if}
  </div>

  <!-- Warning modal temporarily disabled due to startup issues - will fix in next iteration -->
  <!-- <NewLinesWarning
    bind:showWarning={showNewLinesWarning}
    {newLines}
    existingFilterCount={expectedLineCount}
    onAddToFilter={handleAddNewLines}
    onDismiss={handleDismissWarning}
  /> -->
</div>
