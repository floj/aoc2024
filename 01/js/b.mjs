import { createReadStream } from "node:fs";
import { createInterface } from "node:readline";

const input = createReadStream("input.txt");
try {
  const lineReader = createInterface({
    input,
    crlfDelay: Infinity,
  });

  const left = [];
  const right = new Map();
  for await (const line of lineReader) {
    const [l, r] = line
      .split(/\s+/)
      .map((s) => s.trim())
      .map((s) => parseInt(s, 10));
    left.push(l);
    right.set(r, (right.get(r) || 0) + 1);
  }

  const similarity = left.reduce(
    (similarity, v) => similarity + v * (right.get(v) || 0),
    0
  );
  console.log("similarity", similarity);
} finally {
  input.close();
}
