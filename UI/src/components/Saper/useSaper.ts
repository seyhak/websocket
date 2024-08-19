import { createSignal, onCleanup, onMount } from "solid-js";
import { createStore } from "solid-js/store";
import { GAME_IN_PLAY, END_OF_PLAYERS_FLAG } from "../../constants";

const initWebSocket = async () => {
  initUserCookie()
  const url = "ws://localhost:8080/init";
  const webSocket = new WebSocket(url);
  webSocket.binaryType = "arraybuffer";
  console.log({ webSocket });
  // const result = await fetch("localhost:8080/init")
  // console.log(result)
  // await new Promise(async (res, _) => {
  //   while(webSocket.OPEN !== webSocket.readyState ){
  //     console.log("loop")
  //     await new Promise((s, _) => {
  //       setTimeout(() => {
  //         s(true)
  //       }, 500)
  //     })
  //   }
  //   res(true)
  // })
  return webSocket;
};

const getUserCookie = () => {
  return document.cookie.split("user=")?.[1]?.split(";")?.[0]
}
export const getAndSetUserCookie = () => {
  let user = getUserCookie()
  if (!user) {
    user = Date.now().toString()
    document.cookie = `user=${user};`
  }
  return user
}
function initUserCookie() {
  console.log(document.cookie)
  if(document.cookie){

  } else {
    getAndSetUserCookie()
  }
}

export const useSaper = () => {
  let socket: null | WebSocket = null;
  const [matchState, setMatchState] = createSignal<number>(GAME_IN_PLAY)
  const [points, setPoints] = createSignal<number>(0)
  const [hue, setHue] = createSignal<number>(0)
  const [user, setUser] = createSignal<string>(getUserCookie())
  const [field, setField] = createStore<number[][]>([]);

  function handleOnMessage(this: WebSocket, event: MessageEvent<any>): any {
    console.log({ th: this, event });
    // const blob = new Blob(event.data)
    const view = new DataView(event.data);
    // const data = await blob.arrayBuffer()
    console.log({ view });
    // const data = view.
    // view.getUint8 // TODO
    // const data = new Uint8Array(view.buffer)
    const data = new Array(...new Uint8Array(view.buffer));
    // setHeight(() => data?.[0] || 0)
    // setWidth(() => data?.[0] || 0)
    const width = data.shift();
    const matchStateFlag = data.shift()
    const points = data.shift();
    const hue = data.shift();
    const otherPlayersData = []
    console.log({data})
    while(true) {
      const next = data.shift()
      if (next === END_OF_PLAYERS_FLAG ){
        break
      }
      const playerPoints = data.shift()
      const playerHue = data.shift()
      otherPlayersData.push({
        points: playerPoints,
        hue: playerHue
      })
    }
    const field: number[][] = [];
    let tmp = [];
    for (let index = 0; index < data.length; index++) {
      const element = data[index];
      tmp.push(element);
      if (tmp.length === width) {
        field.push([...tmp]);
        tmp = [];
      }
    }
    setField(() => field);
    setMatchState(() => matchStateFlag)
    setPoints(() => points)
    setHue(() => hue)
  }

  const onFieldClick = (x: number, y: number) => {
    if(matchState() === GAME_IN_PLAY) {
      console.log({ x, y });
      const valBuffer = new Uint8Array([x, y]).buffer;
      // console.log({ valBuffer });
      socket!.send(valBuffer);
    } else {
      const valBuffer = new Uint8Array([1]).buffer;
      socket!.send(valBuffer);
    }
  };

  const onConnectClick = async () => {
    socket = await initWebSocket()
    socket.onopen = () => {
      console.log("onopen");
      const val = new ArrayBuffer(0);
      console.log({ val });
      socket!.send(val);
    };
    document.cookie = `user=${user()};`
    console.log({ socket });
    socket.onmessage = handleOnMessage;
  }
  onMount(async () => {
  });
  onCleanup(() => {
    console.log("onCleanup");
    socket?.close();
  });
  return {
    field,
    onFieldClick,
    matchState,
    points,
    user,
    setUser,
    onConnectClick
    // height,
    // width
  };
};
