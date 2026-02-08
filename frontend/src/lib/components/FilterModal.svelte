<script lang="ts">
  import type { ScraperData, ScraperState, URLStatus } from '../types/scraper';
  import { PreviewScrape } from '../../../wailsjs/go/main/App';

  let dialog: HTMLDialogElement;

  let {
    showModal = $bindable(false),
    rawScrapedData = $bindable([]),
    urlStatusList = $bindable([]),
    selectedLines = [],
    scraperState = 'stopped',
    onConfirm
  }: {
    showModal?: boolean;
    rawScrapedData?: ScraperData[];
    urlStatusList?: URLStatus[];
    selectedLines: number[];
    scraperState?: ScraperState;
    onConfirm: (selectedLines: number[]) => void;
  } = $props();

  let internalSelection = $state<number[]>([]);
  let fetchLoading = $state(false);
  let fetchError = $state<string | null>(null);

  const isStopped = $derived(scraperState === 'stopped');
  const needsFetch = $derived(isStopped && rawScrapedData.length === 0);

  // Compute URL-to-line ranges for grouping display
  const urlLineRanges = $derived(() => {
    if (!urlStatusList || urlStatusList.length === 0) return [];
    const ranges: { url: string; startIndex: number; endIndex: number }[] = [];
    let offset = 0;
    for (const status of urlStatusList) {
      if (status.lineCount > 0) {
        ranges.push({
          url: status.url,
          startIndex: offset,
          endIndex: offset + status.lineCount - 1
        });
      }
      offset += status.lineCount;
    }
    return ranges;
  });

  function getUrlForIndex(index: number): string | null {
    const ranges = urlLineRanges();
    for (const r of ranges) {
      if (index >= r.startIndex && index <= r.endIndex) return r.url;
    }
    return null;
  }

  function isFirstLineOfUrl(index: number): boolean {
    const ranges = urlLineRanges();
    return ranges.some(r => r.startIndex === index);
  }

  function getUrlLabel(index: number): string {
    const ranges = urlLineRanges();
    const range = ranges.find(r => r.startIndex === index);
    if (!range) return '';
    try {
      const u = new URL(range.url);
      return u.hostname + u.pathname;
    } catch {
      return range.url;
    }
  }

  $effect(() => {
    if (showModal && dialog) {
      if (!selectedLines || selectedLines.length === 0) {
        internalSelection = Array.from({ length: rawScrapedData.length }, (_, i) => i + 1);
      } else {
        internalSelection = [...selectedLines];
      }
      fetchError = null;
      dialog.showModal();
    } else if (dialog && !showModal) {
      dialog.close();
    }
  });

  function handleClose() {
    showModal = false;
  }

  function handleCancel() {
    showModal = false;
  }

  function handleConfirm() {
    onConfirm(internalSelection);
    showModal = false;
  }

  function handleCheckAll() {
    internalSelection = Array.from({ length: rawScrapedData.length }, (_, i) => i + 1);
  }

  function handleUncheckAll() {
    internalSelection = [];
  }

  function isLineSelected(index: number): boolean {
    return internalSelection.includes(index + 1);
  }

  function toggleLine(index: number) {
    const oneBasedIndex = index + 1;
    if (internalSelection.includes(oneBasedIndex)) {
      internalSelection = internalSelection.filter((i) => i !== oneBasedIndex);
    } else {
      internalSelection = [...internalSelection, oneBasedIndex];
    }
  }

  async function handleFetchData() {
    fetchLoading = true;
    fetchError = null;
    try {
      const result = await PreviewScrape();
      rawScrapedData = result.rawData;
      urlStatusList = result.statuses;
      internalSelection = Array.from({ length: result.rawData.length }, (_, i) => i + 1);
    } catch (e) {
      fetchError = `Failed to fetch data: ${e}`;
    } finally {
      fetchLoading = false;
    }
  }
</script>

<dialog
  bind:this={dialog}
  onclose={handleClose}
  class="bg-gray-900 text-white rounded-lg shadow-2xl max-w-2xl w-full h-[85vh] border border-gray-700 p-0 backdrop:bg-black/75"
  style="margin: auto;"
>
  <div class="flex flex-col h-full">
  <div class="flex-shrink-0 flex items-center justify-between p-4 border-b border-gray-700">
    <h2 class="text-xl font-semibold">Filter Lines</h2>
    <button
      type="button"
      onclick={handleCancel}
      class="text-gray-400 hover:text-white transition-colors"
      aria-label="Close"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M6 18L18 6M6 6l12 12"
        />
      </svg>
    </button>
  </div>

  {#if !needsFetch}
    <div class="flex-shrink-0 flex gap-2 p-4 border-b border-gray-700">
      <button
        type="button"
        onclick={handleCheckAll}
        class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md transition-colors text-sm"
      >
        Check All
      </button>
      <button
        type="button"
        onclick={handleUncheckAll}
        class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md transition-colors text-sm"
      >
        Uncheck All
      </button>
      {#if isStopped}
        <button
          type="button"
          onclick={handleFetchData}
          disabled={fetchLoading}
          class="px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-md transition-colors text-sm disabled:opacity-50 ml-auto"
        >
          {fetchLoading ? 'Fetching...' : 'Refresh Data'}
        </button>
      {/if}
    </div>
  {/if}

  <div class="flex-1 overflow-y-auto p-4 min-h-0">
    {#if needsFetch}
      <div class="flex flex-col items-center justify-center py-12 gap-4">
        <p class="text-gray-400 text-center">
          Scraper is stopped. Fetch data to configure filters.
        </p>
        <button
          onclick={handleFetchData}
          disabled={fetchLoading}
          class="px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-md transition-colors disabled:opacity-50 flex items-center gap-2"
        >
          {#if fetchLoading}
            <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            Fetching...
          {:else}
            Fetch Data
          {/if}
        </button>
        {#if fetchError}
          <p class="text-red-400 text-sm">{fetchError}</p>
        {/if}
      </div>
    {:else if rawScrapedData.length === 0}
      <p class="text-gray-400 text-center py-8">No lines available.</p>
    {:else}
      <div class="space-y-1">
        {#each rawScrapedData as line, index}
          {#if isFirstLineOfUrl(index)}
            <div class="text-xs text-gray-500 font-medium px-2 pt-3 pb-1 border-t border-gray-700 first:border-t-0 first:pt-0 truncate" title={getUrlLabel(index)}>
              {getUrlLabel(index)}
            </div>
          {/if}
          <label
            class="flex items-center gap-2 p-2 hover:bg-gray-800 rounded transition-colors cursor-pointer"
          >
            <input
              type="checkbox"
              checked={isLineSelected(index)}
              onchange={() => toggleLine(index)}
              class="w-4 h-4 rounded accent-blue-500 cursor-pointer"
            />
            <span class="text-sm text-gray-300">
              [{index + 1}] {line.name}: {line.value}
            </span>
          </label>
        {/each}
      </div>
    {/if}
  </div>

  <div class="flex-shrink-0 flex gap-2 p-4 border-t border-gray-700">
    <button
      type="button"
      onclick={handleCancel}
      class="flex-1 px-4 py-2 bg-gray-600 hover:bg-gray-500 rounded-md transition-colors"
    >
      Cancel
    </button>
    <button
      type="button"
      onclick={handleConfirm}
      disabled={needsFetch}
      class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-md transition-colors disabled:opacity-50"
    >
      Confirm
    </button>
  </div>
  </div>
</dialog>
