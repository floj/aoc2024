import { readFile } from "node:fs/promises";

const safe = (await readFile("input.txt"))
  .toString()
  .split("\n")
  .map((line) => line.split(/\s+/).map((field) => parseInt(field.trim(), 10)))
  .filter((levels) => {
    for (const c of relevantCombinations(levels)) {
      if (isSafe(c)) {
        return true;
      }
    }
  }).length;
console.log("safe levels", safe);

/**
 * @param {Array<Number>} levels
 */
function isSafe(levels) {
  let lastSig = 0;
  for (let idx = 1; idx < levels.length; idx++) {
    const distance = levels[idx - 1] - levels[idx];
    // Any two adjacent levels differ by at least one and at most three.
    if (Math.abs(distance) < 1 || Math.abs(distance) > 3) {
      return false;
    }
    // The levels are either all increasing or all decreasing.
    const sig = Math.sign(distance);
    if (lastSig != 0 && Math.sign(distance) != lastSig) {
      return false;
    }
    lastSig = sig;
  }
  return true;
}

/**
 * @param {Array<Number>} levels
 */
function* relevantCombinations(levels) {
  yield [...levels];
  for (let i = 0; i < levels.length; i++) {
    const l = [...levels];
    l.splice(i, 1);
    yield l;
  }
}
