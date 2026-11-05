<script lang="ts">
  import type { Link } from '../types/config';

  let dialog: HTMLDialogElement;
  let urlInput = $state('');
  let usernameInput = $state('');
  let passwordInput = $state('');
  let urlError = $state<string | null>(null);
  let editIndex = $state<number | null>(null);

  let {
    onAdd,
    onEdit
  }: {
    onAdd: (link: Link) => void;
    onEdit?: (index: number, link: Link) => void;
  } = $props();

  const isEditing = $derived(editIndex !== null);
  const credentialsError = $derived(
    Boolean(usernameInput.trim()) !== Boolean(passwordInput)
      ? 'Enter both username and password or leave both empty'
      : null
  );

  export function open(existingLink?: Link, index?: number) {
    reset();
    if (existingLink !== undefined && index !== undefined) {
      urlInput = existingLink.url;
      usernameInput = existingLink.name ?? '';
      passwordInput = existingLink.password ?? '';
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
    usernameInput = '';
    passwordInput = '';
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
    if (!validateURL(urlInput) || credentialsError) {
      return;
    }

    const link: Link = {
      url: urlInput.trim(),
      name: usernameInput.trim() || undefined,
      password: passwordInput || undefined
    };

    if (isEditing && onEdit) {
      onEdit(editIndex!, link);
    } else {
      onAdd(link);
    }
    close();
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
  class="bg-gray-800 text-white rounded-lg shadow-2xl max-w-lg w-full p-6 backdrop:bg-black/75"
  style="margin: auto;"
  onclose={reset}
>
  <h2 class="text-xl font-semibold mb-4">{isEditing ? 'Edit' : 'Add'} Scraping URL</h2>

  <div class="space-y-4">
    <div>
      <label for="scraping-url" class="block text-sm text-gray-300 mb-1">URL</label>
      <textarea
        id="scraping-url"
        bind:value={urlInput}
        placeholder="https://example.com/data"
        rows="3"
        class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 resize-y font-mono text-sm"
      ></textarea>
      {#if urlError}
        <p class="text-red-400 text-sm mt-1">{urlError}</p>
      {/if}
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div>
        <label for="scraping-username" class="block text-sm text-gray-300 mb-1">
          Username <span class="text-gray-500">(optional)</span>
        </label>
        <input
          id="scraping-username"
          type="text"
          bind:value={usernameInput}
          autocomplete="username"
          class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
        />
      </div>

      <div>
        <label for="scraping-password" class="block text-sm text-gray-300 mb-1">
          Password <span class="text-gray-500">(optional)</span>
        </label>
        <input
          id="scraping-password"
          type="password"
          bind:value={passwordInput}
          autocomplete="current-password"
          class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
        />
      </div>
    </div>

    {#if credentialsError}
      <p class="text-red-400 text-sm">{credentialsError}</p>
    {/if}

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
        disabled={!!urlError || !urlInput.trim() || !!credentialsError}
        class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed rounded-md transition-colors"
      >
        {isEditing ? 'Save' : 'Add'}
      </button>
    </div>
  </div>
</dialog>
