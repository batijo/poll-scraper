<script lang="ts">
  import URLModal from '../URLModal.svelte';
  import ConfirmDialog from '../ConfirmDialog.svelte';

  let { links = $bindable([]) }: { links: string[] } = $props();

  let urlModal: URLModal;
  let confirmDialog: ConfirmDialog;
  let deleteIndex = $state<number | null>(null);

  function handleAddUrl(url: string) {
    links.push(url);
  }

  function handleOpenAddModal() {
    urlModal.open();
  }

  function handleMoveUp(index: number) {
    if (index === 0) return;
    [links[index - 1], links[index]] = [links[index], links[index - 1]];
  }

  function handleMoveDown(index: number) {
    if (index === links.length - 1) return;
    [links[index], links[index + 1]] = [links[index + 1], links[index]];
  }

  function handleOpenDeleteDialog(index: number) {
    deleteIndex = index;
    confirmDialog.open();
  }

  function handleDelete() {
    if (deleteIndex !== null) {
      links.splice(deleteIndex, 1);
      deleteIndex = null;
    }
  }
</script>

<div class="space-y-4">
  <div class="flex items-center gap-3">
    <h3 class="text-lg font-medium text-white">Scraping URLs</h3>
    <span class="px-2 py-0.5 bg-gray-700 text-gray-300 text-sm rounded-full">
      {links.length}
    </span>
  </div>

  {#if links.length === 0}
    <div class="bg-gray-700/30 rounded-lg p-6 text-center">
      <p class="text-gray-400 mb-4">No URLs configured</p>
      <button
        type="button"
        onclick={handleOpenAddModal}
        class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-md transition-colors"
      >
        Add URL
      </button>
    </div>
  {:else}
    <div class="space-y-2">
      {#each links as url, index}
        <div class="bg-gray-700/50 rounded p-3 flex items-center gap-2">
          <span class="flex-1 text-sm font-mono truncate text-gray-200">
            {url}
          </span>

          <div class="flex gap-1">
            <button
              type="button"
              onclick={() => handleMoveUp(index)}
              disabled={index === 0}
              class="px-2 py-1 bg-gray-600 hover:bg-gray-500 disabled:opacity-30 disabled:cursor-not-allowed rounded text-xs transition-colors"
              title="Move up"
            >
              ↑
            </button>
            <button
              type="button"
              onclick={() => handleMoveDown(index)}
              disabled={index === links.length - 1}
              class="px-2 py-1 bg-gray-600 hover:bg-gray-500 disabled:opacity-30 disabled:cursor-not-allowed rounded text-xs transition-colors"
              title="Move down"
            >
              ↓
            </button>
            <button
              type="button"
              onclick={() => handleOpenDeleteDialog(index)}
              class="px-2 py-1 bg-red-600/80 hover:bg-red-600 rounded text-xs transition-colors"
              title="Delete"
            >
              Delete
            </button>
          </div>
        </div>
      {/each}
    </div>

    <button
      type="button"
      onclick={handleOpenAddModal}
      class="w-full px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-md transition-colors"
    >
      Add URL
    </button>
  {/if}
</div>

<URLModal bind:this={urlModal} onAdd={handleAddUrl} />
<ConfirmDialog
  bind:this={confirmDialog}
  title="Delete URL"
  message="Are you sure you want to delete this URL?"
  onConfirm={handleDelete}
/>
