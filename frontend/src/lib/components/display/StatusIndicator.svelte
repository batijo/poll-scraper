<script lang="ts">
  import type { ScraperState } from '$lib/types/scraper';

  let { status, hasData }: { status: ScraperState; hasData: boolean } = $props();

  const dotColor = $derived(
    status === 'error' ? 'bg-red-500' :
    hasData ? 'bg-green-500' :
    'bg-blue-500'
  );

  const shouldPulse = $derived(status !== 'error');
</script>

<div class="inline-flex items-center gap-2 text-sm font-medium text-gray-200">
  <span class="relative flex h-3 w-3">
    <span class="{dotColor} absolute inline-flex h-full w-full rounded-full opacity-75 {shouldPulse ? 'animate-ping' : ''}"></span>
    <span class="{dotColor} relative inline-flex rounded-full h-3 w-3"></span>
  </span>
  Scraping
</div>
