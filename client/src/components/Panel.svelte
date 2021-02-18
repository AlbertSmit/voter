<script lang="ts">
  let open: boolean = false;

  const style = {
    wrapper:
      "absolute space-x-2 bg-gray-800 text-white rounded-full flex items-center justify-center right-4 bottom-4 shadow-xl w-auto px-6 h-14",
  };
</script>

<div
  class={style.wrapper}
  on:mouseenter={() => (open = true)}
  on:mouseleave={() => (open = false)}
>
  <svg
    class={`mx-2 h-4 w-4 block text-white transition-all transform ${
      open && "rotate-180"
    }`}
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 20 20"
    fill="currentColor"
  >
    <path
      fill-rule="evenodd"
      d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
      clip-rule="evenodd"
    />
  </svg>
  {#if open}
    <slot />
  {/if}
</div>
