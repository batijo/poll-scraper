<script lang="ts">
  import type { CustomLine } from '../types/config';

  let dialog: HTMLDialogElement;

  let {
    showModal = $bindable(false),
    lines = [],
    onConfirm
  }: {
    showModal?: boolean;
    lines: CustomLine[];
    onConfirm: (lines: CustomLine[]) => void;
  } = $props();

  let internalLines = $state<CustomLine[]>([]);

  $effect(() => {
    if (showModal && dialog) {
      internalLines = lines.map((l) => ({ ...l }));
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
    const cleaned = internalLines.filter((l) => l.name.trim() || l.value.trim());
    onConfirm(cleaned);
    showModal = false;
  }

  function addLine() {
    internalLines = [...internalLines, { name: '', value: '' }];
  }

  function removeLine(index: number) {
    internalLines = internalLines.filter((_, i) => i !== index);
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
    <h2 class="text-xl font-semibold">Custom Lines</h2>
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
      onclick={addLine}
      class="px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-md transition-colors text-sm"
    >
      Add Line
    </button>
  </div>

  <div class="flex-1 overflow-y-auto p-4 min-h-0">
    {#if internalLines.length === 0}
      <p class="text-gray-400 text-center py-8">No custom lines. Click "Add Line" to create one.</p>
    {:else}
      <div class="space-y-2">
        {#each internalLines as line, index}
          <div class="flex items-center gap-2 p-2 bg-gray-800 rounded">
            <span class="text-xs text-gray-500 w-6 text-right shrink-0">{index + 1}</span>
            <input
              type="text"
              bind:value={line.name}
              placeholder="Name"
              class="flex-1 px-2 py-1.5 bg-gray-700 border border-gray-600 rounded text-sm text-white focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
            <input
              type="text"
              bind:value={line.value}
              placeholder="Value"
              class="flex-1 px-2 py-1.5 bg-gray-700 border border-gray-600 rounded text-sm text-white focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
            <button
              type="button"
              onclick={() => removeLine(index)}
              class="text-gray-400 hover:text-red-400 transition-colors shrink-0"
              aria-label="Remove line"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          </div>
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
