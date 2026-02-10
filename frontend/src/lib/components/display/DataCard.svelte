<script lang="ts">
  let { name, value, index }: { name: string; value: string; index: number } = $props();
</script>

<div class="bg-gray-800 border border-gray-700 rounded p-2 flex items-center gap-2">
  <div class="text-xs font-bold text-gray-500 shrink-0 w-6">{index + 1}.</div>
  <div class="text-xs font-bold text-blue-400 shrink-0 min-w-24">{name}:</div>
  <div class="text-xs text-gray-50 wrap-break-word flex-1">
    {value}
  </div>
</div>
