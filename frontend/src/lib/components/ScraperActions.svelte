<script lang="ts">
  import { StartScraper, StopScraper, PreviewScrape } from '../../../wailsjs/go/main/App';
  import type { ScraperState, PreviewResult } from '../types/scraper';
  import ConfirmDialog from './ConfirmDialog.svelte';

  let {
    scraperState,
    onPreviewComplete
  }: {
    scraperState: ScraperState;
    onPreviewComplete: (result: PreviewResult) => void;
  } = $props();

  let loading = $state(false);
  let previewLoading = $state(false);
  let error = $state<string | null>(null);
  let confirmDialog: ConfirmDialog;

  const isRunning = $derived(scraperState === 'scraping' || scraperState === 'idle');

  async function handleStart() {
    loading = true;
    error = null;
    try {
      await StartScraper();
    } catch (e) {
      error = `Failed to start: ${e}`;
    } finally {
      loading = false;
    }
  }

  function handleStopRequest() {
    confirmDialog.open();
  }

  async function handleStopConfirm() {
    loading = true;
    error = null;
    try {
      await StopScraper();
    } catch (e) {
      error = `Failed to stop: ${e}`;
    } finally {
      loading = false;
    }
  }

  async function handlePreview() {
    previewLoading = true;
    error = null;
    try {
      const result = await PreviewScrape();
      onPreviewComplete(result);
    } catch (e) {
      error = `Preview failed: ${e}`;
    } finally {
      previewLoading = false;
    }
  }
</script>

<div class="bg-gray-900 rounded-lg p-3 space-y-2 flex-shrink-0">
  {#if error}
    <div class="bg-red-900/30 border border-red-700 rounded p-2 text-red-200 text-xs">
      {error}
    </div>
  {/if}

  <div class="flex gap-2 justify-end">
    <button
      onclick={handlePreview}
      disabled={previewLoading || loading}
      class="px-3 py-1 text-sm rounded bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
    >
      {#if previewLoading}
        <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        Previewing...
      {:else}
        Preview
      {/if}
    </button>

    {#if isRunning}
      <button
        onclick={handleStopRequest}
        disabled={loading}
        class="px-3 py-1 text-sm rounded bg-red-600 hover:bg-red-500 text-white disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
      >
        {#if loading}
          <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          Stopping...
        {:else}
          Stop
        {/if}
      </button>
    {:else}
      <button
        onclick={handleStart}
        disabled={loading}
        class="px-3 py-1 text-sm rounded bg-green-600 hover:bg-green-500 text-white disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
      >
        {#if loading}
          <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          Starting...
        {:else}
          Start
        {/if}
      </button>
    {/if}
  </div>
</div>

<ConfirmDialog
  bind:this={confirmDialog}
  title="Stop Scraper"
  message="Are you sure you want to stop the scraper?"
  confirmText="Stop"
  confirmClass="bg-red-600 hover:bg-red-500"
  onConfirm={handleStopConfirm}
/>
