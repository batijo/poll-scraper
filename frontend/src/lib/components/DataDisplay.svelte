<script lang="ts">
  import type { ScraperData, ScraperState } from '$lib/types/scraper';
  import StatusIndicator from './display/StatusIndicator.svelte';
  import DataGrid from './display/DataGrid.svelte';
  import SkeletonCard from './display/SkeletonCard.svelte';
  import EmptyState from './display/EmptyState.svelte';
  import ErrorCard from './display/ErrorCard.svelte';
  import LastUpdated from './display/LastUpdated.svelte';

  let displayData = $state<ScraperData[]>([]);
  let scraperState = $state<ScraperState>('idle');
  let lastUpdated = $state<Date | null>(null);
  let errorMessage = $state<string | null>(null);

  const isEmpty = $derived(displayData.length === 0);
  const formattedTime = $derived(
    lastUpdated ? lastUpdated.toLocaleTimeString() : 'Never'
  );

  function handleReset() {
    scraperState = 'idle';
    errorMessage = null;
  }
</script>

<div class="bg-gray-900 rounded-lg p-6 space-y-6">
  <div class="flex items-center justify-between">
    <h2 class="text-lg font-semibold text-white">Data Display</h2>
    <StatusIndicator status={scraperState} />
  </div>

  <div class="space-y-4">
    {#if scraperState === 'scraping'}
      <SkeletonCard />
    {:else if scraperState === 'error'}
      <ErrorCard message={errorMessage || 'Unknown error occurred'} reset={handleReset} />
    {:else if isEmpty}
      <EmptyState />
    {:else}
      <svelte:boundary>
        <DataGrid data={displayData} />
        {#snippet failed(error, reset)}
          <ErrorCard message={error.message} {reset} />
        {/snippet}
      </svelte:boundary>
    {/if}
  </div>

  {#if !isEmpty && scraperState !== 'scraping' && scraperState !== 'error'}
    <div class="pt-2 border-t border-gray-700">
      <LastUpdated time={formattedTime} />
    </div>
  {/if}
</div>
