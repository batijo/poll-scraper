<script lang="ts">
  import URLModal from '../URLModal.svelte';
  import ConfirmDialog from '../ConfirmDialog.svelte';
  import PreviewModal from '../PreviewModal.svelte';

  let {
    links = $bindable([]),
    initialLinks = [],
    urlStatuses = {}
  }: {
    links: string[];
    initialLinks: string[];
    urlStatuses?: Record<string, boolean>;
  } = $props();

  let urlModal: URLModal;
  let confirmDialog: ConfirmDialog;
  let previewModal: PreviewModal;
  let deleteIndex = $state<number | null>(null);

  const isDirty = $derived(JSON.stringify(links) !== JSON.stringify(initialLinks));

  function handleAddUrl(url: string) {
    links = [...links, url];
  }

  function handleEditUrl(index: number, url: string) {
    links = links.map((u, i) => i === index ? url : u);
  }

  function handleOpenAddModal() {
    urlModal.open();
  }

  function handleOpenEditModal(index: number) {
    urlModal.open(links[index], index);
  }

  function handleMoveUp(index: number) {
    if (index === 0) return;
    const updated = [...links];
    [updated[index - 1], updated[index]] = [updated[index], updated[index - 1]];
    links = updated;
  }

  function handleMoveDown(index: number) {
    if (index === links.length - 1) return;
    const updated = [...links];
    [updated[index], updated[index + 1]] = [updated[index + 1], updated[index]];
    links = updated;
  }

  function handleOpenDeleteDialog(index: number) {
    deleteIndex = index;
    confirmDialog.open();
  }

  function handleDelete() {
    if (deleteIndex !== null) {
      links = links.filter((_, i) => i !== deleteIndex);
      deleteIndex = null;
    }
  }
</script>

<section
  class={`
    rounded-lg p-4 space-y-4 transition-colors
    ${isDirty
      ? 'border-2 border-yellow-500 bg-yellow-900/20'
      : 'bg-gray-800 border border-gray-700'
    }
  `}
>
  <div class="flex items-center gap-3">
    <h3 class="text-lg font-medium text-white">Scraping URLs</h3>
    <span class="px-2 py-0.5 bg-gray-700 text-gray-300 text-sm rounded-full">
      {links.length}
    </span>
  </div>

  {#if isDirty}
    <p class="text-yellow-400 text-xs">Unsaved changes</p>
  {/if}

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
        <div class="bg-gray-700/50 rounded p-3 space-y-2">
          <div class="flex items-start gap-2">
            <span
              class={`mt-1.5 w-2.5 h-2.5 rounded-full shrink-0 ${
                urlStatuses?.[url] === true
                  ? 'bg-green-500'
                  : urlStatuses?.[url] === false
                    ? 'bg-red-500'
                    : 'bg-gray-500'
              }`}
              title={
                urlStatuses?.[url] === true
                  ? 'Producing data'
                  : urlStatuses?.[url] === false
                    ? 'No data'
                    : 'Not scraped yet'
              }
            ></span>
            <span class="block text-sm font-mono break-all text-gray-200 leading-normal">
              {url}
            </span>
          </div>

          <div class="flex gap-1">
            <button
              type="button"
              onclick={() => handleMoveUp(index)}
              disabled={index === 0}
              class="w-8 h-8 flex items-center justify-center bg-gray-600 hover:bg-gray-500 disabled:opacity-30 disabled:cursor-not-allowed rounded transition-colors"
              title="Move up"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
              </svg>
            </button>
            <button
              type="button"
              onclick={() => handleMoveDown(index)}
              disabled={index === links.length - 1}
              class="w-8 h-8 flex items-center justify-center bg-gray-600 hover:bg-gray-500 disabled:opacity-30 disabled:cursor-not-allowed rounded transition-colors"
              title="Move down"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <button
              type="button"
              onclick={() => handleOpenEditModal(index)}
              class="w-8 h-8 flex items-center justify-center bg-gray-600 hover:bg-gray-500 rounded transition-colors"
              title="Edit"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
            </button>
            <button
              type="button"
              onclick={() => previewModal.open(url)}
              class="w-8 h-8 flex items-center justify-center bg-gray-600 hover:bg-gray-500 rounded transition-colors"
              title="Preview"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
            </button>
            <button
              type="button"
              onclick={() => handleOpenDeleteDialog(index)}
              class="w-8 h-8 flex items-center justify-center bg-red-600/80 hover:bg-red-600 rounded transition-colors"
              title="Delete"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
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
</section>

<URLModal bind:this={urlModal} onAdd={handleAddUrl} onEdit={handleEditUrl} />
<PreviewModal bind:this={previewModal} />
<ConfirmDialog
  bind:this={confirmDialog}
  title="Delete URL"
  message="Are you sure you want to delete this URL?"
  onConfirm={handleDelete}
/>
