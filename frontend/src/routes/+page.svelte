<script lang="ts">
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
</script>

<div class="flex gap-4 p-4 min-h-screen bg-gray-950 overflow-x-auto">
  <div class="w-80 flex-shrink-0 max-h-screen overflow-y-auto">
    <SettingsPanel />
  </div>
  <div class="flex-shrink-0 w-56">
    <DataDisplay />
  </div>
</div>
