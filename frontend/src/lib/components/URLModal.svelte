<script lang="ts">
  let dialog: HTMLDialogElement;
  let urlInput = $state('');
  let urlError = $state<string | null>(null);
  let editIndex = $state<number | null>(null);

  let { onAdd, onEdit }: { onAdd: (url: string) => void; onEdit?: (index: number, url: string) => void } = $props();

  const isEditing = $derived(editIndex !== null);

  export function open(existingUrl?: string, index?: number) {
    if (existingUrl !== undefined && index !== undefined) {
      urlInput = existingUrl;
      editIndex = index;
    }
    dialog.showModal();
  }

  export function close() {
    dialog.close();
    reset();
  }

  function reset() {
    urlInput = '';
    urlError = null;
    editIndex = null;
  }

  function validateURL(url: string): boolean {
    if (!url.trim()) {
      urlError = 'URL cannot be empty';
      return false;
    }
    try {
      new URL(url);
      urlError = null;
      return true;
    } catch {
      urlError = 'Invalid URL format';
      return false;
    }
  }

  function handleSubmit() {
    if (validateURL(urlInput)) {
      if (isEditing && onEdit) {
        onEdit(editIndex!, urlInput);
      } else {
        onAdd(urlInput);
      }
      close();
    }
  }

  function handleCancel() {
    close();
  }

  $effect(() => {
    if (urlInput) {
      validateURL(urlInput);
    }
  });
</script>

<dialog
  bind:this={dialog}
  class="bg-gray-800 text-white rounded-lg shadow-2xl max-w-md w-full p-6 backdrop:bg-black/75"
  style="margin: auto;"
  onclose={reset}
>
  <h2 class="text-xl font-semibold mb-4">{isEditing ? 'Edit' : 'Add'} Scraping URL</h2>

  <div class="space-y-4">
    <div>
      <input
        type="text"
        bind:value={urlInput}
        placeholder="https://example.com/data"
        class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      {#if urlError}
        <p class="text-red-400 text-sm mt-1">{urlError}</p>
      {/if}
    </div>

    <div class="flex gap-3 pt-4">
      <button
        type="button"
        onclick={handleCancel}
        class="flex-1 px-4 py-2 bg-gray-600 hover:bg-gray-500 rounded-md transition-colors"
      >
        Cancel
      </button>
      <button
        type="button"
        onclick={handleSubmit}
        disabled={!!urlError || !urlInput.trim()}
        class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed rounded-md transition-colors"
      >
        {isEditing ? 'Save' : 'Add'}
      </button>
    </div>
  </div>
</dialog>
