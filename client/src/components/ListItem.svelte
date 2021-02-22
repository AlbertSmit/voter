<script lang="ts">
  import store from "../stores/websocket";
  export let status: any;
  export let user: User;

  type User = {
    uuid: string;
    name: string;
  };

  type Payload = {
    motivation: string;
    from: User;
    for: User;
  };

  let votes: Payload[] = [];
  store.votes((payload: string) => {
    if (!payload) return;
    const { data }: { data: Payload[] } = JSON.parse(payload);
    votes = data;
  });

  // Count the votes.
  $: count = votes.filter((vote) => vote.for.uuid === user.uuid).length;
  $: {
    console.log(`${user.name} vote count -> ${count}`);
  }

  const style = {
    wrapper:
      "p-4 dark:bg-white bg-opacity-10 text-gray-800 dark:text-white rounded-md",
    text: "antialiased text-xs flex justify-between",
  };
</script>

<div data-user={user} on:click class={style.wrapper}>
  <p class={style.text}>
    {user.name}

    {#if status === "VOTING"}
      <span>
        {count}
      </span>
    {/if}
  </p>
</div>
