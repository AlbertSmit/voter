import { writable } from "svelte/store";

export type Status = "WAITING" | "VOTING" | "PRESENTING";
const statusStore = writable<string>("");
const messageStore = writable<string>("");
const userStore = writable<string>("");
const voteStore = writable<string>("");

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
        case "vote":
          voteStore.update(() => event.data);
          break;
        default:
          break;
      }
    };
  };

  socket.onclose = () => {
    setTimeout(() => {
      setSocket();
    }, 3000);
  };

  socket.onerror = () => {
    socket.close();
  };
};

const sendMessage = (from: string = "default", msg: any): void => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        type: "message",
        data: {
          message: msg,
          from,
        },
      })
    );
  }
};

const vote = (user: { uuid: string; name: string; role?: number }): void => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        type: "vote",
        data: {
          for: user,
          motivation: "hey",
        },
      })
    );
  }
};

const updateUser = (name: any = "default"): void => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
        type: "update",
        data: {
          name,
        },
      })
    );
  }
};

const changeRoomStatus = (status: Status = "WAITING") => {
  if (socket.readyState === 1) {
    socket.send(
      JSON.stringify({
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
  votes: voteStore.subscribe,
  vote,
  sendMessage,
  updateUser,
  changeRoomStatus,
};
