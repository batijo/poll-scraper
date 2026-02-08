<script lang="ts">
  import { onMount } from 'svelte';
  import { GetConfig, UpdateConfig } from '../../../wailsjs/go/main/App';
  import SettingsSidebar from './SettingsSidebar.svelte';
  import FormActions from './FormActions.svelte';
  import GeneralSettings from './forms/GeneralSettings.svelte';
  import ServerSettings from './forms/ServerSettings.svelte';
  import URLList from './forms/URLList.svelte';
  import ScrapingSettings from './forms/ScrapingSettings.svelte';
  import FilterSettings from './forms/FilterSettings.svelte';
  import OutputSettings from './forms/OutputSettings.svelte';
  import StatusSection from './StatusSection.svelte';
  import type { Config } from '../types/config';
  import type { ScraperData, ScraperState, LogEntry } from '../types/scraper';
  import { createDefaultConfig } from '../types/config';

  type Section = 'settings' | 'scraping' | 'status';

  let {
    displayData = $bindable([]),
    rawScrapedData = $bindable([]),
    formState = $bindable(),
    savedConfig = $bindable(),
    urlStatuses = {},
    scraperState = 'stopped',
    logEntries = [],
    lastError = null
  }: {
    displayData?: ScraperData[];
    rawScrapedData?: ScraperData[];
    formState?: Config;
    savedConfig?: Config;
    urlStatuses?: Record<string, boolean>;
    scraperState?: ScraperState;
    logEntries?: LogEntry[];
    lastError?: string | null;
  } = $props();

  let activeSection: Section = $state('status');
  let initialState: Config = $state(formState ? { ...formState } : createDefaultConfig());
  let loading = $state(false);
  let error = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  const isDirty = $derived(JSON.stringify(formState) !== JSON.stringify(initialState));

  onMount(async () => {
    try {
      const config = await GetConfig();
      formState = { ...config };
      initialState = JSON.parse(JSON.stringify(config));
      if (savedConfig) {
        savedConfig = JSON.parse(JSON.stringify(config));
      }
    } catch (e) {
      error = `Failed to load configuration: ${e}`;
    }
  });

  async function handleSubmit() {
    loading = true;
    error = null;
    successMessage = null;

    try {
      await UpdateConfig(formState);
      initialState = JSON.parse(JSON.stringify(formState));
      if (savedConfig) {
        savedConfig = JSON.parse(JSON.stringify(formState));
      }
      successMessage = 'Configuration saved successfully';
      setTimeout(() => {
        successMessage = null;
      }, 3000);
    } catch (e) {
      error = `Failed to save configuration: ${e}`;
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    formState = JSON.parse(JSON.stringify(initialState));
  }

  export function addToFilterLines(indices: number[]) {
    formState.filter_lines = [...formState.filter_lines, ...indices];
  }
</script>

<div class="h-full flex flex-col max-w-full">
  <SettingsSidebar bind:activeSection />

  <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
    <main class="flex-1 min-h-0 p-4 {activeSection === 'status' ? 'flex flex-col overflow-hidden' : 'overflow-y-auto'}">
      {#if activeSection === 'settings'}
        <div class="space-y-4 w-full">
          <GeneralSettings bind:config={formState} initialConfig={initialState} />
          <ServerSettings bind:config={formState} initialConfig={initialState} />
          <OutputSettings bind:config={formState} initialConfig={initialState} />
        </div>
      {:else if activeSection === 'scraping'}
        <div class="space-y-4 w-full">
          <URLList bind:links={formState.links} initialLinks={initialState.links} {urlStatuses} />
          <ScrapingSettings bind:config={formState} initialConfig={initialState} />
          <FilterSettings bind:config={formState} initialConfig={initialState} bind:rawScrapedData {scraperState} />
        </div>
      {:else if activeSection === 'status'}
        <StatusSection config={savedConfig} {displayData} {scraperState} {logEntries} {lastError} {urlStatuses} bind:rawScrapedData />
      {/if}
    </main>

    {#if activeSection !== 'status'}
      <div class="flex-shrink-0">
        <FormActions
          {isDirty}
          {loading}
          {error}
          {successMessage}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
        />
      </div>
    {/if}
  </div>
</div>
