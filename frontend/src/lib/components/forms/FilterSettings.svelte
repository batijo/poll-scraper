<script lang="ts">
  import type { Config } from '../../types/config';
  import FilterModal from '../FilterModal.svelte';

  let {
    config = $bindable(),
    initialConfig
  }: {
    config: Config;
    initialConfig: Config;
  } = $props();

  let showModal = $state(false);

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }

  const hiddenLinesCount = $derived(() => {
    return config.filter_lines.length;
  });

  const filterStatus = $derived(() => {
    if (config.filter_lines.length === 0) {
      return 'No filters applied (all lines visible)';
    }
    return `${config.filter_lines.length} line${config.filter_lines.length === 1 ? '' : 's'} hidden`;
  });

  const hiddenLinesList = $derived(() => {
    if (config.filter_lines.length === 0) return '';
    return config.filter_lines.sort((a, b) => a - b).join(', ');
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
      <p class="text-gray-400 mt-1">Hidden lines: {hiddenLinesList()}</p>
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
  availableLines={[]}
  selectedLines={config.filter_lines}
  onConfirm={handleConfirm}
/>
