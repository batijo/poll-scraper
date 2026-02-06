<script lang="ts">
  type Section = 'settings' | 'scraping' | 'output';

  let { activeSection = $bindable() }: { activeSection: Section } = $props();

  const sections: { id: Section; label: string }[] = [
    { id: 'settings', label: 'Settings' },
    { id: 'scraping', label: 'Scraping' },
    { id: 'output', label: 'Output' },
  ];
</script>

<div class="bg-gray-800 border-b border-gray-700">
  <div class="px-4 py-3 border-b border-gray-700">
    <h1 class="text-lg font-bold text-white">Poll Scraper</h1>
  </div>

  <nav class="px-4">
    <ul class="flex gap-1">
      {#each sections as section}
        <li>
          <button
            onclick={() => (activeSection = section.id)}
            class="px-4 py-2 text-sm font-medium transition-colors {activeSection === section.id
              ? 'text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-gray-200'}"
          >
            {section.label}
          </button>
        </li>
      {/each}
    </ul>
  </nav>
</div>
