<script lang="ts">
  import { fly } from "svelte/transition";
  import { infinity } from "../directives/infinity";

  const style = {
    wrapper:
      "p-4 h-screen bg-spruce w-full flex flex-col items-start mx-auto overflow-auto text-gray-900",
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

  let selected: number | null;
  const select = (id: number | null, block: number = 1): void => {
    const w = document.getElementById("wrapper");
    const b = document.getElementById("block");

    const travel = b!.offsetHeight * block;
    const split = b!.offsetHeight / data.length;
    const amount = split * id!;

    if (id != null)
      (w as HTMLDivElement).scrollTo({
        top: travel + amount,
        behavior: "smooth",
      });

    selected = id;
  };
</script>

<div id="wrapper" use:infinity class={style.wrapper}>
  {#if selected != null}
    <div
      transition:fly={{ x: 50, duration: 250 }}
      class="absolute right-0 top-0 h-full w-3/4 bg-white shadow-3xl p-4 z-10"
    >
      <button
        on:click={() => select(null)}
        class="absolute right-0 top-0 z-20 text-sm font-medium antialiased py-2 px-4 text-gray-900 m-4"
      >
        <svg
          class=" h-8 w-8"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
      <h1 class="text-8xl antialiased font-light">{data[selected]}</h1>
      <hr class="border-gray-800 border-2 mb-8" />
      <p class="text-sm antialiased">is wie je hebt geselecteerd.</p>
    </div>
  {/if}
  <div id="content" class="w-full">
    <div id="block">
      {#each data as name, index}
        <h1 on:click={() => select(index, 0)} class={style.h1}>
          {name} <span class="opacity-10">{index}</span>
        </h1>
      {/each}
    </div>
    {#each data as name, index}
      <h1 on:click={() => select(index, 1)} class={style.h1}>
        {name} <span class="opacity-10">{index}</span>
      </h1>
    {/each}
    {#each data as name, index}
      <h1 on:click={() => select(index, 2)} class={style.h1}>
        {name} <span class="opacity-10">{index}</span>
      </h1>
    {/each}
    {#each data as name, index}
      <h1 on:click={() => select(index, 3)} class={style.h1}>
        {name} <span class="opacity-10">{index}</span>
      </h1>
    {/each}
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
