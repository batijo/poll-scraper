<script lang="ts">
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
</script>

<div class="flex flex-col gap-4 p-4 min-h-screen bg-gray-950">
  <div class="w-full max-w-md">
    <SettingsPanel />
  </div>
  <div class="w-full max-w-xs">
    <DataDisplay />
  </div>
</div>
