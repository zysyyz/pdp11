package pdp11

import "os"

//import "fmt"

var TKS, TPS, keybuf int

func clearterminal() {
	TKS = 0
	TPS = 1 << 7
}

func writeterminal(c int) {
	os.Stdout.Write([]byte{byte(c)})
}

func addchar(c int) {
	TKS |= 0x80
	keybuf = c
	if TKS&(1<<6) != 0 {
		interrupt(INTTTYIN, 4)
	}
}

func specialchar(c int) {
	switch c {
	case 42:
		keybuf = 4
		break
	case 19:
		keybuf = 034
		break
	case 46:
		keybuf = 127
		break
	default:
		return
	}
	TKS |= 0x80
	if TKS&(1<<6) != 0 {
		interrupt(INTTTYIN, 4)
	}
}

var input = []int{'u', 'n', 'i', 'x', '\n'}

func getchar() int {
	if TKS&0x80 == 0x80 {
		TKS &= 0xff7e
		return keybuf
	}
	return 0
}

func consread16(a int) int {
	//fmt.Printf("consread16: %o\n", a)
	switch a {
	case 0777560:
		if len(input) > 0 {
			addchar(input[0])
			input = input[1:]
		}
		return TKS
	case 0777562:
		return getchar()
	case 0777564:
		return TPS
	case 0777566:
		return 0
	}
	panic("read from invalid address " + ostr(a, 6))
}

func conswrite16(a, v int) {
	//fmt.Printf("conswrite16: %o, %o\n", a,v )
	switch a {
	case 0777560:
		if v&(1<<6) != 0 {
			TKS |= 1 << 6
		} else {
			TKS &= ^(1 << 6)
		}
	case 0777564:
		if v&(1<<6) != 0 {
			TPS |= 1 << 6
		} else {
			TPS &= ^(1 << 6)
		}
	case 0777566:
		v &= 0xFF
		if !(TPS&0x80 == 0x80) {
			break
		}
		switch v {
		case 13:
			break
		default:
			writeterminal(v & 0x7F)
		}
		TPS &= 0xff7f
		if TPS&(1<<6) != 0 {
			TPS |= 0x80
			interrupt(INTTTYOUT, 4)
		} else {
			TPS |= 0x80
		}
	default:
		panic("write to invalid address " + ostr(a, 6))
	}
}
