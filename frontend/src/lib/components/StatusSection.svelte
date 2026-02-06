<script lang="ts">
  import type { Config } from '../types/config';

  let { config }: { config: Config } = $props();

  const urlCount = $derived(config.links.length);
  const serverAddress = $derived(`${config.ip}:${config.port}`);
</script>

<section class="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
  <h2 class="text-lg font-semibold mb-4 text-white">Scraper Status</h2>

  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <span class="text-sm text-gray-400">Status</span>
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
        <span class="text-sm text-white">Running</span>
      </div>
    </div>

    <div class="flex items-center justify-between">
      <span class="text-sm text-gray-400">Active URLs</span>
      <span class="text-sm text-white">{urlCount} configured</span>
    </div>

    <div class="flex items-center justify-between">
      <span class="text-sm text-gray-400">Server</span>
      <span class="text-sm text-white font-mono">{serverAddress}</span>
    </div>

    <div class="flex items-center justify-between">
      <span class="text-sm text-gray-400">Update Interval</span>
      <span class="text-sm text-white">{config.update_interval}ms</span>
    </div>
  </div>
</section>
