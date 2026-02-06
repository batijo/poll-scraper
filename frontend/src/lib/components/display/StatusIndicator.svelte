<script lang="ts">
  import type { ScraperState } from '$lib/types/scraper';

  let { status }: { status: ScraperState } = $props();

  const isError = $derived(status === 'error');
  const dotColor = $derived(isError ? 'bg-red-500' : 'bg-green-500');
</script>

<div class="inline-flex items-center gap-2 text-sm font-medium text-gray-200">
  <span class="relative flex h-3 w-3">
    <span class="{dotColor} absolute inline-flex h-full w-full rounded-full opacity-75 {isError ? '' : 'animate-ping'}"></span>
    <span class="{dotColor} relative inline-flex rounded-full h-3 w-3"></span>
  </span>
  Scraping
</div>
