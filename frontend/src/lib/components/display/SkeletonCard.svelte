<script lang="ts">
  let { count = 3 }: { count?: number } = $props();
</script>

{#each Array(count) as _, i (i)}
  <div class="bg-gray-800 border border-gray-700 rounded-lg p-4 animate-pulse">
    <div class="h-4 bg-gray-700 rounded mb-3 w-1/3"></div>
    <div class="space-y-2">
      <div class="h-6 bg-gray-700 rounded"></div>
      <div class="h-6 bg-gray-700 rounded w-5/6"></div>
    </div>
  </div>
{/each}
