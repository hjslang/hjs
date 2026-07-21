export const createContext = () => {
  const stack = [];
  return {
    with:
      (values, fn) =>
      (...args) => {
        stack.push(values);
        try {
          return fn(...args);
        } finally {
          stack.pop();
        }
      },
    use: () => {
      if (stack.length == 0) {
        throw new Error("missing context");
      }
      return stack[stack.length - 1];
    },
  };
};
