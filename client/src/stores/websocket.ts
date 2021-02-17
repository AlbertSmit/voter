import { writable } from "svelte/store";

export type Status = "WAITING" | "VOTING" | "PRESENTING";
const statusStore = writable<string>("");
const messageStore = writable<string>("");
const userStore = writable<string>("");

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
  if (socket) return;

  socket = new WebSocket(`${socketUri}/${room}`);
  socket.onopen = () => {
    socket.onmessage = (event: MessageEvent) => {
      const { type } = JSON.parse(event.data);

      switch (type) {
        case "message":
          messageStore.set(event.data);
          break;
        case "status":
          statusStore.set(event.data);
          break;
        case "update":
          userStore.update(() => event.data);
          break;
        default:
          break;
      }
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
        type: "message",
        data: {
          message: msg,
          from,
        },
      })
    );
  }
};

const updateUser = (name: any = "default", room: string = "test"): void => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        room,
        type: "update",
        data: {
          name,
        },
      })
    );
  }
};

const changeRoomStatus = (
  status: Status = "WAITING",
  room: string = "test"
) => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        room,
        type: "status",
        state: {
          status,
        },
      })
    );
  }
};

export default {
  setSocket,
  subscribe: messageStore.subscribe,
  status: statusStore.subscribe,
  users: userStore.subscribe,
  sendMessage,
  updateUser,
  changeRoomStatus,
};
