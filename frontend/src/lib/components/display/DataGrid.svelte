<script lang="ts">
  import type { ScraperData } from '$lib/types/scraper';
  import DataCard from './DataCard.svelte';

  let { data }: { data: ScraperData[] } = $props();
</script>

<div class="space-y-2">
  {#each data as item}
    <DataCard name={item.name} value={item.value} />
  {/each}
</div>
