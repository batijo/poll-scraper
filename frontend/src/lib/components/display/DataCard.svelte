<script lang="ts">
  let { name, value }: { name: string; value: string } = $props();
</script>

<div class="bg-gray-800 border border-gray-700 rounded p-2 hover:border-gray-600 transition-colors">
  <div class="text-xs font-bold text-blue-400 uppercase tracking-wide mb-1 truncate">{name}</div>
  <div class="text-xs text-gray-50 break-words leading-snug">
    {value}
  </div>
</div>
