<script lang="ts">
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import DataDisplay from '$lib/components/DataDisplay.svelte';
</script>

<div class="flex gap-3 p-3 min-h-screen min-w-[540px] bg-gray-950">
  <div class="w-[360px] flex-shrink-0">
    <SettingsPanel />
  </div>
  <div class="w-40 flex-shrink-0">
    <DataDisplay />
  </div>
</div>
