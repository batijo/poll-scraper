<script lang="ts">
  type Section = 'settings' | 'scraping' | 'output';

  let { activeSection = $bindable() }: { activeSection: Section } = $props();

  const sections: { id: Section; label: string }[] = [
    { id: 'settings', label: 'Settings' },
    { id: 'scraping', label: 'Scraping' },
    { id: 'output', label: 'Output' },
  ];
</script>

<aside class="w-56 bg-gray-800 border-r border-gray-700 h-full flex flex-col">
  <div class="p-6 border-b border-gray-700">
    <h1 class="text-xl font-bold text-white">Poll Scraper</h1>
  </div>

  <nav class="flex-1 p-4">
    <ul class="space-y-1">
      {#each sections as section}
        <li>
          <button
            onclick={() => (activeSection = section.id)}
            class="w-full text-left px-4 py-2 rounded transition-colors {activeSection ===
            section.id
              ? 'bg-gray-700 border-l-2 border-blue-500'
              : 'hover:bg-gray-700/50'}"
          >
            {section.label}
          </button>
        </li>
      {/each}
    </ul>
  </nav>
</aside>
