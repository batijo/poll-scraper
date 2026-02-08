<script lang="ts">
  import type { Config } from '../../types/config';

  let { config, initialConfig }: { config: Config; initialConfig: Config } = $props();

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }
</script>

<section class="space-y-6">
  <h2 class="text-lg font-semibold mb-4 text-white">Output Settings</h2>

  <div class="space-y-4">
    <h3 class="text-md font-medium text-gray-300">CSV Output</h3>

    <div>
      <label class="flex items-center gap-3 cursor-pointer">
        <input
          type="checkbox"
          bind:checked={config.write_to_csv}
          class={`
            w-5 h-5 rounded cursor-pointer
            transition-colors
            ${
              isFieldDirty('write_to_csv')
                ? 'accent-yellow-500 ring-2 ring-yellow-500'
                : 'accent-blue-500'
            }
          `}
        />
        <span class="text-sm font-medium text-gray-300">Write to CSV</span>
      </label>
      {#if isFieldDirty('write_to_csv')}
        <p class="text-yellow-400 text-xs mt-1 ml-8">Unsaved change</p>
      {/if}
    </div>

    <div>
      <label for="csv-path" class="block text-sm font-medium text-gray-300 mb-1">
        CSV Path
      </label>
      <input
        id="csv-path"
        type="text"
        bind:value={config.csv_path}
        disabled={!config.write_to_csv}
        placeholder="e.g., ./output/data.csv"
        class={`
          w-full px-3 py-2 rounded text-white
          transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
          ${
            !config.write_to_csv
              ? 'bg-gray-800 opacity-50 cursor-not-allowed border border-gray-600'
              : isFieldDirty('csv_path')
                ? 'border-2 border-yellow-500 bg-yellow-900/20'
                : 'border border-gray-600 bg-gray-700'
          }
        `}
      />
      {#if config.write_to_csv && isFieldDirty('csv_path')}
        <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
      {/if}
    </div>
  </div>

  <div class="space-y-4">
    <h3 class="text-md font-medium text-gray-300">TXT Output</h3>

    <div>
      <label class="flex items-center gap-3 cursor-pointer">
        <input
          type="checkbox"
          bind:checked={config.write_to_txt}
          class={`
            w-5 h-5 rounded cursor-pointer
            transition-colors
            ${
              isFieldDirty('write_to_txt')
                ? 'accent-yellow-500 ring-2 ring-yellow-500'
                : 'accent-blue-500'
            }
          `}
        />
        <span class="text-sm font-medium text-gray-300">Write to TXT</span>
      </label>
      {#if isFieldDirty('write_to_txt')}
        <p class="text-yellow-400 text-xs mt-1 ml-8">Unsaved change</p>
      {/if}
    </div>

    <div>
      <label for="txt-path" class="block text-sm font-medium text-gray-300 mb-1">
        TXT Path
      </label>
      <input
        id="txt-path"
        type="text"
        bind:value={config.txt_path}
        disabled={!config.write_to_txt}
        placeholder="e.g., ./output/data.txt"
        class={`
          w-full px-3 py-2 rounded text-white
          transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
          ${
            !config.write_to_txt
              ? 'bg-gray-800 opacity-50 cursor-not-allowed border border-gray-600'
              : isFieldDirty('txt_path')
                ? 'border-2 border-yellow-500 bg-yellow-900/20'
                : 'border border-gray-600 bg-gray-700'
          }
        `}
      />
      {#if config.write_to_txt && isFieldDirty('txt_path')}
        <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
      {/if}
    </div>

    <div>
      <label for="dataset-name" class="block text-sm font-medium text-gray-300 mb-1">
        Dataset Name{#if config.write_to_txt}<span class="text-red-400 ml-0.5">*</span>{/if}
      </label>
      <input
        id="dataset-name"
        type="text"
        bind:value={config.dataset_name}
        disabled={!config.write_to_txt}
        placeholder="e.g., Q1 2026 Polls"
        class={`
          w-full px-3 py-2 rounded text-white
          transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
          ${
            !config.write_to_txt
              ? 'bg-gray-800 opacity-50 cursor-not-allowed border border-gray-600'
              : isFieldDirty('dataset_name')
                ? 'border-2 border-yellow-500 bg-yellow-900/20'
                : 'border border-gray-600 bg-gray-700'
          }
        `}
      />
      {#if config.write_to_txt && isFieldDirty('dataset_name')}
        <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
      {/if}
      {#if config.write_to_txt}
        <p class="text-xs text-gray-500 mt-1">Required when TXT output is enabled</p>
      {/if}
    </div>
  </div>
</section>
