<script lang="ts">
  import type { Config } from '../../types/config';
  import type { ScraperData, ScraperState, URLStatus } from '../../types/scraper';
  import FilterModal from '../FilterModal.svelte';

  let {
    config = $bindable(),
    initialConfig,
    rawScrapedData = $bindable([]),
    urlStatusList = $bindable([]),
    scraperState = 'stopped'
  }: {
    config: Config;
    initialConfig: Config;
    rawScrapedData?: ScraperData[];
    urlStatusList?: URLStatus[];
    scraperState?: ScraperState;
  } = $props();

  let showModal = $state(false);

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }

  const totalAvailable = $derived(rawScrapedData.length);
  const hasFilterConfig = $derived(config.filter_lines !== undefined && config.filter_lines !== null);
  const visibleCount = $derived(
    !hasFilterConfig
      ? totalAvailable
      : config.filter_lines.filter((idx) => idx >= 1 && idx <= totalAvailable).length
  );

  const filterStatus = $derived(() => {
    if (totalAvailable === 0) {
      return 'No data available';
    }
    if (!hasFilterConfig || config.filter_lines.length === 0) {
      return `All ${totalAvailable} lines shown (no filters)`;
    }
    return `${visibleCount} of ${totalAvailable} lines shown`;
  });

  function handleOpenModal() {
    showModal = true;
  }

  function handleConfirm(newSelectedLines: number[]) {
    config.filter_lines = newSelectedLines;
  }
</script>

<div class="pt-4 border-t space-y-3 {isFieldDirty('filter_lines') ? 'border-yellow-500/50' : 'border-gray-700'}">
  <div>
    <h4 class="text-sm font-medium text-gray-300 mb-2">
      Filter Lines
      {#if isFieldDirty('filter_lines')}
        <span class="text-xs text-yellow-400">* Unsaved changes</span>
      {/if}
    </h4>
    <p class="text-sm text-gray-400 mb-3">{filterStatus()}</p>
    <button
      type="button"
      onclick={handleOpenModal}
      class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md transition-colors text-sm text-white"
    >
      Manage Filters
    </button>
  </div>
</div>

<FilterModal
  bind:showModal
  bind:rawScrapedData
  bind:urlStatusList
  selectedLines={config.filter_lines}
  {scraperState}
  onConfirm={handleConfirm}
/>
