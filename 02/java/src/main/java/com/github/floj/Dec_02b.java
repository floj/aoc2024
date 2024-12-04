package com.github.floj;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Iterator;
import java.util.LinkedList;
import java.util.List;
import java.util.stream.Stream;

public class Dec_02b {

  public static void main(String... args) throws IOException {
    var saveLevels =
        Files.lines(Path.of("src/main/resources/input.txt"))
            .map(String::trim)
            .filter(line -> !line.isBlank())
            .map(line -> lineToLevels(line))
            .map(levels -> checkSafety(levels) || checkSafetyWithTolerance(levels))
            .filter(safe -> safe)
            .count();

    System.out.format("save levels: %d\n", saveLevels);
  }

  private static boolean checkSafetyWithTolerance(List<Integer> levels) {
    for (var combinations = combinationsWithOneDropped(levels); combinations.hasNext(); ) {
      if (checkSafety(combinations.next())) {
        return true;
      }
    }
    return false;
  }

  private static boolean checkSafety(List<Integer> levels) {
    var itr = levels.iterator();

    var prevNum = itr.next();
    int prevInc = 0;

    while (itr.hasNext()) {
      var num = itr.next();
      var inc = prevNum - num;
      // Any two adjacent levels differ by at least one and at most three.
      if (Math.abs(inc) < 1 || Math.abs(inc) > 3) {
        return false;
      }

      // The levels are either all increasing or all decreasing
      if (prevInc != 0 && Integer.signum(prevInc) != Integer.signum(inc)) {
        return false;
      }
      prevNum = num;
      prevInc = inc;
    }

    return true;
  }

  private static Iterator<List<Integer>> combinationsWithOneDropped(List<Integer> levels) {
    return new Iterator<List<Integer>>() {
      int drop = 0;

      @Override
      public boolean hasNext() {
        return drop < levels.size();
      }

      @Override
      public List<Integer> next() {
        var list = new LinkedList<>(levels);
        list.remove(drop);
        drop++;
        return list;
      }
    };
  }

  private static List<Integer> lineToLevels(String line) {
    var fields = line.split("\\s+");
    return Stream.of(fields).map(String::trim).map(field -> Integer.parseInt(field, 10)).toList();
  }
}
