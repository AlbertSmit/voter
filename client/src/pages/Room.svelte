<script lang="ts">
  import { meta } from "tinro";
  import { onMount } from "svelte";
  import store from "../stores/websocket";
  import { Modal, State, Button, Panel, ListItem } from "../components";
  import type { Status } from "../stores/websocket";
  import iam from "../stores/iam";

  type MessageBody = {
    message: string;
    from: string;
  };

  type MessageResponse = {
    type: "message";
    data: MessageBody;
  };

  type StatusBody = {
    status: Status;
  };

  type StatusResponse = {
    type: "status";
    data: StatusBody;
  };

  let modalIsOpen: boolean = false;

  let message: string;
  let messages: MessageBody[] = [];

  const route = meta();
  let room: string = route.params.id;

  let name: any;
  let submitted: boolean = false;

  const style = {
    wrapper:
      "p-4 h-screen flex flex-col items-start justify-center mx-auto my-auto bg-white dark:bg-gray-900",
    cloak:
      "z-50 inset-0 p-4 antialiased absolute dark:text-white bg-white dark:bg-gray-900",
    title:
      "text-3xl antialiased font-bold tracking-tight text-gray-800 dark:text-white",
    span: "flex w-full mb-4 items-center",
    icon:
      "h-5 w-5 ml-4 text-gray-400 hover:text-gray-500 transition-all cursor-pointer",
    body: "antialiased text-xs",
    url: "text-blue-600",
    container: "inset-0 w-full flex-1",
    messages: "p-4 space-y-2 flex-1 flex-col",
    chat: "p-4 w-full flex flex-col bottom-0 inset-x-0 absolute",
    input:
      "text-gray-800 bg-white dark:text-white dark:bg-gray-700 p-2 my-1 outline-none block rounded-md",
    button:
      "px-6 py-2 text-xs antialiased font-medium rounded-md whitespace-nowrap",
  };

  onMount((): void => {
    store.setSocket(room);
    store.subscribe((payload) => {
      if (!payload) return;
      const { data }: MessageResponse = JSON.parse(payload);
      messages = [...messages, { ...data }];
    });
  });

  let status: Status;
  store.status((payload) => {
    if (!payload) return;
    const { data }: StatusResponse = JSON.parse(payload);
    status = data.status;
  });

  type User = {
    uuid: string;
    name: string;
  };

  let users: User[] = [];
  store.users((payload) => {
    if (!payload) return;
    const { data }: { data: User[] } = JSON.parse(payload);
    users = data;
  });

  function onSendMessage(): void {
    if (message.length > 0) {
      store.sendMessage(name, message);
      message = "";
    }
  }

  function onCastVote(user: {
    uuid: string;
    name: string;
    role?: number;
  }): void {
    store.vote(user.uuid);
  }

  function promptForName(): void {
    const data = prompt("Please enter your name", "Harry Potter");
    if (data != null) {
      name = data;
      void onFinalizeName();
    }
  }

  function onFinalizeName(): void {
    submitted = true;
    void store.updateUser(name);
  }

  async function copyCode(): Promise<void> {
    await navigator.clipboard.writeText(location.href);
  }

  function setRoom(status: Status): void {
    void store.changeRoomStatus(status);
  }
</script>

<main class={style.wrapper}>
  {#if !submitted}
    <div class={style.cloak}>
      <h1 class="text-gray-800 dark:text-white">Hey!</h1>
      <p class="text-gray-800 dark:text-white">What's your name?</p>

      <input class={style.input} type="text" bind:value={name} />
      <Button on:submit={onFinalizeName} on:click={onFinalizeName}>
        Submit
      </Button>
    </div>
  {/if}

  {#if modalIsOpen}
    <Modal close={() => (modalIsOpen = false)}>
      <h1>Hey!</h1>
    </Modal>
  {/if}

  <span class={style.span}>
    <div class="flex items-center flex-1">
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
    </div>
    <div class="flex-1 flex justify-end">
      <button
        class={`text-white dark:bg-gray-800 ${style.button}`}
        on:click={promptForName}>change name</button
      >
    </div>
  </span>

  <State {status} />

  {#if $iam}
    <Panel>
      <p
        class={status !== "WAITING"
          ? "text-gray-500 antialiased"
          : "text-white antialiased"}
        on:click={() => setRoom("WAITING")}
      >
        Wait
      </p>
      <p
        class={status !== "VOTING"
          ? "text-gray-500 antialiased"
          : "text-white antialiased"}
        on:click={() => setRoom("VOTING")}
      >
        Vote
      </p>
      <p
        class={status !== "PRESENTING"
          ? "text-gray-500 antialiased"
          : "text-white antialiased"}
        on:click={() => setRoom("PRESENTING")}
      >
        Present
      </p>
    </Panel>
  {/if}

  <div class={style.container}>
    <ol class={style.messages}>
      {#each messages as content}
        <li class={style.body}>
          <b>{content.from}</b>
          {content.message}
        </li>
      {/each}
    </ol>
    <div class="w-full flex flex-col space-y-1">
      {#each users as user}
        <ListItem {status} {user} on:click={() => onCastVote(user)} />
      {/each}
    </div>
  </div>

  <!-- <div class={style.chat}>
    <input class={style.input} type="text" bind:value={message} />
    <button
      class={`bg-gray-800 text-white dark:text-gray-800 dark:bg-white ${style.button}`}
      on:submit={onSendMessage}
      on:click={onSendMessage}
    >
      Send Message
    </button>
  </div> -->
</main>
