<script lang="ts">
  let { name, value }: { name: string; value: string } = $props();
</script>

<div class="bg-gray-800 border border-gray-700 rounded-lg p-4 hover:border-gray-600 transition-colors">
  <dt class="text-sm font-semibold text-gray-300 mb-2">{name}</dt>
  <dd class="text-base text-white break-words whitespace-normal">
    {value}
  </dd>
</div>
