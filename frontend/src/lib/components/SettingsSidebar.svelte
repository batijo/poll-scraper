<script lang="ts">
  type Section = 'settings' | 'scraping' | 'status';

  let { activeSection = $bindable() }: { activeSection: Section } = $props();

  const sections: { id: Section; label: string }[] = [
    { id: 'status', label: 'Status' },
    { id: 'settings', label: 'Settings' },
    { id: 'scraping', label: 'Scraping' },
  ];
</script>

<nav class="shrink-0 border-b border-gray-700/50 px-2">
  <ul class="flex">
    {#each sections as section}
      <li>
        <button
          onclick={() => (activeSection = section.id)}
          class="px-4 py-2.5 text-sm font-medium transition-colors relative {activeSection === section.id
            ? 'text-white'
            : 'text-gray-400 hover:text-gray-200'}"
        >
          {section.label}
          {#if activeSection === section.id}
            <span class="absolute bottom-0 left-2 right-2 h-0.5 bg-blue-500 rounded-full"></span>
          {/if}
        </button>
      </li>
    {/each}
  </ul>
</nav>
