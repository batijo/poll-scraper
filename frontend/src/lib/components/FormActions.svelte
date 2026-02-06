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

<div class="border-t border-gray-700 bg-gray-800 p-4 space-y-3">
  {#if error}
    <div class="bg-red-900/30 border border-red-700 rounded p-3 text-red-200 text-sm">
      {error}
    </div>
  {/if}

  {#if successMessage}
    <div class="bg-green-900/30 border border-green-700 rounded p-3 text-green-200 text-sm flex items-center gap-2">
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
          clip-rule="evenodd"
        />
      </svg>
      {successMessage}
    </div>
  {/if}

  <div class="flex gap-3 justify-end">
    <button
      onclick={onCancel}
      disabled={!isDirty || loading}
      class="px-4 py-2 rounded bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
    >
      Cancel
    </button>
    <button
      onclick={onSubmit}
      disabled={!isDirty || loading}
      class="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
    >
      {#if loading}
        <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
        Saving...
      {:else}
        Submit
      {/if}
    </button>
  </div>
</div>
