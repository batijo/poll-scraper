<script lang="ts">
  import type { Config } from '../../types/config';

  let { config = $bindable(), initialConfig }: { config: Config; initialConfig: Config } = $props();

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }
</script>

<section class="space-y-4">
  <h2 class="text-lg font-semibold mb-4 text-white">General Settings</h2>

  <div>
    <label for="update-interval" class="block text-sm font-medium text-gray-300 mb-1">
      Update Interval (ms)
    </label>
    <input
      id="update-interval"
      type="number"
      min="100"
      step="100"
      bind:value={config.update_interval}
      class={`
        w-full px-3 py-2 rounded text-white
        transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
        ${
          isFieldDirty('update_interval')
            ? 'border-2 border-yellow-500 bg-yellow-900/20'
            : 'border border-gray-600 bg-gray-700'
        }
      `}
    />
    {#if isFieldDirty('update_interval')}
      <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
    {/if}
  </div>

  <div>
    <label class="flex items-center gap-3 cursor-pointer">
      <input
        type="checkbox"
        bind:checked={config.debug}
        class={`
          w-5 h-5 rounded cursor-pointer
          transition-colors
          ${
            isFieldDirty('debug')
              ? 'accent-yellow-500 ring-2 ring-yellow-500'
              : 'accent-blue-500'
          }
        `}
      />
      <span class="text-sm font-medium text-gray-300">Debug Mode</span>
    </label>
    {#if isFieldDirty('debug')}
      <p class="text-yellow-400 text-xs mt-1 ml-8">Unsaved change</p>
    {/if}
  </div>

  <div>
    <label class="flex items-center gap-3 cursor-pointer">
      <input
        type="checkbox"
        bind:checked={config.stop_on_line_count_change}
        class={`
          w-5 h-5 rounded cursor-pointer
          transition-colors
          ${
            isFieldDirty('stop_on_line_count_change')
              ? 'accent-yellow-500 ring-2 ring-yellow-500'
              : 'accent-blue-500'
          }
        `}
      />
      <span class="text-sm font-medium text-gray-300">Stop on URL line count change</span>
    </label>
    <p class="text-xs text-gray-500 mt-1 ml-8">Stops the scraper if any URL returns a different number of lines than the first scrape</p>
    {#if isFieldDirty('stop_on_line_count_change')}
      <p class="text-yellow-400 text-xs mt-1 ml-8">Unsaved change</p>
    {/if}
  </div>
</section>
