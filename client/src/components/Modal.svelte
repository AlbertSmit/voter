<script>
  const style = {
    cloak: "z-50 flex inset-0 absolute bg-gray-900 bg-opacity-70",
    card: "p-4 rounded-md bg-white mx-auto my-auto",
  };
</script>

<div class={style.cloak}>
  <div class={style.card}>
    <slot />
  </div>
</div>
