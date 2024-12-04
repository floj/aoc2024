import { readFile } from "node:fs/promises";

const input = (await readFile("input.txt")).toString();

let state;
let a;
let b;
let sum = 0;

function reset() {
  state = "find-m";
  a = undefined;
  b = undefined;
}

function isDigit(c) {
  return c >= "0" && c <= "9";
}

reset();

for (let i = 0; i < input.length; ) {
  let c = input[i];
  //console.log(i, c, state, a, b, sum);
  switch (state) {
    case "find-m":
      if (c === "m") {
        state = "read-u";
      }
      ++i;
      continue;
    case "read-u":
      if (c === "u") {
        state = "read-l";
        ++i;
        continue;
      }
      reset();
      continue;
    case "read-l":
      if (c === "l") {
        state = "read-(";
        ++i;
        continue;
      }
      reset();
      continue;
    case "read-(":
      if (c === "(") {
        state = "read-num-a";
        i++;
        continue;
      }
      reset();
      continue;
    case "read-num-a":
      // handle mul(,
      if (!isDigit(c) && a === undefined) {
        reset();
        continue;
      }
      if (!isDigit(c)) {
        state = "read-num-sep";
        continue;
      }

      if (a === undefined) {
        a = 0;
      }
      a = a * 10 + parseInt(c, 10);
      ++i;
      continue;
    case "read-num-sep":
      if (c === ",") {
        state = "read-num-b";
        ++i;
        continue;
      }
      reset();
      continue;
    case "read-num-b":
      // handle mul(1,)
      if (!isDigit(c) && b === undefined) {
        reset();
        continue;
      }
      if (!isDigit(c)) {
        state = "read-)";
        continue;
      }

      if (b === undefined) {
        b = 0;
      }
      b = b * 10 + parseInt(c, 10);
      ++i;
      continue;
    case "read-)":
      if (c === ")") {
        ++i;
        state = "finalize";
        continue;
      }
      reset();
      continue;
    case "finalize":
      sum += a * b;
      reset();
  }
}

console.log("sum", sum);
