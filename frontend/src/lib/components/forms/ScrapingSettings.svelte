<script lang="ts">
  import type { Config } from '../../types/config';

  let { config = $bindable(), initialConfig }: {
    config: Config;
    initialConfig: Config;
  } = $props();

  function isFieldDirty(field: keyof Config): boolean {
    return config[field] !== initialConfig[field];
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

    <div class="pt-4 border-t border-gray-700 space-y-3">
      <div>
        <h4 class="text-sm font-medium text-gray-300 mb-2">Filter Lines</h4>
        <div class="bg-gray-700/30 rounded-md p-3">
          <p class="text-sm text-gray-400 italic">
            {config.filter_lines.length > 0
              ? config.filter_lines.join(', ')
              : 'No filters configured'}
          </p>
          <p class="text-xs text-gray-500 mt-2">Interactive editing coming soon (Phase 4)</p>
        </div>
      </div>

      <div>
        <h4 class="text-sm font-medium text-gray-300 mb-2">Custom Lines</h4>
        <div class="bg-gray-700/30 rounded-md p-3">
          <p class="text-sm text-gray-400 italic">
            {config.add_lines.length > 0
              ? config.add_lines.join(', ')
              : 'No custom lines configured'}
          </p>
          <p class="text-xs text-gray-500 mt-2">Interactive editing coming soon (Phase 5)</p>
        </div>
      </div>
    </div>
  </div>
</div>
