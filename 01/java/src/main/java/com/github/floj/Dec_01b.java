package com.github.floj;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.LinkedList;
import java.util.Optional;
import java.util.TreeMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Dec_01b {

  public static void main(String... args) throws IOException {
    var leftList = new LinkedList<Integer>();
    var rightAppearances = new TreeMap<Integer, AtomicInteger>();

    Files.lines(Path.of("src/main/resources/input.txt"))
        .map(String::trim)
        .filter(line -> !line.isBlank())
        .map(line -> line.trim().split("\\s+"))
        .forEach(
            pair -> {
              var leftNum = Integer.parseInt(pair[0].trim(), 10);
              var rightNum = Integer.parseInt(pair[1].trim(), 10);
              leftList.add(leftNum);
              rightAppearances
                  .computeIfAbsent(rightNum, x -> new AtomicInteger())
                  .incrementAndGet();
            });

    var similarity = 0;
    for (var leftNum : leftList) {
      var appearances =
          Optional.ofNullable(rightAppearances.get(leftNum)).map(AtomicInteger::get).orElse(0);
      similarity += leftNum * appearances;
    }

    System.out.format("similarity: %d\n", similarity);
  }
}
