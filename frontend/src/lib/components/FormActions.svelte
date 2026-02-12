<script lang="ts">
  let {
    isDirty,
    loading,
    error,
    successMessage,
    onSubmit,
    onCancel,
  }: {
    isDirty: boolean;
    loading: boolean;
    error: string | null;
    successMessage: string | null;
    onSubmit: () => void;
    onCancel: () => void;
  } = $props();
</script>

<div class="shrink-0 border-t border-gray-700/50 px-3 py-2.5 space-y-2">
  {#if error}
    <div class="bg-red-900/30 border border-red-700/50 rounded px-2.5 py-1.5 text-red-300 text-xs">
      {error}
    </div>
  {/if}

  {#if successMessage}
    <div class="bg-green-900/30 border border-green-700/50 rounded px-2.5 py-1.5 text-green-300 text-xs flex items-center gap-2">
      <svg class="w-3.5 h-3.5 shrink-0" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
          clip-rule="evenodd"
        />
      </svg>
      {successMessage}
    </div>
  {/if}

  <div class="flex gap-2 justify-end">
    <button
      onclick={onCancel}
      disabled={!isDirty || loading}
      class="px-3 py-1.5 text-sm rounded bg-gray-700 hover:bg-gray-600 text-gray-200 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
    >
      Cancel
    </button>
    <button
      onclick={onSubmit}
      disabled={!isDirty || loading}
      class="px-3 py-1.5 text-sm rounded bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center gap-1.5"
    >
      {#if loading}
        <svg class="animate-spin h-3.5 w-3.5" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        Saving...
      {:else}
        Save
      {/if}
    </button>
  </div>
</div>
