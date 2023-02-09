package formatter

import (
	"bytes"
	"fmt"
	"github.com/lackone/gin-scaffold/internal/contract"
	"time"
)

func TextFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	Separator := "\t"

	prefix := Prefix(level)

	bf.WriteString(prefix)
	bf.WriteString(Separator)

	ts := t.Format(time.RFC3339)
	bf.WriteString(ts)
	bf.WriteString(Separator)

	bf.WriteString("\"")
	bf.WriteString(msg)
	bf.WriteString("\"")
	bf.WriteString(Separator)

	bf.WriteString(fmt.Sprint(fields))
	return bf.Bytes(), nil
}
