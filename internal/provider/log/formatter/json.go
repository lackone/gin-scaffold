package formatter

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lackone/gin-scaffold/internal/contract"
	"time"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["msg"] = msg
	fields["level"] = level
	fields["datetime"] = t.Format(time.RFC3339)
	c, err := json.Marshal(fields)
	if err != nil {
		return bf.Bytes(), errors.New("")
	}
	bf.Write(c)
	return bf.Bytes(), nil
}
