<script lang="ts">
  import type { ScraperData, ScraperState, LogEntry, PreviewResult, URLStatus } from '$lib/types/scraper';
  import type { Config } from '$lib/types/config';
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

  function handlePreviewComplete(result: PreviewResult) {
    displayData = result.data;
    rawScrapedData = result.rawData;
    urlStatusList = result.statuses;
    const statusMap: Record<string, boolean> = {};
    for (const s of result.statuses) {
      statusMap[s.url] = s.hasData;
    }
    urlStatuses = statusMap;
  }

  function handleAddNewLines(indices: number[]) {
    settingsPanel?.addToFilterLines?.(indices);
  }
</script>

<div class="flex gap-3 p-3 h-screen min-w-[900px] bg-gray-950 overflow-hidden">
  <div class="flex-1 min-w-[480px] max-w-2xl h-full overflow-hidden">
    <SettingsPanel bind:this={settingsPanel} bind:formState bind:savedConfig bind:displayData bind:rawScrapedData {urlStatuses} bind:urlStatusList {scraperState} {logEntries} {lastError} />
  </div>
  <div class="w-96 shrink-0 h-full flex flex-col gap-3 overflow-hidden">
    <div class="flex-1 overflow-y-auto">
      <DataDisplay bind:displayData bind:rawScrapedData bind:urlStatuses bind:urlStatusList bind:currentState={scraperState} bind:logEntries bind:lastError filterConfig={savedConfig} onNewLinesAdded={handleAddNewLines} />
    </div>
    <ScraperActions {scraperState} onPreviewComplete={handlePreviewComplete} />
  </div>
</div>
