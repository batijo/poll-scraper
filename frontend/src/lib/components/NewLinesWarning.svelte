<script lang="ts">
  import type { ScraperData } from '../types/scraper';

  let dialog: HTMLDialogElement;

  let {
    showWarning = $bindable(false),
    newLines = [],
    existingFilterCount = 0,
    onAddToFilter,
    onDismiss
  }: {
    showWarning?: boolean;
    newLines: ScraperData[];
    existingFilterCount: number;
    onAddToFilter: () => void;
    onDismiss: () => void;
  } = $props();

  $effect(() => {
    // Only show modal if there are actually new lines
    const hasNewLines = newLines.length > 0;

    if (showWarning && hasNewLines && dialog) {
      dialog.showModal();
    } else if (dialog) {
      try {
        if (dialog.open) {
          dialog.close();
        }
      } catch (e) {
        // Dialog might already be closed, ignore error
      }
    }
  });

  function handleDismiss() {
    showWarning = false;
    onDismiss?.();
  }

  function handleShowNewLines() {
    onAddToFilter?.();
    showWarning = false;
  }
</script>

<dialog
  bind:this={dialog}
  onclick={(e) => {
    // Close modal when clicking on backdrop (outside the dialog element)
    if (e.target === dialog) {
      handleDismiss();
    }
  }}
  class="bg-gray-900 text-white rounded-lg shadow-2xl max-w-2xl w-full max-h-[80vh] border-2 border-yellow-600 p-0 backdrop:bg-black/75"
>
  <div class="flex flex-col h-full">
  <div class="shrink-0 flex items-center justify-between p-4 border-b border-yellow-600">
    <div class="flex items-center gap-2">
      <svg class="w-6 h-6 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
        />
      </svg>
      <h2 class="text-xl font-semibold">New Lines Detected</h2>
    </div>
    <button
      type="button"
      onclick={handleDismiss}
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

  <div class="flex-1 overflow-y-auto p-4 min-h-0">
    <p class="text-gray-300 mb-4">
      The scraper found <strong class="text-yellow-400">{newLines.length}</strong> new line(s) that
      weren't in your original configuration. These lines are currently hidden by default.
    </p>

    <div class="space-y-1 max-h-64 overflow-y-auto">
      {#each newLines as line, index}
        <div class="flex items-start gap-2 p-2 bg-yellow-900/20 rounded border border-yellow-800">
          <span class="text-xs text-yellow-400 font-semibold shrink-0">[New]</span>
          <span class="text-sm text-gray-300">
            {line.name}: {line.value}
          </span>
        </div>
      {/each}
    </div>
  </div>

  <div class="shrink-0 flex gap-2 p-4 border-t border-yellow-600">
    <button
      type="button"
      onclick={handleDismiss}
      class="flex-1 px-4 py-2 bg-gray-600 hover:bg-gray-500 rounded-md transition-colors"
    >
      Keep Hidden
    </button>
    <button
      type="button"
      onclick={handleShowNewLines}
      class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-md transition-colors"
    >
      Show New Lines
    </button>
  </div>
  </div>
</dialog>
