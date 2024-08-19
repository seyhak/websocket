import { For } from "solid-js";
import { getAndSetUserCookie, useSaper } from "./useSaper";
import { BOMB_VALUE, CLEAN_VALUE } from "../../constants";

const getValue = (val: number) => {
  switch (val) {
    case BOMB_VALUE:
      return "💣";
    case CLEAN_VALUE:
      return "🍀";
    default:
      break;
  }
};
export function Saper() {
  const {
    field,
    onFieldClick,
    matchState,
    points,
    user,
    onConnectClick,
    setUser,
  } = useSaper();
  return (
    <div>
      <h1>{"SAPER " + getAndSetUserCookie()}</h1>
      <h2>{"POINTS: " + points()}</h2>
      <h2>{"STATUS: " + matchState()}</h2>
      <div class="flex gap-2 justify-center items-center">
        <h2>{"username: "}</h2>
        <input
          type="text"
          value={user()}
          onChange={(e) => {
            setUser(() => e.target.value);
          }}
        />
        <button onClick={onConnectClick}>Connect</button>
      </div>
      <div>
        <table class="bg-slate-300">
          <For each={field}>
            {(row, y) => (
              <tr>
                <For each={row}>
                  {(col, x) => (
                    <td
                      class={`
                  w-8 h-8 border-green-600
                  border ${
                    matchState() === 3 ? "cursor-pointer" : "cursor-crosshair"
                  }
                  `}
                      onClick={() => onFieldClick(x(), y())}
                    >
                      {getValue(col)}
                    </td>
                  )}
                </For>
              </tr>
            )}
          </For>
        </table>
      </div>
    </div>
  );
}
