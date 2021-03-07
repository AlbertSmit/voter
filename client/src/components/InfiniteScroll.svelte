<script lang="ts">
  import Multiplier from "../components/Multiplier.svelte";
  import { infinity } from "../directives/infinity";

  export let data: any[] = [];
  export let callback: (index: number) => any;

  const style = {
    wrapper:
      "p-4 h-screen bg-spruce w-full flex flex-col items-start mx-auto overflow-auto text-gray-900",
    h1: "text-8xl w-full flex justify-between antialiased font-light",
    dot: "ml-6 mt-5 h-6 w-6 rounded-full bg-green-300 inline-block",
    center: "flex items-center",
    num: "opacity-10",
  };
</script>

<div id="wrapper" use:infinity class={style.wrapper}>
  <div id="content" class="w-full">
    <Multiplier times={4}>
      {#each data as name, index}
        <h1 on:click={() => callback(index)} class={style.h1}>
          <span class={style.center}>
            <span>{name}</span>
            <span class={style.dot} />
          </span>
          <span class={style.num}>{index}</span>
        </h1>
      {/each}
    </Multiplier>
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
