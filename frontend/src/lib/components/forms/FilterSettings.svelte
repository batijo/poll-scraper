<script lang="ts">
  import type { Config, SymbolLine } from '../../types/config';
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

  const symbolCount = $derived(config.symbol_lines?.length ?? 0);
  const isDirty = $derived(isFieldDirty('filter_lines') || isFieldDirty('symbol_lines'));

  function handleOpenModal() {
    showModal = true;
  }

  function handleConfirm(newSelectedLines: number[], newSymbolLines: SymbolLine[]) {
    config.filter_lines = newSelectedLines;
    config.symbol_lines = newSymbolLines;
  }
</script>

<div class="pt-4 border-t space-y-3 {isDirty ? 'border-yellow-500/50' : 'border-gray-700'}">
  <div>
    <h4 class="text-sm font-medium text-gray-300 mb-2">
      Filter Lines &amp; Symbols
      {#if isDirty}
        <span class="text-xs text-yellow-400">* Unsaved changes</span>
      {/if}
    </h4>
    <p class="text-sm text-gray-400 mb-3">
      {filterStatus()}{#if symbolCount > 0}, {symbolCount} with symbol{/if}
    </p>
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
  customLines={config.add_lines}
  symbolLines={config.symbol_lines ?? []}
  {scraperState}
  onConfirm={handleConfirm}
/>
