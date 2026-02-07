<script lang="ts">
  import type { Config } from '../../types/config';
  import type { ScraperData } from '../../types/scraper';
  import FilterModal from '../FilterModal.svelte';

  let {
    config = $bindable(),
    initialConfig,
    displayData = []
  }: {
    config: Config;
    initialConfig: Config;
    displayData?: ScraperData[];
  } = $props();

  let showModal = $state(false);

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }

  const totalAvailable = $derived(displayData.length);
  const hasFilterConfig = $derived(config.filter_lines !== undefined && config.filter_lines !== null);
  const visibleCount = $derived(
    !hasFilterConfig ? displayData.length : config.filter_lines.length
  );
  const hiddenCount = $derived(totalAvailable - visibleCount);

  const filterStatus = $derived(() => {
    if (!hasFilterConfig) {
      return `All ${totalAvailable} lines shown (no filters)`;
    }
    if (config.filter_lines.length === 0) {
      return `0 of ${totalAvailable} lines shown (all hidden)`;
    }
    return `${visibleCount} of ${totalAvailable} lines shown`;
  });

  const hiddenLinesList = $derived(() => {
    if (!hasFilterConfig || config.filter_lines.length === 0) return '';
    const allIndices = Array.from({ length: totalAvailable }, (_, i) => i + 1);
    const hidden = allIndices.filter(idx => !config.filter_lines.includes(idx));
    if (hidden.length === 0) return '';
    return hidden.length === 1 ? `Line ${hidden[0]} hidden` : `Lines ${hidden.join(', ')} hidden`;
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
    <p class="text-sm text-gray-400 mb-1">{filterStatus()}</p>
    {#if hiddenLinesList()}
      <p class="text-sm text-gray-400 mb-3">{hiddenLinesList()}</p>
    {:else}
      <div class="mb-3"></div>
    {/if}
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
  availableLines={displayData}
  selectedLines={config.filter_lines}
  onConfirm={handleConfirm}
/>
