<script lang="ts">
  const style = {
    wrapper: "p-4 flex flex-col items-start justify-center mx-auto my-auto",
    title: "text-3xl antialiased font-bold tracking-tight",
    body: "antialiased text-xs",
    url: "text-blue-600",
  };
</script>

<main class={style.wrapper}>
  <h1 class={style.title}>Hello visitor!</h1>
  <p class={style.body}>
    Visit the <a class={style.url} href="https://svelte.dev/tutorial"
      >Svelte tutorial</a
    > to learn how to build Svelte apps.
  </p>
</main>
