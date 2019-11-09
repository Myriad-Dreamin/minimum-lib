package logger

import (
	"errors"
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"testing"
)
type fields struct {
	srcLogger WrappedStdLoggerInterface
	keyvals   []interface{}
}

func getNormalField() fields {
	logger := NewStdLogger().(stdLogger)
	return fields {
		srcLogger: logger.srcLogger,
		keyvals: logger.keyvals,
	}
}

func Test_stdLogger_Fatal(t *testing.T) {
	type args struct {
		msg     string
		keyvals []interface{}
	}
	var tests = []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name   : "test_easy",
			fields : getNormalField(),
			args : args {
				msg: "panic",
				keyvals: []interface{}{"err", errors.New("errrrrr")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stdLogger{
				srcLogger: tt.fields.srcLogger,
				keyvals:   tt.fields.keyvals,
			}
			defer func() {
				if err := recover(); err != nil {
					sugar.PrintStack()
					fmt.Println(err)
				}
			}()
			s.Fatal(tt.args.msg, tt.args.keyvals...)
		})
	}
}