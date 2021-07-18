package sha1

import (
	"strings"
)

var hexcase = 0 /* 十六进制输出格式。0 -小写；1 -大写 */
var chrsz = 8   /* 每个输入字符的位数。8 - ASCII；16 -统一码 */

func Hex_sha1(s string) string {
	return binb2hex(core_sha1(str2binb(s), int(len(s)*chrsz)))
}

func binb2hex(binarray []int) string {
	var hex_tab []int32
	if hexcase > 0 {
		hex_tab = []int32("0123456789ABCDEF")
	} else {
		hex_tab = []int32("0123456789abcdef")
	}

	str := make([]string, 0, len(binarray)*4)
	for i := 0; i < len(binarray)*4; i++ {
		str = append(str, string(hex_tab[(binarray[int(i>>2)]>>((3-i%4)*8+4))&0xF]), string(hex_tab[int(binarray[int(i>>2)]>>((3-i%4)*8))&0xF]))
	}
	return strings.Trim(strings.Join(str, ""), " ")
}

func str2binb(str string) []int {
	int32Arr := []int32(str)
	bin := make([]int, (len(int32Arr)*chrsz)>>5+1)
	mask := int(1<<chrsz) - 1
	for i := 0; i < len(int32Arr)*chrsz; i += chrsz {
		bin[int(i>>5)] |= int((int32Arr[i/chrsz] & int32(mask)) << (24 - i%32))
	}
	return bin
}

func rol(num, cnt int) int {
	return int(int((num << cnt)) | int(int(uint32(num))>>(32-cnt)))
}

func safe_add(x, y int) int {
	var lsw = (x & 0xFFFF) + (y & 0xFFFF)
	var msw = int((x >> 16)) + int(y>>16) + int(lsw>>16)
	return int((msw << 16)) | int((lsw & 0xFFFF))
}

func sha1_ft(t, b, c, d int) int {
	if t < 20 {
		return (b & c) | ((^b) & d)
	}
	if t < 40 {
		return b ^ c ^ d
	}
	if t < 60 {
		return (b & c) | (b & d) | (c & d)
	}
	return b ^ c ^ d
}

func sha1_kt(t int) int {
	if t < 20 {
		return 1518500249
	}
	if t < 40 {
		return 1859775393
	}
	if t < 60 {
		return -1894007588
	}
	return -899497514
}

func core_sha1(x []int, lenth int) []int {
	x[int(lenth>>5)] |= int(0x80 << (24 - lenth%32))
	newLen := (((lenth + 64) >> 9) << 4) + 15
	if len(x) < int(newLen) {
		xx := make([]int, int(newLen)+1)
		copy(xx[0:newLen], x[:])
		x = xx
	}
	x[newLen] = int(lenth)
	xx := make([]int, int32(len(x)/16+80))
	copy(xx, x)
	w := make([]int, 80)
	var a int = 1732584193
	var b int = -271733879
	var c int = -1732584194
	var d int = 271733878
	var e int = -1009589776
	for i := 0; i < len(x); i += 16 {
		olda := a
		oldb := b
		oldc := c
		oldd := d
		olde := e
		for j := 0; j < 80; j++ {
			if j < 16 {
				w[j] = xx[i+j]
			} else {
				w[j] = rol(int(int(int(w[j-3]^w[j-8])^w[j-14])^w[j-16]), 1)
			}
			t := safe_add(safe_add(rol(int(a), 5), sha1_ft(int(j), b, c, d)), safe_add(safe_add(e, w[j]), sha1_kt(int(j))))
			e = d
			d = c
			c = rol(b, 30)
			b = a
			a = t
		}
		a = safe_add(a, olda)
		b = safe_add(b, oldb)
		c = safe_add(c, oldc)
		d = safe_add(d, oldd)
		e = safe_add(e, olde)
	}
	return []int{a, b, c, d, e}
}
