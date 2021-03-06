<script lang="ts">
  import { infinity } from "../directives/infinity";

  const style = {
    wrapper:
      "p-4 h-screen w-full flex flex-col items-start mx-auto overflow-auto text-gray-800",
    h1: "text-8xl w-full flex justify-between antialiased font-light",
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

  let selected: number;
  const select = (id: number) => {
    selected = id;
  };
</script>

<div id="wrapper" use:infinity class={style.wrapper}>
  <div id="content" class="w-full">
    <div id="block">
      {#each data as name, index}
        <h1 on:click={() => select(index)} class={style.h1}>
          {name} <span class="opacity-10">{index}</span>
        </h1>
      {/each}
    </div>
    {#each data.concat(...data, ...data) as name, index}
      <h1 on:click={() => select(data.length + index)} class={style.h1}>
        {name} <span class="opacity-10">{data.length + index}</span>
      </h1>
    {/each}
  </div>
</div>

<style>
  #wrapper {
    background: #ffa50047;
  }

  #wrapper::-webkit-scrollbar {
    display: none;
  }

  #wrapper {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
