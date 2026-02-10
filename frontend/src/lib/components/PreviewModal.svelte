<script lang="ts">
  import type { ScraperData } from '$lib/types/scraper';
  import { PreviewURL } from '../../../wailsjs/go/main/App';
  import DataGrid from './display/DataGrid.svelte';

  let dialog: HTMLDialogElement;
  let previewUrl = $state('');
  let previewData = $state<ScraperData[]>([]);
  let loading = $state(false);
  let error = $state<string | null>(null);

  export function open(url: string) {
    previewUrl = url;
    previewData = [];
    error = null;
    loading = true;
    dialog.showModal();

    PreviewURL(url)
      .then((data) => {
        previewData = data ?? [];
      })
      .catch((e) => {
        error = `Failed to scrape: ${e}`;
      })
      .finally(() => {
        loading = false;
      });
  }

  export function close() {
    dialog.close();
  }

  function handleClose() {
    previewUrl = '';
    previewData = [];
    error = null;
  }
</script>

<dialog
  bind:this={dialog}
  onclose={handleClose}
  class="bg-gray-900 text-white rounded-lg shadow-2xl max-w-2xl w-full h-[85vh] border border-gray-700 p-0 backdrop:bg-black/75"
  style="margin: auto;"
>
  <div class="flex flex-col h-full">
    <div class="shrink-0 flex items-center justify-between p-4 border-b border-gray-700">
      <h2 class="text-xl font-semibold">Preview</h2>
      <button
        type="button"
        onclick={close}
        class="text-gray-400 hover:text-white transition-colors"
        aria-label="Close"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="shrink-0 px-4 py-2 border-b border-gray-700">
      <span class="text-sm font-mono text-gray-400 break-all">{previewUrl}</span>
    </div>

    <div class="flex-1 overflow-y-auto p-4 min-h-0">
      {#if loading}
        <div class="flex flex-col items-center justify-center py-12 gap-3">
          <svg class="w-6 h-6 text-gray-500 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          <p class="text-sm text-gray-400">Scraping...</p>
        </div>
      {:else if error}
        <div class="bg-red-900/30 border border-red-700 rounded-lg p-4">
          <p class="text-red-400 text-sm">{error}</p>
        </div>
      {:else if previewData.length === 0}
        <p class="text-sm text-gray-400 text-center py-12">No data returned</p>
      {:else}
        <DataGrid data={previewData} />
      {/if}
    </div>

    <div class="shrink-0 flex gap-2 p-4 border-t border-gray-700">
      <button
        type="button"
        onclick={close}
        class="flex-1 px-4 py-2 bg-gray-600 hover:bg-gray-500 rounded-md transition-colors"
      >
        Close
      </button>
    </div>
  </div>
</dialog>
