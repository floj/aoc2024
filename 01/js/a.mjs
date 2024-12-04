import { createReadStream } from "node:fs";
import { createInterface } from "node:readline";

const input = createReadStream("input.txt");
try {
  const lineReader = createInterface({
    input,
    crlfDelay: Infinity,
  });

  const left = [];
  const right = [];
  for await (const line of lineReader) {
    const [l, r] = line
      .split(/\s+/)
      .map((s) => s.trim())
      .map((s) => parseInt(s, 10));
    left.push(l);
    right.push(r);
  }

  left.sort();
  right.sort();

  const distance = right.reduce(
    (distance, _, idx) => (distance += Math.abs(right[idx] - left[idx])),
    0
  );
  console.log("distance", distance);
} finally {
  input.close();
}
