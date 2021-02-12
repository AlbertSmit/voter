<script lang="ts">
  import { router } from "tinro";

  let roomUri: string;

  const style = {
    wrapper:
      "p-4 w-full h-screen flex flex-col items-center justify-center mx-auto my-auto",
    title: "text-3xl antialiased font-bold tracking-tight",
    body: "antialiased text-xs",
    url: "text-blue-600",
    input:
      "border p-2 my-1 focus:ring-indigo-500 focus:border-indigo-500 block sm:text-sm border-gray-300 rounded-md",
    button:
      "bg-green-100 px-6 py-2 text-xs antialiased font-medium rounded-md text-green-500",
  };

  async function requestRoomGeneration(): Promise<void> {
    const d = `localhost:1323`;
    const p = `votevotevotevote.herokuapp.com`;
    const { protocol } = window.location;
    const socketUri = `${protocol}//${
      // @ts-ignore
      import.meta.env.MODE === "development" ? d : p
    }`;

    const response = await fetch(`${socketUri}/room`);
    roomUri = await response.json();
  }

  function redirectToRoom(): void {
    router.goto(`room/${roomUri}`);
  }
</script>

<main class={style.wrapper}>
  <h1 class={style.title}>v4te</h1>
  <hr />

  {#if !roomUri}
    <button class={style.button} on:click={requestRoomGeneration}>
      Create Room
    </button>
  {/if}

  {#if roomUri}
    <button class={style.button} on:click={redirectToRoom}>
      Go to the room
    </button>
  {/if}
</main>
