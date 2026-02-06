<script>
  import '../app.css';
  let { children } = $props();
</script>

<main class="h-screen bg-gray-900 text-white">
  {@render children()}
</main>
