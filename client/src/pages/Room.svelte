<script lang="ts">
  import { meta } from "tinro";
  import { onMount } from "svelte";
  import store from "../stores/websocket";

  let message: any;
  let messages: any = [];

  const route = meta();
  let room: string = route.params.room;

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

  function setRoom(room: string): void {
    store.setSocket(room);
  }

  onMount((): void => {
    setRoom(room);
    store.subscribe((currentMessage) => {
      messages = [...messages, currentMessage];
    });
  });

  function onSendMessage(): void {
    if (message.length > 0) {
      store.sendMessage(message, room);
      message = "";
    }
  }

  function onSetRoom(): void {
    store.setSocket(room);
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
  <h2>Room: {room}</h2>
  <hr />
  <ol>
    {#each messages as message}
      <li class={style.body}>
        {message}
      </li>
    {/each}
  </ol>

  <label for="message">message</label>
  <input class={style.input} type="text" bind:value={message} />
  <button
    class={style.button}
    on:submit={onSendMessage}
    on:click={onSendMessage}
  >
    Send Message
  </button>

  <label for="room">room</label>
  <input class={style.input} type="text" bind:value={room} />
  <button class={style.button} on:submit={onSetRoom} on:click={onSetRoom}>
    Change Room
  </button>
</main>
