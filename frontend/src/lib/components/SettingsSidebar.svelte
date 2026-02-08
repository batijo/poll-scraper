<script lang="ts">
  type Section = 'settings' | 'scraping' | 'status';

  let { activeSection = $bindable() }: { activeSection: Section } = $props();

  const sections: { id: Section; label: string }[] = [
    { id: 'status', label: 'Status' },
    { id: 'settings', label: 'Settings' },
    { id: 'scraping', label: 'Scraping' },
  ];
</script>

<div class="bg-gray-800 border-b border-gray-700">
  <nav class="px-4 py-2">
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
