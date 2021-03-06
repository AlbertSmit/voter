<script lang="ts">
  import { Route } from "tinro";
  import Home from "./pages/Home.svelte";
  import Room from "./pages/Room.svelte";
  import Test from "./pages/Test.svelte";
</script>

<Route path="/"><Home /></Route>
<Route path="/room/:id"><Room /></Route>
<Route path="/test"><Test /></Route>
<Route fallback redirect="/" />
