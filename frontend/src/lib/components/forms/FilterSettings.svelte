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

<section
  class={`
    rounded-lg p-4 space-y-3 transition-colors
    ${
      isFieldDirty('filter_lines')
        ? 'border-2 border-yellow-500 bg-yellow-900/20'
        : 'bg-gray-800 border border-gray-700'
    }
  `}
>
  <h3 class="text-lg font-semibold text-white">Filter Lines</h3>

  <div class="text-sm text-gray-300">
    <p class="font-medium">{filterStatus()}</p>
    {#if hiddenLinesList()}
      <p class="text-gray-400 mt-1">{hiddenLinesList()}</p>
    {/if}
  </div>

  {#if isFieldDirty('filter_lines')}
    <p class="text-yellow-400 text-xs">Unsaved changes</p>
  {/if}

  <button
    type="button"
    onclick={handleOpenModal}
    class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-white text-sm font-medium"
  >
    Manage Filters
  </button>
</section>

<FilterModal
  bind:showModal
  availableLines={displayData}
  selectedLines={config.filter_lines}
  onConfirm={handleConfirm}
/>
