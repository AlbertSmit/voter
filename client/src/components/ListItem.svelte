<script lang="ts">
  export let status: any;
  export let user: any;

  const handle = () => console.log(user.name);

  const style = {
    wrapper:
      "p-4 dark:bg-white bg-opacity-10 text-gray-800 dark:text-white rounded-md",
    text: "antialiased text-xs",
  };
</script>

<div on:click={handle} class={style.wrapper}>
  <p class={style.text}>
    {user.name}
    {#if status === "VOTING"}
      (vote)
    {/if}
  </p>
</div>
