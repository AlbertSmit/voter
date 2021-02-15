<script lang="ts">
  import { meta } from "tinro";
  import { onMount } from "svelte";
  import store from "../stores/websocket";
  import type { Status } from "../stores/websocket";

  type Message = {
    from: string;
    msg: string;
  };

  // type WebSocketResponse = {
  //   type: "message" | "notify" | "pong";
  //   body: Message;
  //   Room: string;
  // };

  let message: string;
  let messages: string[] = [];

  const route = meta();
  let room: string = route.params.id;

  let name: any;
  let submitted: boolean = false;

  const style = {
    wrapper:
      "p-4 h-screen flex flex-col items-start justify-center mx-auto my-auto",
    cloak: "z-50 inset-0 p-4 antialiased absolute bg-gray-800",
    title: "text-3xl antialiased font-bold tracking-tight",
    span: "flex mb-4 items-center",
    icon:
      "h-5 w-5 ml-4 text-gray-400 hover:text-gray-500 transition-all cursor-pointer",
    body: "antialiased text-xs",
    url: "text-blue-600",
    container: "inset-0 w-full flex-1",
    messages: "p-4 space-y-2 flex-1 flex-col bg-gray-50 rounded-xl",
    chat: "p-4 w-full flex flex-col bottom-0 inset-x-0 absolute",
    input:
      "border p-2 my-1 focus:ring-indigo-500 focus:border-indigo-500 block sm:text-sm border-gray-300 rounded-md",
    button:
      "bg-green-100 px-6 py-2 text-xs antialiased font-medium rounded-md text-green-500",
  };

  onMount((): void => {
    store.setSocket(room);
    store.subscribe((currentMessage) => {
      if (!currentMessage) return;
      // const { body }: WebSocketResponse = JSON.parse(currentMessage);
      const body: string = currentMessage;
      messages = [...messages, body];
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

  async function copyCode(): Promise<void> {
    await navigator.clipboard.writeText(location.href);
  }

  function setRoom(status: Status): void {
    void store.changeRoomStatus(status, room);
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

  <div>
    <button class={style.button} on:click={() => setRoom("WAITING")}>
      Waiting
    </button>
    <button class={style.button} on:click={() => setRoom("VOTING")}>
      Voting
    </button>
    <button class={style.button} on:click={() => setRoom("PRESENTING")}>
      Presenting
    </button>
  </div>

  <hr />
  <div class={style.container}>
    <ol class={style.messages}>
      {#each messages as message}
        <li class={style.body}>
          <b>{message}</b>
        </li>
      {/each}
    </ol>
  </div>

  <div class={style.chat}>
    <input class={style.input} type="text" bind:value={message} />
    <button
      class={style.button}
      on:submit={onSendMessage}
      on:click={onSendMessage}
    >
      Send Message
    </button>
  </div>
</main>
