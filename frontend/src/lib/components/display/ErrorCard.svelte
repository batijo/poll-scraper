<script lang="ts">
  let { message, reset }: { message: string; reset: () => void } = $props();
</script>

<div class="bg-red-900/20 border border-red-700 rounded-lg p-4">
  <p class="font-semibold text-red-100 mb-2">Error Loading Data</p>
  <p class="text-sm text-red-100/90 mb-4">{message}</p>
  <button
    onclick={reset}
    class="text-xs bg-red-700 hover:bg-red-600 text-white px-3 py-1 rounded transition-colors"
  >
    Try Again
  </button>
</div>
