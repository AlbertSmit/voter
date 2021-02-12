<script lang="ts">
  import { meta } from "tinro";
  import { onMount } from "svelte";
  import store from "../stores/websocket";

  type Message = {
    from: string;
    msg: string;
  };

  type WebSocketResponse = {
    type: "message" | "notify" | "pong";
    body: Message;
    Room: string;
  };

  let message: string;
  let messages: Message[] = [];

  const route = meta();
  let room: string = route.params.id;

  let name: any;
  let submitted: boolean = false;

  const style = {
    wrapper: "p-4 flex flex-col items-start justify-center mx-auto my-auto",
    cloak: "inset-0 p-4 antialiased absolute bg-gray-800",
    title: "text-3xl antialiased font-bold tracking-tight",
    span: "flex mb-4 items-center",
    icon:
      "h-5 w-5 ml-4 text-gray-400 hover:text-gray-500 transition-all cursor-pointer",
    body: "antialiased text-xs",
    url: "text-blue-600",
    input:
      "border p-2 my-1 focus:ring-indigo-500 focus:border-indigo-500 block sm:text-sm border-gray-300 rounded-md",
    button:
      "bg-green-100 px-6 py-2 text-xs antialiased font-medium rounded-md text-green-500",
  };

  onMount((): void => {
    store.setSocket(room);
    store.subscribe((currentMessage) => {
      if (!currentMessage) return;
      const { body }: WebSocketResponse = JSON.parse(currentMessage);
      messages = [...messages, { ...body }];
    });
  });

  function onSendMessage(): void {
    if (message.length > 0) {
      store.sendMessage(name, message, room);
      message = "";
    }
  }

  function onFinalizeName(): void {
    submitted = true;
  }

  function copyCode(): void {
    window.navigator.clipboard.writeText(window.location.href);
  }
</script>

<main class={style.wrapper}>
  {#if !submitted}
    <div class={style.cloak}>
      <h1 class="text-white">Hey!</h1>
      <p class="text-white">What's your name?</p>

      <input class={style.input} type="text" bind:value={name} />
      <button
        class={style.button}
        on:submit={onFinalizeName}
        on:click={onFinalizeName}
      >
        Submit
      </button>
    </div>
  {/if}

  <span class={style.span}>
    <h1 class={style.title}>v4te</h1>
    <svg
      class={style.icon}
      on:click={copyCode}
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 20 20"
      fill="currentColor"
    >
      <path
        d="M15 8a3 3 0 10-2.977-2.63l-4.94 2.47a3 3 0 100 4.319l4.94 2.47a3 3 0 10.895-1.789l-4.94-2.47a3.027 3.027 0 000-.74l4.94-2.47C13.456 7.68 14.19 8 15 8z"
      />
    </svg>
  </span>
  <hr />
  <ol>
    {#each messages as message}
      <li class={style.body}>
        <b>{message.from}</b>
        {message.msg}
      </li>
    {/each}
  </ol>

  <input class={style.input} type="text" bind:value={message} />
  <button
    class={style.button}
    on:submit={onSendMessage}
    on:click={onSendMessage}
  >
    Send Message
  </button>
</main>
