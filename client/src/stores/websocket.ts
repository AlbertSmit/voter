import { writable } from "svelte/store";

const messageStore = writable("");

let socket: WebSocket;
const setSocket = (room: string) => {
  socket = new WebSocket(
    `wss://votevotevotevote.herokuapp.com/ws?room=${room}`
  );

  // Connection opened
  socket.addEventListener("open", () => {
    console.log("WS open.");
  });

  // Listen for messages
  socket.addEventListener("message", (event: any) => {
    messageStore.set(event.data);
  });
};

const sendMessage = (message: any) => {
  if (socket.readyState <= 1) {
    socket.send(message);
  }
};

export default {
  setSocket,
  subscribe: messageStore.subscribe,
  sendMessage,
};
