<script lang="ts">
  import { Route } from "tinro";
  import Home from "./pages/Home.svelte";
  import Room from "./pages/Room.svelte";
</script>

<Route path="/"><Home /></Route>
<Route path="/room/:id"><Room /></Route>
<Route fallback redirect="/" />
