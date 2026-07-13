export const createContext = () => {
  const ctx = {
    currentEffect: {
      id: "",
      fn: undefined,
      effectsId: "",
    },
    cleaners: new Map(),
  };

  return {
    signal: (val) => createSignal(ctx, val),
    effects: () => createEffects(ctx),
  };
};

const createSignal = (ctx, val) => {
  const signal = {
    id: genId(),
    val,
    effects: new Map(),
  };

  function cleaner(effectsId) {
    for (const [id, effect] of signal.effects) {
      if (effectsId === effect.effectsId) {
        signal.effects.delete(id);
      }
    }
  }

  return {
    get: () => {
      const { id, fn, effectsId } = ctx.currentEffect;
      if (id) {
        const cleaners = ctx.cleaners;
        signal.effects.set(id, { fn, effectsId, cleaners });
        cleaners.set(signal.id, cleaner);
      }
      return signal.val;
    },
    set: (val) => {
      signal.val = val;
      const prevEffect = { ...ctx.currentEffect };
      const prevCleaners = ctx.cleaners;
      try {
        for (const [id, effect] of signal.effects) {
          const { fn, effectsId, cleaners } = effect;
          ctx.currentEffect = { id, fn, effectsId };
          ctx.cleaners = cleaners ?? prevCleaners;
          fn?.();
        }
      } finally {
        ctx.currentEffect = prevEffect;
        ctx.cleaners = prevCleaners;
      }
    },
  };
};

const createEffects = (ctx) => {
  const effects = {
    id: genId(),
    signalCleaners: new Map(),
  };

  return {
    effect: (fn) => {
      const prevEffect = { ...ctx.currentEffect };
      const prevCleaners = ctx.cleaners;
      try {
        const id = genId();
        ctx.currentEffect = { id, fn, effectsId: effects.id };
        ctx.cleaners = effects.signalCleaners;
        fn();
      } finally {
        ctx.currentEffect = prevEffect;
        ctx.cleaners = prevCleaners;
      }
    },
    clean: () => {
      for (const [key, cleanerFn] of effects.signalCleaners) {
        cleanerFn(effects.id);
      }
    },
  };
};

let _id = 0;
const genId = () => (++_id).toString();
