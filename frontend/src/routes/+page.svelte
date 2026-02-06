<script lang="ts">
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
</script>

<div class="flex gap-3 p-3 min-h-screen min-w-[700px] bg-gray-950 overflow-x-hidden">
  <div class="flex-1 min-w-[480px] max-w-2xl">
    <SettingsPanel />
  </div>
  <div class="w-40 flex-shrink-0">
    <DataDisplay />
  </div>
</div>
