<script lang="ts">
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
</script>

<div class="flex lg:flex-row flex-col gap-6 p-6 min-h-screen bg-gray-950">
  <div class="lg:w-1/2 w-full">
    <SettingsPanel />
  </div>
  <div class="lg:w-1/2 w-full">
    <DataDisplay />
  </div>
</div>
