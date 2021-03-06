<script>
  import { onMount } from "svelte";

  const style = {
    wrapper:
      "p-4 h-screen flex flex-col items-start mx-auto overflow-auto text-gray-800",
    h1: "text-6xl antialiased font-light",
  };
  const data = [
    "Albert",
    "Geert",
    "Jan",
    "Jelle",
    "Pieter",
    "Steve",
    "Igor",
    "Hendrik",
    "Willem",
    "Arend",
  ];

  let wrapper;
  let content;
  let block;

  const scroller = () => {
    if (wrapper.scrollTop + wrapper.offsetHeight > content.offsetHeight) {
      wrapper.scrollTop = block.offsetHeight - wrapper.offsetHeight;
      console.log(wrapper.offsetHeight - block.offsetHeight);
    }
  };

  onMount(() => {
    const x = document.getElementById("wrapper");
    const y = document.getElementById("content");
    const z = document.getElementById("block");

    wrapper = x;
    content = y;
    block = z;
    wrapper.addEventListener("scroll", scroller, false);
  });
</script>

<div id="wrapper" class={style.wrapper}>
  <div id="content">
    <div id="block">
      {#each data as name}
        <h1 class={style.h1}>{name}</h1>
      {/each}
    </div>
    {#each data as name}
      <h1 class={style.h1}>{name}</h1>
    {/each}
    {#each data as name}
      <h1 class={style.h1}>{name}</h1>
    {/each}
    {#each data as name}
      <h1 class={style.h1}>{name}</h1>
    {/each}
  </div>
</div>

<style>
  #wrapper::-webkit-scrollbar {
    display: none;
  }

  #wrapper {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
