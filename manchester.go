package manchester

//import "fmt"

const preamble_len = 16
const allowed_errors = 5
const quality = 10

const MESSAGEGO = 253
const OVERWRITE = 254
const BADSAMPLE = 255

func preamble(buf []uint16, i int) bool {
  var low, high uint16 = 0, 65535
  for i2 := 0; i2<preamble_len; i2++ {
    switch i2 {
      case 0, 2, 7, 9:
        high = buf[i+i2]
        break
      default:
        low = buf[i+i2]
        break
    }
    if high <= low {
      return false
    }
  }
  return true
}

func bool_to_int(a bool) uint16 {
  if a {
    return 1
  } else {
    return 0
  }
}

func single_manchester(a, b, c, d uint16 ) uint16 {
  var bit, bit_p bool
  bit_p = a > b
  bit   = c > d

  if quality == 0 {
    return bool_to_int(bit)
  }

  if quality == 5 {
    if bit && bit_p && b > c {
      return BADSAMPLE
    }
    if !bit && !bit_p && b < c {
      return BADSAMPLE
    }
    return bool_to_int(bit)
  }

  if quality == 10 {
    if bit && bit_p && c > b {
      return 1
    }
    if bit && !bit_p && d < b {
      return 1
    }
    if !bit && bit_p && d > b {
      return 0
    }
    if !bit && !bit_p && c < b {
      return 0
    }
    return BADSAMPLE
  }
  if bit &&  bit_p && c > b && d < a {
    return 1
  }
  if bit && !bit_p && c > a && d < b {
    return 1
  }
  if !bit &&  bit_p && c < a && d > b {
    return 0
  }
  if !bit && !bit_p && c < b && d > a {
    return 0
  }
  return BADSAMPLE;
}

func Manchester(buf []uint16) {
  var a, b uint16 = 0, 0
  var bit uint16
  var i2, errors int
  var maximum_i int = len(buf)-1

  //fmt.Printf("length of array is %d\n",len(buf))
  for i := 0; i<maximum_i; {
    for ;i<(len(buf)-preamble_len); i++ {
      if !preamble(buf, i) {
        //fmt.Printf("  not preamble\n")
        continue;
      }
      //fmt.Printf("  FOUND PREAMBLE!\n")
      a = buf[i]
      b = buf[i+1]
      for i2 := 0; i2<preamble_len; i2++ {
        buf[i+i2] = MESSAGEGO;
      }
      i += preamble_len
      break;
    }
    i2 = i
    errors = 0

    //fmt.Printf("i,i2 %d,%d\n",i,i2)
    for i<maximum_i {
      bit = single_manchester(a, b, buf[i], buf[i+1])
      a = buf[i]
      b = buf[i+1]
      //fmt.Printf("a,b,bit: %b,%b,%b\n",a,b,bit)
      if bit == BADSAMPLE {
        //fmt.Printf("  BADSAMPLE\n")
        errors++
        if errors > allowed_errors {
          buf[i2] = BADSAMPLE
          break
        } else {
          if a > b {
            bit = 1
          } else {
            bit = 0
          }
          a = 0
          b = 65535
        }
      }
      buf[i] = OVERWRITE
      buf[i+1] = OVERWRITE
      buf[i2] = bit
      //fmt.Printf("  putting data at index %d\n",i2)
      i+=2
      i2++
    }
  }
}
