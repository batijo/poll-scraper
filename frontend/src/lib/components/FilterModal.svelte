<script lang="ts">
  import type { ScraperData } from '../types/scraper';

  let dialog: HTMLDialogElement;

  let {
    showModal = $bindable(false),
    availableLines = [],
    selectedLines = [],
    onConfirm
  }: {
    showModal?: boolean;
    availableLines: ScraperData[];
    selectedLines: number[];
    onConfirm: (selectedLines: number[]) => void;
  } = $props();

  let internalSelection = $state<number[]>([]);

  $effect(() => {
    if (showModal && dialog) {
      // If no lines are selected (empty or undefined), default to ALL lines selected
      if (!selectedLines || selectedLines.length === 0) {
        internalSelection = Array.from({ length: availableLines.length }, (_, i) => i + 1);
      } else {
        internalSelection = [...selectedLines];
      }
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
    internalSelection = Array.from({ length: availableLines.length }, (_, i) => i + 1);
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
  </div>

  <div class="flex-1 overflow-y-auto p-4 min-h-0">
    {#if availableLines.length === 0}
      <p class="text-gray-400 text-center py-8">No lines available. Start scraping to see data.</p>
    {:else}
      <div class="space-y-1">
        {#each availableLines as line, index}
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
      class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-md transition-colors"
    >
      Confirm
    </button>
  </div>
  </div>
</dialog>
