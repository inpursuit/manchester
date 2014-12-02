package manchester

import (
  //"log"
  "container/list"
)

const LONG_FRAME = 112
const SHORT_FRAME = 56

var squares []uint16 = make([]uint16, 256)

func init() {
  squares_precompute()
}

func abs8(x int) uint8 {
  if x >= 127 {
    return uint8(x - 127)
  } else {
    return uint8(127 - x)
  }
}

func squares_precompute() {
  var j uint8
  for i := 0; i<256; i++ {
    j = abs8(i)
    squares[i] = uint16(j*j)
  }
}

func Magnitute(buf []byte) []uint16 {
  m := make([]uint16, len(buf)/2)
  var i, j int = 0, 0
  for ; i<len(buf); i, j = i+2, j+1 {
    m[j] = uint16(squares[buf[i]] + squares[buf[i+1]])
  }
  return m
}

/**
 * frame should be an int array of size 14
 **/
func ReadMessages(buf []uint16) *list.List {
  var all *list.List = list.New()
  numBytes := len(buf)
  var frame_len, data_i, index, i int
  var shift uint8
  //log.Printf("ReadMessages numBytes %d\n", numBytes)
  for i=0; i<numBytes; i++ {
    if buf[i] > 1 {
      continue
    }
    frame_len = LONG_FRAME
    data_i = 0
    var frame []int = make([]int, 14)
//    log.Printf("outer loop, i is %d\n", i)
    for ;i<numBytes && buf[i] <= 1 && data_i<frame_len; i, data_i = i+1, data_i+1 {
  //    log.Printf("\tinner loop, i is %d\n", i)
      if buf[i] != 0 {
        index = data_i / 8
        shift = uint8( (7 - (data_i % 8)) )
        frame[index] |= int(uint8((1<<shift)))
      }
      if data_i == 7 {
        if frame[0] == 0 {
          break
        }
        if (frame[0] & 0x80) != 0 {
          frame_len = LONG_FRAME
        } else {
          frame_len = SHORT_FRAME
        }
      }
    }
 //   log.Printf("\t\tdone inner loop\n")
    if data_i < (frame_len-1) {
      //log.Printf("\tIgnore, data_i is less than frame_len\n")
      //log.Printf("data_i is less than frame length, ignoring %d %d\n", data_i, frame_len-1)
      continue
    }
    //log.Printf("adding frame to all List\n")
    arrSize := (frame_len+7)/8
    all.PushBack(frame[:arrSize])
    //log.Printf("\tFound Message! There are now %d messages\n", all.Len())
  }
  //log.Printf("done outer loop - returning\n")
  return all
}
