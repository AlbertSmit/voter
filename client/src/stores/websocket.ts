import { writable } from "svelte/store";

type Payload = {
  payloadType: string;
  message: string;
  from: string;
};

// Message gets send around.
type Message = {
  data: Payload;
  room: string;
};

const messageStore = writable<string>("");

const d = `localhost:1323`;
const p = `votevotevotevote.herokuapp.com`;
const { protocol } = window.location;
const l = protocol === "https:" ? "wss:" : "ws";
const socketUri = `${l}://${
  // @ts-ignore
  import.meta.env.MODE === "development" ? d : p
}/api/socket`;

let socket: WebSocket;
const setSocket = (room: string = "default") => {
  socket = new WebSocket(`${socketUri}/${room}`);

  socket.onopen = () => {
    socket.onmessage = (event: MessageEvent) => {
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
        room,
        data: {
          type: "message",
          message: msg,
          from,
        },
      })
    );
  }
};

export type Status = "WAITING" | "VOTING" | "PRESENTING";
const changeRoomStatus = (
  status: Status = "WAITING",
  room: string = "test"
) => {
  if (socket.readyState === 1) {
    socket.send(status);
  }
};

export default {
  setSocket,
  subscribe: messageStore.subscribe,
  sendMessage,
  changeRoomStatus,
};
