<script lang="ts">
  import { onMount } from "svelte";
  import { infinity } from "../directives/infinity";

  export let data: any[] = [];
  export let callback: (id: number | null) => void;

  let wrapper: number;
  let block: number;

  let multiplier: number[] = [1];
  onMount(() => {
    multiplier = Array.from(new Array(Math.ceil(wrapper / block) * 2));
  });

  const select = (id: number | null, idx: number = 1): void => {
    const w = document.getElementById("wrapper");

    const travel = block * idx;
    const split = block / data.length;
    const amount = split * id!;

    if (id != null)
      (w as HTMLDivElement).scrollTo({
        top: travel + amount,
        behavior: "smooth",
      });

    callback(id);
  };

  const style = {
    wrapper:
      "relative px-4 h-screen bg-spruce w-full flex flex-col items-start mx-auto overflow-auto text-gray-900",
    h1: "text-8xl w-full flex justify-between antialiased font-light",
    dot: "ml-6 mt-5 h-6 w-6 rounded-full bg-green-300 inline-block",
    center: "flex items-center",
    num: "opacity-10",
  };
</script>

<div
  id="wrapper"
  use:infinity
  class={style.wrapper}
  bind:offsetHeight={wrapper}
>
  <div class="w-full">
    {#each multiplier as _}
      <div bind:offsetHeight={block}>
        {#each data as name, index}
          <h1 on:click={() => select(index)} class={style.h1}>
            <span class={style.center}>
              <span>{name}</span>
              <span class={style.dot} />
            </span>
            <span class={style.num}>{index}</span>
          </h1>
        {/each}
      </div>
    {/each}
  </div>
  <div class="fixed flex flex-col text-xs bg-white top-0 right-0">
    <span>wrapper: {wrapper}</span>
    <span>block: {block}</span>
  </div>
</div>

<style>
  ::selection {
    background: yellowgreen;
  }

  #wrapper::-webkit-scrollbar {
    display: none;
  }

  #wrapper {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
