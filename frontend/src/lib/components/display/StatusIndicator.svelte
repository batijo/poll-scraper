<script lang="ts">
  import type { ScraperState } from '$lib/types/scraper';

  let { status }: { status: ScraperState } = $props();

  const statusConfig = {
    idle: { color: 'bg-gray-700 text-gray-200', icon: '○', label: 'Ready' },
    scraping: { color: 'bg-blue-900 text-blue-100', icon: '⟳', label: 'Scraping...' },
    error: { color: 'bg-red-900 text-red-100', icon: '✕', label: 'Error' },
  };

  const config = $derived(statusConfig[status]);
</script>

<div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium {config.color}">
  <span class={status === 'scraping' ? 'animate-spin' : ''}>
    {config.icon}
  </span>
  {config.label}
</div>
