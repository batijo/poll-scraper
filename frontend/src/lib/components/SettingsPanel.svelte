<script lang="ts">
  import { onMount } from 'svelte';
  import { GetConfig, UpdateConfig } from '../../../wailsjs/go/main/App';
  import SettingsSidebar from './SettingsSidebar.svelte';
  import FormActions from './FormActions.svelte';
  import GeneralSettings from './forms/GeneralSettings.svelte';
  import ServerSettings from './forms/ServerSettings.svelte';
  import URLList from './forms/URLList.svelte';
  import ScrapingSettings from './forms/ScrapingSettings.svelte';
  import type { Config } from '../types/config';
  import { createDefaultConfig } from '../types/config';

  type Section = 'settings' | 'scraping' | 'output';

  let activeSection: Section = $state('settings');
  let formState: Config = $state(createDefaultConfig());
  let initialState: Config = $state(createDefaultConfig());
  let loading = $state(false);
  let error = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  const isDirty = $derived(JSON.stringify(formState) !== JSON.stringify(initialState));

  onMount(async () => {
    try {
      const config = await GetConfig();
      formState = { ...config };
      initialState = structuredClone(config);
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
      initialState = structuredClone(formState);
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
    formState = structuredClone(initialState);
  }
</script>

<div class="h-full flex">
  <SettingsSidebar bind:activeSection />

  <div class="flex-1 flex flex-col">
    <main class="flex-1 overflow-y-auto p-6">
      {#if activeSection === 'settings'}
        <div class="space-y-8 max-w-2xl">
          <GeneralSettings bind:config={formState} initialConfig={initialState} />
          <ServerSettings bind:config={formState} initialConfig={initialState} />
        </div>
      {:else if activeSection === 'scraping'}
        <div class="text-gray-400">
          <p>URL management and scraping settings will appear here</p>
        </div>
      {:else if activeSection === 'output'}
        <div class="text-gray-400">
          <p>Output settings and status will appear here</p>
        </div>
      {/if}
    </main>

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
