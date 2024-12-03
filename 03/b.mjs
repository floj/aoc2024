import { readFile } from "node:fs/promises";

const input = (await readFile("input-b.txt")).toString();

let state;
let a;
let b;
let enabled = true;
let sum = 0;

function reset() {
  state = "read-instruction";
  a = undefined;
  b = undefined;
}

function isDigit(c) {
  return c >= "0" && c <= "9";
}

reset();

const MUL = "mul(";
const DO = "do()";
const DONT = "don't()";

for (let i = 0; i < input.length; ) {
  //console.log(i, input[i], state, a, b, sum);
  switch (state) {
    case "read-instruction":
      if (input.substring(i, i + MUL.length) === MUL) {
        state = "read-num-a";
        i += 4;
        continue;
      }
      if (input.substring(i, i + DO.length) === DO) {
        i += DO.length;
        enabled = true;
        continue;
      }
      if (input.substring(i, i + DONT.length) === DONT) {
        i += DONT.length;
        enabled = false;
        continue;
      }
      i++;
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
        if (enabled) {
          sum += a * b;
        }
        ++i;
      }
      reset();
  }
}

console.log("sum", sum);
