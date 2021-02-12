import { writable } from "svelte/store";

type Message = {
  type: string;
  room?: string;
  body: { from: string; msg: string };
};

const messageStore = writable("");

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

  socket.onopen = function () {
    console.log("Connected");

    socket.onmessage = (event: MessageEvent) => {
      const data: Message = JSON.parse(event.data);
      console.log("Message!", data.body);
      messageStore.set(data.body.msg);
    };
  };

  socket.onclose = (event: CloseEvent) => {
    console.log(
      "Socket is closed. Reconnect will be attempted in 3 second.",
      event.reason
    );
    setTimeout(function () {
      setSocket();
    }, 3000);
  };

  socket.onerror = (error: any) => {
    console.error(
      "Socket encountered error: ",
      error.message,
      "Closing socket"
    );
    socket.close();
  };
};

const sendMessage = (message: any, room: string = "test"): void => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        type: "message",
        room,
        body: {
          from: "Albert",
          msg: message,
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
