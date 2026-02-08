<script lang="ts">
  let dialog: HTMLDialogElement;

  let {
    title,
    message,
    confirmText = 'Delete',
    confirmClass = 'bg-red-600 hover:bg-red-500',
    onConfirm
  }: {
    title: string;
    message: string;
    confirmText?: string;
    confirmClass?: string;
    onConfirm: () => void;
  } = $props();

  export function open() {
    dialog.showModal();
  }

  export function close() {
    dialog.close();
  }

  function handleConfirm() {
    onConfirm();
    close();
  }

  function handleCancel() {
    close();
  }
</script>

<dialog
  bind:this={dialog}
  class="bg-gray-800 text-white rounded-lg shadow-2xl max-w-md w-full p-6 backdrop:bg-black/75"
  style="margin: auto;"
>
  <h2 class="text-xl font-semibold mb-4">{title}</h2>

  <p class="text-gray-300 mb-6">{message}</p>

  <div class="flex gap-3">
    <button
      type="button"
      onclick={handleCancel}
      class="flex-1 px-4 py-2 bg-gray-600 hover:bg-gray-500 rounded-md transition-colors"
    >
      Cancel
    </button>
    <button
      type="button"
      onclick={handleConfirm}
      class="flex-1 px-4 py-2 {confirmClass} rounded-md transition-colors"
    >
      {confirmText}
    </button>
  </div>
</dialog>
