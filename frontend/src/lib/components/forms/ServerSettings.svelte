<script lang="ts">
  import type { Config } from '../../types/config';

  let { config, initialConfig }: { config: Config; initialConfig: Config } = $props();

  function isFieldDirty(field: keyof Config): boolean {
    return JSON.stringify(config[field]) !== JSON.stringify(initialConfig[field]);
  }

  let domainsText = $derived(config.domains.join('\n'));

  function updateDomains(event: Event) {
    const target = event.target as HTMLTextAreaElement;
    const text = target.value;
    config.domains = text
      .split('\n')
      .map((line) => line.trim())
      .filter((line) => line.length > 0);
  }
</script>

<section class="space-y-4">
  <h2 class="text-lg font-semibold mb-4 text-white">Server Settings</h2>

  <div>
    <label for="port" class="block text-sm font-medium text-gray-300 mb-1">
      Port
    </label>
    <input
      id="port"
      type="number"
      min="1"
      max="65535"
      required
      bind:value={config.port}
      class={`
        w-full px-3 py-2 rounded text-white
        transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
        invalid:border-red-500 invalid:bg-red-900/20
        ${
          isFieldDirty('port')
            ? 'border-2 border-yellow-500 bg-yellow-900/20'
            : 'border border-gray-600 bg-gray-700'
        }
      `}
    />
    {#if config.port < 1 || config.port > 65535}
      <p class="text-red-400 text-xs mt-1">Port must be 1-65535</p>
    {:else if isFieldDirty('port')}
      <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
    {/if}
  </div>

  <div>
    <label for="ip" class="block text-sm font-medium text-gray-300 mb-1">
      IP Address
    </label>
    <input
      id="ip"
      type="text"
      bind:value={config.ip}
      placeholder="localhost"
      class={`
        w-full px-3 py-2 rounded text-white
        transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
        ${
          isFieldDirty('ip')
            ? 'border-2 border-yellow-500 bg-yellow-900/20'
            : 'border border-gray-600 bg-gray-700'
        }
      `}
    />
    {#if isFieldDirty('ip')}
      <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
    {/if}
  </div>

  <div>
    <label for="domains" class="block text-sm font-medium text-gray-300 mb-1">
      CORS Domains (one per line)
    </label>
    <textarea
      id="domains"
      rows="3"
      value={domainsText}
      oninput={updateDomains}
      placeholder="https://example.com"
      class={`
        w-full px-3 py-2 rounded text-white
        transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none
        resize-vertical
        ${
          isFieldDirty('domains')
            ? 'border-2 border-yellow-500 bg-yellow-900/20'
            : 'border border-gray-600 bg-gray-700'
        }
      `}
    ></textarea>
    {#if isFieldDirty('domains')}
      <p class="text-yellow-400 text-xs mt-1">Unsaved change</p>
    {/if}
  </div>
</section>
