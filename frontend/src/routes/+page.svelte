<script lang="ts">
  import type { ScraperData } from '$lib/types/scraper';
  import type { Config } from '$lib/types/config';
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
  import { createDefaultConfig } from '$lib/types/config';

  let settingsPanel: any;
  let displayData = $state<ScraperData[]>([]);
  let formState = $state<Config>(createDefaultConfig());
  let savedConfig = $state<Config>(createDefaultConfig());

  function handleAddNewLines(indices: number[]) {
    settingsPanel?.addToFilterLines?.(indices);
  }
</script>

<div class="flex gap-3 p-3 h-screen min-w-[900px] bg-gray-950 overflow-hidden">
  <div class="flex-1 min-w-[480px] max-w-2xl h-full overflow-hidden">
    <SettingsPanel bind:this={settingsPanel} bind:formState bind:savedConfig displayData={displayData} />
  </div>
  <div class="w-96 flex-shrink-0 h-full overflow-y-auto">
    <DataDisplay bind:displayData filterConfig={savedConfig} onNewLinesAdded={handleAddNewLines} />
  </div>
</div>
