import { readFile } from "node:fs/promises";

const input = (await readFile("input-b.txt")).toString();

let state;
let a;
let b;
let sum = 0;

function reset() {
  state = "enabled";
  a = undefined;
  b = undefined;
}

function isDigit(c) {
  return c >= "0" && c <= "9";
}

function peek(input, idx, match) {
  return input.substring(idx, idx + match.length) === match;
}

reset();

const MUL = "mul(";
const DO = "do()";
const DONT = "don't()";

for (let i = 0; i < input.length; ) {
  // console.log(i, input[i], state, a, b, sum);
  switch (state) {
    case "enabled":
      if (peek(input, i, MUL)) {
        state = "read-num-a";
        i += MUL.length;
        continue;
      }
      if (peek(input, i, DONT)) {
        state = "disabled";
        i += DONT.length;
        continue;
      }
      // else advance one
      ++i;
      continue;
    case "disabled":
      if (peek(input, i, DO)) {
        state = "enabled";
        i += DO.length;
        continue;
      }
      // else advance one
      ++i;
      continue;
    case "read-num-a":
      // handle mul(,
      if (a === undefined && !isDigit(input[i])) {
        reset();
        continue;
      }
      a = 0;
      while (isDigit(input[i])) {
        a = a * 10 + parseInt(input[i]);
        ++i;
      }
      state = "read-num-sep";
      continue;
    case "read-num-sep":
      if (input[i] === ",") {
        state = "read-num-b";
        ++i;
        continue;
      }
      reset();
      continue;
    case "read-num-b":
      // handle mul(1,)
      if (b === undefined && !isDigit(input[i])) {
        reset();
        continue;
      }
      b = 0;
      while (isDigit(input[i])) {
        b = b * 10 + parseInt(input[i]);
        ++i;
      }
      state = "read-)";
      continue;
    case "read-)":
      if (input[i] === ")") {
        // finalize
        sum += a * b;
        ++i;
      }
      reset();
      continue;
  }
}

console.log("sum", sum);
