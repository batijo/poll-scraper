<script>
  import '../app.css';
  let { children } = $props();
</script>

<main class="min-h-screen bg-gray-900 text-white p-6">
  {@render children()}
</main>
