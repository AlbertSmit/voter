<script>
  export let status;

  const style = {
    wrapper: "w-full gap-x-1 grid grid-cols-3 overflow-hidden rounded mb-2",
    base:
      "transition-all px-4 py-2 bg-gray-800 text-white text-xs antialiased font-medium",
  };
</script>

<div class={style.wrapper}>
  <span class={`${status !== "WAITING" && "text-gray-500"} ${style.base}`}>
    Waiting
  </span>
  <span class={`${status !== "VOTING" && "text-gray-500"} ${style.base}`}>
    Voting
  </span>
  <span class={`${status !== "PRESENTING" && "text-gray-500"} ${style.base}`}>
    Presenting
  </span>
</div>
