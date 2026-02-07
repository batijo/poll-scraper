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
  import StatusSection from './StatusSection.svelte';
  import OutputSettings from './forms/OutputSettings.svelte';
  import type { Config } from '../types/config';
  import type { ScraperData } from '../types/scraper';
  import { createDefaultConfig } from '../types/config';

  type Section = 'settings' | 'scraping' | 'output';

  let {
    displayData = [],
    formState = $bindable()
  }: {
    displayData?: ScraperData[];
    formState?: Config;
  } = $props();

  let activeSection: Section = $state('settings');
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
    <main class="flex-1 overflow-y-auto p-4">
      {#if activeSection === 'settings'}
        <div class="space-y-4 w-full">
          <GeneralSettings bind:config={formState} initialConfig={initialState} />
          <ServerSettings bind:config={formState} initialConfig={initialState} />
        </div>
      {:else if activeSection === 'scraping'}
        <div class="space-y-4 w-full">
          <URLList bind:links={formState.links} />
          <ScrapingSettings bind:config={formState} initialConfig={initialState} />
          <FilterSettings bind:config={formState} initialConfig={initialState} {displayData} />
        </div>
      {:else if activeSection === 'output'}
        <div class="space-y-4 w-full">
          <StatusSection config={formState} />
          <OutputSettings bind:config={formState} initialConfig={initialState} />
        </div>
      {/if}
    </main>

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
  </div>
</div>
