import { writable } from "svelte/store";

const iamStore = writable<string>("");

export default {
  set: iamStore.set,
  subscribe: iamStore.subscribe,
};
