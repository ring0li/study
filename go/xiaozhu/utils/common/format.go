package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/syyongx/php2go"
	"strings"
)

func FormatName(name string) string {
	if name == "" {
		return ""
	}

	nameRune := []rune(name)
	return string(nameRune[0:1]) + php2go.StrRepeat("*", len(nameRune)-1)
}

//返回值是string不是float，因为float前端展示可能丢失.00
func FormatMoney(money float64) string {
	return php2go.NumberFormat(money/100, 2, ".", "")
}

//unicode转换为字符串string
func Unicode2String(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}
