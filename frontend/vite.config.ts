import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit()
  ],
  server: {
    // Wails generates bindings in ../wailsjs, outside SvelteKit's default allow list.
    fs: {
      allow: ['..']
    }
  }
});
