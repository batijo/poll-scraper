<script lang="ts">
  import type { ScraperData } from '$lib/types/scraper';
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';

  let settingsPanel: any;
  let displayData = $state<ScraperData[]>([]);
  let formState = $state<any>({});

  function handleAddNewLines(indices: number[]) {
    settingsPanel?.addToFilterLines?.(indices);
  }
</script>

<div class="flex gap-3 p-3 h-screen min-w-[900px] bg-gray-950 overflow-hidden">
  <div class="flex-1 min-w-[480px] max-w-2xl h-full overflow-hidden">
    <SettingsPanel bind:this={settingsPanel} bind:formState displayData={displayData} />
  </div>
  <div class="w-96 flex-shrink-0 h-full overflow-y-auto">
    <DataDisplay bind:displayData bind:formState onNewLinesAdded={handleAddNewLines} />
  </div>
</div>
