package com.github.floj;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Collections;
import java.util.LinkedList;

public class Dec_01a {

  public static void main(String... args) throws IOException {
    var leftList = new LinkedList<Integer>();
    var rightList = new LinkedList<Integer>();

    Files.lines(Path.of("src/main/resources/dec_01/input.txt"))
        .map(String::trim)
        .filter(line -> !line.isBlank())
        .map(line -> line.trim().split("\\s+"))
        .forEach(
            pair -> {
              leftList.add(Integer.parseInt(pair[0].trim(), 10));
              rightList.add(Integer.parseInt(pair[1].trim(), 10));
            });

    Collections.sort(leftList);
    Collections.sort(rightList);

    var totalDistance = 0;
    while (!leftList.isEmpty() && !rightList.isEmpty()) {
      var smallestLeft = leftList.pollFirst();
      var smallestRight = rightList.pollFirst();
      totalDistance += Math.abs(smallestLeft - smallestRight);
    }
    System.out.format("total distance: %d\n", totalDistance);
  }
}
