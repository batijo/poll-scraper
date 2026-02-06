<script lang="ts">
  import type { ScraperData } from '$lib/types/scraper';
  import DataCard from './DataCard.svelte';

  let { data }: { data: ScraperData[] } = $props();
</script>

<div class="space-y-4">
  {#each data as item (item.name)}
    <DataCard name={item.name} value={item.value} />
  {/each}
</div>
