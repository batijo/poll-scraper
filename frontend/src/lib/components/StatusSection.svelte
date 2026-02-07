<script lang="ts">
  import { tick } from 'svelte';
  import type { Config } from '../types/config';
  import type { ScraperState, LogEntry } from '../types/scraper';

  let {
    config,
    scraperState = 'idle',
    logEntries = [],
    lastError = null,
    urlStatuses = {}
  }: {
    config: Config;
    scraperState?: ScraperState;
    logEntries?: LogEntry[];
    lastError?: string | null;
    urlStatuses?: Record<string, boolean>;
  } = $props();

  let logContainer: HTMLDivElement;
  let autoScroll = $state(true);

  const urlCount = $derived(config.links.length);
  const urlsWithData = $derived(
    Object.values(urlStatuses).filter((v) => v === true).length
  );
  const urlsNoData = $derived(
    Object.values(urlStatuses).filter((v) => v === false).length
  );
  const serverAddress = $derived(`${config.ip}:${config.port}`);

  const stateLabel = $derived(
    scraperState === 'scraping' ? 'Scraping'
    : scraperState === 'error' ? 'Error'
    : 'Idle'
  );
  const stateDotColor = $derived(
    scraperState === 'error' ? 'bg-red-500'
    : scraperState === 'scraping' ? 'bg-blue-500 animate-pulse'
    : 'bg-green-500'
  );

  const outputTargets = $derived(
    [
      config.write_to_csv ? `CSV (${config.csv_path})` : null,
      config.write_to_txt ? `TXT (${config.txt_path})` : null,
    ].filter(Boolean)
  );

  // Auto-scroll log when new entries arrive
  $effect(() => {
    if (logEntries.length && autoScroll && logContainer) {
      tick().then(() => {
        logContainer.scrollTop = logContainer.scrollHeight;
      });
    }
  });

  function handleLogScroll() {
    if (!logContainer) return;
    const { scrollTop, scrollHeight, clientHeight } = logContainer;
    autoScroll = scrollHeight - scrollTop - clientHeight < 30;
  }

  function levelColor(level: string): string {
    switch (level) {
      case 'ERROR': return 'text-red-400';
      case 'WARN': return 'text-yellow-400';
      case 'DEBUG': return 'text-gray-500';
      default: return 'text-blue-400';
    }
  }
</script>

<div class="space-y-4">
  <section class="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
    <h2 class="text-lg font-semibold mb-4 text-white">Scraper Status</h2>

    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-400">State</span>
        <div class="flex items-center gap-2">
          <span class="w-2 h-2 {stateDotColor} rounded-full"></span>
          <span class="text-sm text-white">{stateLabel}</span>
        </div>
      </div>

      {#if lastError}
        <div class="flex items-start justify-between gap-2">
          <span class="text-sm text-gray-400">Last Error</span>
          <span class="text-sm text-red-400 text-right max-w-[60%] break-words">{lastError}</span>
        </div>
      {/if}

      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-400">URLs</span>
        <span class="text-sm text-white">
          {urlCount} configured
          {#if urlsWithData > 0 || urlsNoData > 0}
            <span class="text-gray-400">
              ({urlsWithData} <span class="text-green-400">ok</span>{#if urlsNoData > 0}, {urlsNoData} <span class="text-red-400">fail</span>{/if})
            </span>
          {/if}
        </span>
      </div>

      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-400">Server</span>
        <span class="text-sm text-white">
          {#if config.enable_server}
            <span class="font-mono">{serverAddress}</span>
            <span class="text-green-400 text-xs ml-1">On</span>
          {:else}
            <span class="text-gray-500">Off</span>
          {/if}
        </span>
      </div>

      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-400">Update Interval</span>
        <span class="text-sm text-white">{config.update_interval}ms</span>
      </div>

      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-400">Output</span>
        <span class="text-sm text-white">
          {#if outputTargets.length > 0}
            {outputTargets.join(', ')}
          {:else}
            None
          {/if}
        </span>
      </div>

      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-400">Debug</span>
        <span class="text-sm text-white">{config.debug ? 'On' : 'Off'}</span>
      </div>
    </div>
  </section>

  <section class="bg-gray-800/50 rounded-lg border border-gray-700 flex flex-col" style="height: calc(100vh - 380px); min-height: 200px;">
    <div class="flex items-center justify-between px-4 py-2 border-b border-gray-700 flex-shrink-0">
      <h2 class="text-lg font-semibold text-white">Logs</h2>
      <span class="text-xs text-gray-500">{logEntries.length} entries</span>
    </div>

    <div
      bind:this={logContainer}
      onscroll={handleLogScroll}
      class="flex-1 overflow-y-auto p-3 font-mono text-xs leading-relaxed min-h-0"
    >
      {#if logEntries.length === 0}
        <p class="text-gray-500 text-center py-8">No log entries yet</p>
      {:else}
        {#each logEntries as entry}
          <div class="flex gap-2 hover:bg-gray-700/30 px-1 rounded">
            <span class="text-gray-600 flex-shrink-0">{entry.time}</span>
            <span class="{levelColor(entry.level)} flex-shrink-0 w-12">{entry.level}</span>
            <span class="text-gray-300 break-all">{entry.message}</span>
          </div>
        {/each}
      {/if}
    </div>
  </section>
</div>
