import { writable } from "svelte/store";
import { io, Socket } from "socket.io-client";

const messageStore = writable("");
const socketUri = `wss://votevotevotevote.herokuapp.com/socket`;

let socket: Socket;
let socket2: Socket;

const setSocket = (room: string) => {
  socket = io(`${socketUri}/${room}`, { transports: ["websocket"] });
  socket2 = io(`${socketUri}/${room}/chat`, { transports: ["websocket"] });

  socket.on("connect", () => {
    console.log(socket.id);
  });

  socket.on("disconnect", () => {
    console.log(socket.id); // undefined
  });

  socket.on("reply", (msg: string) => {
    console.log("Message!", msg);
    messageStore.set(msg);
  });
};

const sendMessage = (message: any) => {
  if (socket.connected) {
    socket2.emit("msg", message);
    socket.emit("notice", message);
  }
};

export default {
  setSocket,
  subscribe: messageStore.subscribe,
  sendMessage,
};
