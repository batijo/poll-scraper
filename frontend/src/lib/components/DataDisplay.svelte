<script lang="ts">
  import { onMount } from 'svelte';
  import type { ScraperData, ScraperState, ScraperPayload, ScraperErrorPayload } from '$lib/types/scraper';
  import { EventsOn } from '../../../wailsjs/runtime';
  import { GetConfig } from '../../../wailsjs/go/main/App';
  import StatusIndicator from './display/StatusIndicator.svelte';
  import DataGrid from './display/DataGrid.svelte';
  import SkeletonCard from './display/SkeletonCard.svelte';
  import EmptyState from './display/EmptyState.svelte';
  import ErrorCard from './display/ErrorCard.svelte';
  import LastUpdated from './display/LastUpdated.svelte';
  import NewLinesWarning from './NewLinesWarning.svelte';

  let {
    onNewLinesAdded,
    displayData = $bindable<ScraperData[]>([]),
    formState = $bindable<any>({})
  }: {
    onNewLinesAdded?: (indices: number[]) => void;
    displayData?: ScraperData[];
    formState?: any;
  } = $props();
  let scraperState = $state<ScraperState>('idle');
  let lastUpdated = $state<Date | null>(null);
  let errorMessage = $state<string | null>(null);
  let expectedLineCount = $state<number>(0);
  let newLines = $state<ScraperData[]>([]);
  let showNewLinesWarning = $state(false);

  let isFirstUpdate = true;

  onMount(async () => {
    // Ensure Wails is ready before registering listeners
    if (typeof EventsOn !== 'function') {
      console.warn('Wails EventsOn not available yet');
      return;
    }

    // Listen for scraper data updates
    EventsOn('polled:data', (payload: ScraperPayload) => {
      console.log('[DataDisplay] Received polled:data', payload);
      try {
        if (Array.isArray(payload.data)) {
          // On first update, initialize expected count to current data length
          if (isFirstUpdate) {
            expectedLineCount = payload.data.length;
            isFirstUpdate = false;
            console.log('[DataDisplay] First update - setting expectedLineCount to', payload.data.length);
          } else {
            // Check for new lines beyond expected count (only after first update)
            if (payload.data.length > expectedLineCount) {
              const newLineCount = payload.data.length - expectedLineCount;
              newLines = payload.data.slice(expectedLineCount);
              console.log('[DataDisplay] New lines detected:', newLineCount, newLines);
              showNewLinesWarning = true;
            }
          }

          // Update display data
          const dataChanged = JSON.stringify(displayData) !== JSON.stringify(payload.data);
          if (dataChanged) {
            displayData = payload.data;
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
      }
    });

    // Listen for scraper errors
    EventsOn('polled:error', (payload: ScraperErrorPayload) => {
      console.log('[DataDisplay] Received polled:error', payload);
      scraperState = 'error';
      errorMessage = payload.message;
    });
  });

  function filterByIndices(data: ScraperData[], filterLines: number[]): ScraperData[] {
    // Empty or undefined filter_lines means show all
    if (!filterLines || filterLines.length === 0) {
      return data;
    }
    // filter_lines uses 1-based indexing
    return data.filter((_, idx) => filterLines.includes(idx + 1));
  }

  const filteredData = $derived(filterByIndices(displayData, formState?.filter_lines ?? []));
  const isEmpty = $derived(filteredData.length === 0);
  const hasData = $derived(filteredData.length > 0);
  const totalLineCount = $derived(displayData.length);
  const filteredLineCount = $derived(
    formState?.filter_lines?.length === 0 || !formState?.filter_lines
      ? displayData.length
      : formState.filter_lines.length
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
    <div class="text-xs text-gray-400 flex-1">{formattedTime}</div>
    {#if totalLineCount > 0}
      <div class="text-xs text-gray-400">
        {filteredLineCount} of {totalLineCount} lines
      </div>
    {/if}
    <StatusIndicator status={scraperState} {hasData} />
  </div>

  <div class="space-y-2">
    {#if scraperState === 'error'}
      <ErrorCard message={errorMessage || 'Unknown error occurred'} reset={handleReset} />
    {:else if isEmpty && scraperState === 'scraping'}
      <SkeletonCard />
    {:else if isEmpty}
      <EmptyState />
    {/if}

    {#if !isEmpty && scraperState !== 'error'}
      <DataGrid data={filteredData} />
    {/if}
  </div>

  <NewLinesWarning
    bind:showWarning={showNewLinesWarning}
    {newLines}
    existingFilterCount={expectedLineCount}
    onAddToFilter={handleAddNewLines}
    onDismiss={handleDismissWarning}
  />
</div>
