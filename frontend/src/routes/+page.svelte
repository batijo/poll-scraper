<script lang="ts">
  import { onMount } from 'svelte';
  import type { ScraperData, ScraperState, LogEntry, PreviewResult, URLStatus } from '$lib/types/scraper';
  import type { Config } from '$lib/types/config';
  import { PreviewScrape } from '../../wailsjs/go/main/App';
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
  import ScraperActions from '$lib/components/ScraperActions.svelte';
  import { createDefaultConfig } from '$lib/types/config';

  let settingsPanel: any;
  let displayData = $state<ScraperData[]>([]);
  let rawScrapedData = $state<ScraperData[]>([]);
  let formState = $state<Config>(createDefaultConfig());
  let savedConfig = $state<Config>(createDefaultConfig());
  let urlStatuses = $state<Record<string, boolean>>({});
  let urlStatusList = $state<URLStatus[]>([]);
  let scraperState = $state<ScraperState>('stopped');
  let logEntries = $state<LogEntry[]>([]);
  let lastError = $state<string | null>(null);

  function updateFromResult(result: PreviewResult) {
    displayData = result.data;
    rawScrapedData = result.rawData;
    urlStatusList = result.statuses;
    const map: Record<string, boolean> = {};
    for (const s of result.statuses) {
      map[s.url] = s.hasData;
    }
    urlStatuses = map;
  }

  // Fetch URL statuses on startup
  onMount(async () => {
    try {
      updateFromResult(await PreviewScrape());
    } catch {
      // URLs will show gray dots if initial fetch fails
    }
  });

  function handleAddNewLines(indices: number[]) {
    settingsPanel?.addToFilterLines?.(indices);
  }
</script>

<div class="flex gap-3 p-3 h-screen min-w-[900px] bg-gray-950 overflow-hidden">
  <div class="flex-1 min-w-[480px] max-w-2xl h-full bg-gray-900 rounded-lg overflow-hidden">
    <SettingsPanel bind:this={settingsPanel} bind:formState bind:savedConfig bind:displayData bind:rawScrapedData bind:urlStatusList {scraperState} {logEntries} {lastError} onConfigSaved={async () => { try { updateFromResult(await PreviewScrape()); } catch {} }} />
  </div>
  <div class="w-96 shrink-0 h-full bg-gray-900 rounded-lg overflow-hidden flex flex-col">
    <DataDisplay bind:displayData bind:rawScrapedData bind:urlStatuses bind:urlStatusList bind:currentState={scraperState} bind:logEntries bind:lastError filterConfig={savedConfig} onNewLinesAdded={handleAddNewLines} />
    <ScraperActions {scraperState} onPreviewComplete={updateFromResult} />
  </div>
</div>
