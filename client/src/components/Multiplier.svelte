<script lang="ts">
  export let idx: number = 1;
  export let times: number = 1;
  $: multiplier = Array.from(new Array(times));

  console.log(idx);
</script>

<div>
  {#each multiplier as _, index}
    <slot />
  {/each}
</div>
