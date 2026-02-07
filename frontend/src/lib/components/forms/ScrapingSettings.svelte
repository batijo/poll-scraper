<script lang="ts">
  import type { Config, CustomLine } from '../../types/config';
  import CustomLinesModal from '../CustomLinesModal.svelte';

  let { config = $bindable(), initialConfig }: {
    config: Config;
    initialConfig: Config;
  } = $props();

  let showCustomLinesModal = $state(false);

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }

  function handleCustomLinesConfirm(lines: CustomLine[]) {
    config.add_lines = lines;
  }
</script>

<div class="space-y-6">
  <h3 class="text-lg font-medium text-white">Scraping Strategy</h3>

  <div class="space-y-4">
    <div class="space-y-2">
      <label class="flex items-center gap-2 text-gray-200">
        <input
          type="radio"
          name="scraping-strategy"
          checked={!config.with_eq}
          onchange={() => (config.with_eq = false)}
          class="w-4 h-4"
        />
        <span>Table-based parsing</span>
        {#if isFieldDirty('with_eq')}
          <span class="text-xs text-yellow-400">*</span>
        {/if}
      </label>

      <label class="flex items-center gap-2 text-gray-200">
        <input
          type="radio"
          name="scraping-strategy"
          checked={config.with_eq}
          onchange={() => (config.with_eq = true)}
          class="w-4 h-4"
        />
        <span>Equals-based parsing</span>
      </label>
    </div>

    <div class="pt-4 border-t border-gray-700">
      <h4 class="text-sm font-medium text-gray-300 mb-3">Sum Calculation</h4>

      <div class="space-y-3">
        <label class="flex items-center gap-2 text-gray-200">
          <input
            type="checkbox"
            bind:checked={config.add_sum}
            class="w-4 h-4"
          />
          <span>Add sum calculation</span>
          {#if isFieldDirty('add_sum')}
            <span class="text-xs text-yellow-400">*</span>
          {/if}
        </label>

        {#if config.add_sum}
          <div>
            <label for="sum-symbols" class="block text-sm text-gray-300 mb-1">
              Sum Symbols
              {#if isFieldDirty('sum_symbols')}
                <span class="text-xs text-yellow-400">*</span>
              {/if}
            </label>
            <input
              id="sum-symbols"
              type="text"
              bind:value={config.sum_symbols}
              placeholder="e.g., +, -, *"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
            />
          </div>
        {/if}
      </div>
    </div>

    <div class="pt-4 border-t space-y-3 {isFieldDirty('add_lines') ? 'border-yellow-500/50' : 'border-gray-700'}">
      <div>
        <h4 class="text-sm font-medium text-gray-300 mb-2">
          Custom Lines
          {#if isFieldDirty('add_lines')}
            <span class="text-xs text-yellow-400">* Unsaved changes</span>
          {/if}
        </h4>
        <p class="text-sm text-gray-400 mb-3">
          {config.add_lines.length > 0
            ? `${config.add_lines.length} custom line${config.add_lines.length === 1 ? '' : 's'} configured`
            : 'No custom lines configured'}
        </p>
        <button
          type="button"
          onclick={() => (showCustomLinesModal = true)}
          class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md transition-colors text-sm text-white"
        >
          Manage Custom Lines
        </button>
      </div>
    </div>

    <CustomLinesModal
      bind:showModal={showCustomLinesModal}
      lines={config.add_lines}
      onConfirm={handleCustomLinesConfirm}
    />
  </div>
</div>
