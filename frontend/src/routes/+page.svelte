<script lang="ts">
  import type { ScraperData, ScraperState, LogEntry } from '$lib/types/scraper';
  import type { Config } from '$lib/types/config';
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
  import { createDefaultConfig } from '$lib/types/config';

  let settingsPanel: any;
  let displayData = $state<ScraperData[]>([]);
  let formState = $state<Config>(createDefaultConfig());
  let savedConfig = $state<Config>(createDefaultConfig());
  let urlStatuses = $state<Record<string, boolean>>({});
  let scraperState = $state<ScraperState>('idle');
  let logEntries = $state<LogEntry[]>([]);
  let lastError = $state<string | null>(null);

  function handleAddNewLines(indices: number[]) {
    settingsPanel?.addToFilterLines?.(indices);
  }
</script>

<div class="flex gap-3 p-3 h-screen min-w-[900px] bg-gray-950 overflow-hidden">
  <div class="flex-1 min-w-[480px] max-w-2xl h-full overflow-hidden">
    <SettingsPanel bind:this={settingsPanel} bind:formState bind:savedConfig displayData={displayData} {urlStatuses} {scraperState} {logEntries} {lastError} />
  </div>
  <div class="w-96 flex-shrink-0 h-full overflow-y-auto">
    <DataDisplay bind:displayData bind:urlStatuses bind:currentState={scraperState} bind:logEntries bind:lastError filterConfig={savedConfig} onNewLinesAdded={handleAddNewLines} />
  </div>
</div>
