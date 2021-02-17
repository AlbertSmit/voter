<script>
  const style = {
    button:
      "bg-gray-800 text-white dark:text-gray-800 dark:bg-white px-6 py-2 text-xs antialiased font-medium rounded-md whitespace-nowrap",
  };
</script>

<button on:click on:submit class={style.button}> <slot /> </button>
