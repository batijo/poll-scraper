<script lang="ts">
  let { time }: { time: string } = $props();
</script>

<p class="text-sm text-gray-400">
  Last updated: <time datetime={time}>{time}</time>
</p>
