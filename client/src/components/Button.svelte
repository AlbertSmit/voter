<script>
  const style = {
    button:
      "bg-gray-800 text-white dark:text-white dark:bg-gray-800 px-6 py-2 text-xs antialiased font-medium rounded-md whitespace-nowrap",
  };
</script>

<button on:click on:submit class={style.button}> <slot /> </button>
