<script lang="ts">
  import { onMount } from 'svelte';
  import { Greet, GetConfig } from '../../wailsjs/go/main/App';

  let greetResult = $state('');
  let nameInput = $state('World');
  let configData = $state<any>(null);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      configData = await GetConfig();
    } catch (e) {
      error = `Failed to load config: ${e}`;
    }
  });

  async function handleGreet() {
    try {
      greetResult = await Greet(nameInput);
      error = null;
    } catch (e) {
      error = `Greet failed: ${e}`;
    }
  }
</script>

<div class="max-w-2xl mx-auto space-y-8">
  <h1 class="text-3xl font-bold">Poll Scraper</h1>

  <section class="bg-gray-800 rounded-lg p-6 space-y-4">
    <h2 class="text-xl font-semibold">Go Communication Test</h2>
    <div class="flex gap-3">
      <input
        type="text"
        bind:value={nameInput}
        placeholder="Enter name"
        class="bg-gray-700 text-white rounded px-3 py-2 flex-1"
      />
      <button
        onclick={handleGreet}
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
      >
        Greet
      </button>
    </div>
    {#if greetResult}
      <p class="text-green-400">{greetResult}</p>
    {/if}
  </section>

  <section class="bg-gray-800 rounded-lg p-6 space-y-4">
    <h2 class="text-xl font-semibold">Current Configuration</h2>
    {#if configData}
      <dl class="grid grid-cols-2 gap-2 text-sm">
        <dt class="text-gray-400">Dataset Name</dt>
        <dd>{configData.dataset_name || 'N/A'}</dd>
        <dt class="text-gray-400">Port</dt>
        <dd>{configData.port}</dd>
        <dt class="text-gray-400">Update Interval</dt>
        <dd>{configData.update_interval}ms</dd>
        <dt class="text-gray-400">Links</dt>
        <dd>{configData.links?.length || 0} configured</dd>
        <dt class="text-gray-400">Debug</dt>
        <dd>{configData.debug ? 'Enabled' : 'Disabled'}</dd>
      </dl>
    {:else if !error}
      <p class="text-gray-400">Loading configuration...</p>
    {/if}
  </section>

  {#if error}
    <div class="bg-red-900/50 border border-red-700 rounded-lg p-4">
      <p class="text-red-400">{error}</p>
    </div>
  {/if}
</div>
