package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/inspect"
)

var MainInsp = inspect.NewInspector("webFingerDetect", 9999)

var (
	Level, _ = MainInsp.NewType("level", func(i interface{}) string {
		text := i.(string)
		return "  " + text + "  |"
	})
	Text, _ = MainInsp.NewType("text", func(i interface{}) string {
		if str, ok := i.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", i)
	})
	FileName, _ = MainInsp.NewType("file_name", func(i interface{}) string {
		return "  " + i.(string) + "  :  "
	}, decorators.Magenta)
	FatalError, _ = MainInsp.NewType("fatal_error", func(i interface{}) string {
		return i.(error).Error()
	}, decorators.Red)
	WarnError, _ = MainInsp.NewType("warn_error", func(i interface{}) string {
		return i.(error).Error()
	}, decorators.Yellow)
	Debug, _ = MainInsp.NewType("debug", func(i interface{}) string {
		return i.(string)
	})
	Json, _ = MainInsp.NewType("json", func(i interface{}) string {
		var out bytes.Buffer
		if obj, ok := i.(string); ok {
			err := json.Indent(&out, []byte(obj), "", "  ")
			if err != nil {
				return err.Error()
			}
			return "\n" + out.String() + "\n"
		} else {
			marshal, err := json.MarshalIndent(i, "", "  ")
			if err != nil {
				return err.Error()
			}
			return "\n" + string(marshal) + "\n"
		}

	}, decorators.Cyan)
)

var (
	LEVEL_INFO    = Level("INFO", decorators.Blue)
	LEVEL_WARNING = Level("WARN", decorators.Yellow)
	LEVEL_ERROR   = Level("ERRO", decorators.Red)
	LEVEL_DEBUG   = Level("DEBUG", decorators.Cyan)
)

func initDecoration() {
	MainInsp.SetTypeDecorations("_start", decorators.Yellow)
	MainInsp.SetTypeDecorations("_func", decorators.Magenta, decorators.Invisible)
	MainInsp.SetOrders("_time", Level, "_func", Text, FatalError)
	MainInsp.SetSeparator("")
}
