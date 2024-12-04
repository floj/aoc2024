package com.github.floj;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Stream;

public class Dec_02a {

  public static void main(String... args) throws IOException {
    var saveLevels =
        Files.lines(Path.of("src/main/resources/input.txt"))
            .map(String::trim)
            .filter(line -> !line.isBlank())
            .map(line -> lineToLevels(line))
            .map(levels -> checkSafety(levels))
            .filter(safe -> safe)
            .count();

    System.out.format("save levels: %d\n", saveLevels);
  }

  private static boolean checkSafety(List<Integer> levels) {
    var itr = levels.iterator();

    var increments = new ArrayList<Integer>(4);
    var prev = itr.next();

    while (itr.hasNext()) {
      var num = itr.next();
      var inc = prev - num;
      // Any two adjacent levels differ by at least one and at most three.
      if (Math.abs(inc) < 1 || Math.abs(inc) > 3) {
        return false;
      }
      increments.add(inc);
      prev = num;
    }

    // The levels are either all increasing or all decreasing
    return increments.stream().allMatch(v -> v < 0) || increments.stream().allMatch(v -> v > 0);
  }

  private static List<Integer> lineToLevels(String line) {
    var fields = line.split("\\s+");
    return Stream.of(fields).map(String::trim).map(field -> Integer.parseInt(field, 10)).toList();
  }
}
