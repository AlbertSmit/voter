<script lang="ts">
  import { router } from "tinro";
  import iam from "../stores/iam";

  const style = {
    wrapper:
      "p-4 w-full h-screen flex flex-col items-center justify-center mx-auto my-auto bg-white dark:bg-gray-900",
    title:
      "text-3xl antialiased font-bold tracking-tight text-gray-800 dark:text-white",
    body: "antialiased text-xs",
    url: "text-blue-600",
    input:
      "border p-2 my-1 focus:ring-indigo-500 focus:border-indigo-500 block sm:text-sm border-gray-300 rounded-md",
    button:
      "bg-green-100 dark:bg-white text-green-500 dark:text-gray-800 px-6 py-2 text-xs antialiased font-medium rounded-md",
  };

  async function requestRoomGeneration(): Promise<void> {
    const d = `localhost:1323`;
    const p = `votevotevotevote.herokuapp.com`;
    const { protocol } = window.location;
    const socketUri = `${protocol}//${
      // @ts-ignore
      import.meta.env.MODE === "development" ? d : p
    }`;

    const response = await fetch(`${socketUri}/api/room`);
    const admin = response.headers.get("X-Super-Admin");
    if (admin === "Absolutely!") {
      iam.set(admin);
    }

    router.goto(`room/${await response.json()}`);
  }
</script>

<main class={style.wrapper}>
  <h1 class={style.title}>v4te</h1>
  <hr />

  <button class={style.button} on:click={requestRoomGeneration}>
    Create Room
  </button>
</main>
