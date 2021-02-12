import { writable } from "svelte/store";

type Message = {
  type: string;
  room?: string;
  body: { from: string; msg: string };
};

const messageStore = writable<string>("");

const d = `localhost:1323`;
const p = `votevotevotevote.herokuapp.com`;
const { protocol } = window.location;
const l = protocol === "https:" ? "wss:" : "ws";
const socketUri = `${l}://${
  // @ts-ignore
  import.meta.env.MODE === "development" ? d : p
}/socket`;

let socket: WebSocket;
const setSocket = (room: string = "default") => {
  socket = new WebSocket(`${socketUri}/${room}`);

  socket.onopen = () => {
    socket.onmessage = (event: MessageEvent) => {
      const data: Message = JSON.parse(event.data);
      messageStore.set(event.data);
    };
  };

  socket.onclose = (event: CloseEvent) => {
    setTimeout(() => {
      setSocket();
    }, 3000);
  };

  socket.onerror = (error: any) => {
    socket.close();
  };
};

const sendMessage = (
  from: string = "default",
  msg: any,
  room: string = "test"
): void => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        type: "message",
        room,
        body: {
          from,
          msg,
        },
      })
    );
  }
};

export default {
  setSocket,
  subscribe: messageStore.subscribe,
  sendMessage,
};
