<script lang="ts">
  let { effectiveStatus = 'stopped' }: { effectiveStatus: 'ok' | 'stopped' | 'error' } = $props();

  const dotColor = $derived(
    effectiveStatus === 'error' ? 'bg-red-500' :
    effectiveStatus === 'stopped' ? 'bg-gray-500' :
    'bg-green-500'
  );

  const label = $derived(
    effectiveStatus === 'error' ? 'Error' :
    effectiveStatus === 'stopped' ? 'Stopped' :
    'Scraping'
  );
</script>

<div class="inline-flex items-center gap-2 text-sm font-medium text-gray-200">
  <span class="{dotColor} inline-flex rounded-full h-3 w-3"></span>
  {label}
</div>
