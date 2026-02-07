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
    onNewLinesAdded
  }: {
    onNewLinesAdded?: (indices: number[]) => void;
  } = $props();

  let displayData = $state<ScraperData[]>([]);
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

    // Load initial expected line count from config
    try {
      const config = await GetConfig();
      // If filter_lines is configured, expect only those lines
      // If empty or not configured, expect current displayData length
      expectedLineCount = config.filter_lines.length > 0 ? config.filter_lines.length : displayData.length;
    } catch (e) {
      console.error('Failed to load initial line count', e);
    }

    // Listen for scraper data updates
    EventsOn('polled:data', (payload: ScraperPayload) => {
      console.log('[DataDisplay] Received polled:data', payload);
      try {
        if (Array.isArray(payload.data)) {
          // Only update if data actually changed to prevent unnecessary re-renders
          const dataChanged = JSON.stringify(displayData) !== JSON.stringify(payload.data);
          if (dataChanged) {
            displayData = payload.data;

            // Check for new lines beyond expected count
            if (payload.data.length > expectedLineCount && expectedLineCount > 0) {
              newLines = payload.data.slice(expectedLineCount);
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
      }
    });

    // Listen for scraper errors
    EventsOn('polled:error', (payload: ScraperErrorPayload) => {
      console.log('[DataDisplay] Received polled:error', payload);
      scraperState = 'error';
      errorMessage = payload.message;
    });
  });

  const isEmpty = $derived(displayData.length === 0);
  const hasData = $derived(displayData.length > 0);
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
      <DataGrid data={displayData} />
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
