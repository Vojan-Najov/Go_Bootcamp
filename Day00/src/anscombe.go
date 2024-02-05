package main

import (
  "io"
  "os"
  "fmt"
  "math"
  "slices"
)

const errmsg string =
`Incorrect input:
expected a non-empty sequence of integers,
strictly between -100000 and 100000.`

func readNums() []int {
  var nums []int

  for {
    var num int

    _, err := fmt.Scanln(&num)
    if err == io.EOF {
      break
    }
    if err != nil || math.Abs(float64(num)) > 100000.0 {
      return nil
    }

    nums = append(nums, num)
  }

  return nums
}

func average(nums []int) float64 {
  var avg float64
  var n float64 = float64(len(nums))
  for _, num := range nums {
    avg += float64(num) / n
  }
  return avg
}

func median(nums []int) float64 {
  var med float64
  idx := len(nums) / 2

  if len(nums) % 2 != 0 {
    med = float64(nums[idx])
  } else {
    med = float64(nums[idx - 1] + nums[idx]) / 2.0
  }

  return med
}

func mode(nums []int) int {
  m := make(map[int]int)
  for _, num := range nums {
    m[num] = m[num] + 1
  }

  num := 0
  count := 0
  for k, v := range m {
    if v > count {
      count = v
      num = k
    }  
  }

  return num
}

func standardDeviation(nums []int) float64 {
  avg := average(nums)

  var variance float64
  for _, num := range nums {
    tmp := float64(num) - avg
    variance += tmp * tmp 
  }

  return math.Sqrt(variance / float64(len(nums)))
}

func main() {
  nums := readNums()
  if nums == nil {
    fmt.Fprintln(os.Stderr, errmsg)
    return
  }

  slices.Sort(nums)

  fmt.Printf("Mean: %.2f\n", average(nums))
  fmt.Printf("Median: %.2f\n", median(nums))
  fmt.Printf("Mode: %d\n", mode(nums))
  fmt.Printf("SD: %.2f\n", standardDeviation(nums))
}
