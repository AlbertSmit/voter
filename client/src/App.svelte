<script lang="ts">
  import { onMount } from "svelte";
  import store from "./stores/websocket";
  let message: any;
  let messages: any = [];

  const style = {
    wrapper: "p-4 flex flex-col items-start justify-center mx-auto my-auto",
    title: "text-3xl antialiased font-bold tracking-tight",
    body: "antialiased text-xs",
    url: "text-blue-600",
    input:
      "border p-2 my-1 focus:ring-indigo-500 focus:border-indigo-500 block sm:text-sm border-gray-300 rounded-md",
    button:
      "bg-green-100 px-6 py-2 text-xs antialiased font-medium rounded-md text-green-500",
  };

  onMount(() => {
    store.setSocket("test");
    store.subscribe((currentMessage) => {
      messages = [...messages, currentMessage];
    });
  });

  function onSendMessage() {
    if (message.length > 0) {
      store.sendMessage(message);
      message = "";
    }
  }
</script>

<main class={style.wrapper}>
  <h1 class={style.title}>Hello visitor!</h1>
  <p class={style.body}>
    Visit the <a class={style.url} href="https://svelte.dev/tutorial"
      >Svelte tutorial</a
    > to learn how to build Svelte apps.
  </p>
  <hr />
  <ol>
    {#each messages as message}
      <li class={style.body}>
        {message}
      </li>
    {/each}
  </ol>

  <input class={style.input} type="text" bind:value={message} />
  <button class={style.button} on:click={onSendMessage}> Send Message </button>
</main>
