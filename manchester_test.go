package manchester

import "testing"
//import "fmt"

func TestManchester(t *testing.T) {
  var preamble []uint16 = []uint16{1,0,1,0,0,0,0,1,0,1,0,0,0,0,0,0}
  //equates to data: 11110100
  var v []uint16 = []uint16{1,0,1,0,1,0,1,0,0,1,1,0,0,1,0,1}
  var expected []uint16 = []uint16{1,1,1,1,0,1,0,0}

  withpre := append(preamble, v...)

  //fmt.Printf("%X\n",withpre)    
  Manchester(withpre[:])
  //fmt.Printf("%X\n",withpre)
  preslice := withpre[:16]
  dataslice := withpre[16:24]
  postdata := withpre[24:]
  //fmt.Printf("%X\n",preslice)
  //fmt.Printf("%X\n",dataslice)
  //fmt.Printf("%X\n",postdata)

  for i:=0; i<len(preslice); i++ {
    if preslice[i] != MESSAGEGO {
      t.Errorf("Expected %X saw %X\n", MESSAGEGO, preslice[i])
    }
  }
  for i:=0; i<len(postdata); i++ {
    if postdata[i] != OVERWRITE {
      t.Errorf("Expected %X saw %Xi\n", OVERWRITE, postdata[i])
    }
  }

  for i:=0; i<len(expected); i++ {
    if dataslice[i] != expected[i] {
      t.Errorf("%X does not equal expected %X\n", dataslice, expected)
    }
  }
}
