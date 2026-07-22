import udomdiff from "udomdiff";
import { createContext } from "./index.js";

const ctx = createContext();
let currentEffects = ctx.effects();

export const signal = ctx.signal;
export const effect = (...args) => currentEffects.effect(...args);

export function Render(fn) {
  const node = document.createTextNode("");
  effect((_) => {
    node.textContent = fn();
  });
  return node;
}

const effectsMap = new WeakMap();
export function domdiff(parentNode, rows, insert) {
  if (!effectsMap.has(parentNode)) {
    effectsMap.set(parentNode, Object.create(null));
  }
  const savedEffects = effectsMap.get(parentNode);
  const prevEffects = currentEffects;
  try {
    const ids = {};
    const rowIdSet = new Set(rows.map((r) => `${r.id}`));
    for (const child of parentNode.children) {
      const id = child.dataset.id;
      if (id === undefined) {
        throw new Error("All children must have dataset.id");
      }
      if (!rowIdSet.has(id) && savedEffects[id]) {
        // cleaning effects group
        savedEffects[id].clean();
        delete savedEffects[id];
      }
      ids[id] = child;
    }
    let nodes = [];
    for (const row of rows) {
      if (row.id === undefined || row.id === null) {
        throw new Error("All rows must have an id property");
      }
      if (!savedEffects[row.id]) {
        // creates an effects group per row
        savedEffects[row.id] = ctx.effects();
      }
      currentEffects = savedEffects[row.id];
      nodes.push(ids[`${row.id}`] || insert(row));
    }
    udomdiff(parentNode, [...parentNode.children], nodes, (n) => n);
  } finally {
    currentEffects = prevEffects;
  }
}
