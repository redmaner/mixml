package stringsxml

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Write writes android strings.xml
func (res *Resources) Write() {

	if len(res.Keys) == 0 {
		return
	}

	if _, err := os.Stat(res.FilePath); err == nil {
		err := os.Remove(res.FilePath)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}

	f, err := os.Create(res.FilePath)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	io.WriteString(f, fmt.Sprintf("<?xml version='1.0' encoding='UTF-8'?>\n"))
	io.WriteString(f, fmt.Sprintf("<resources>\n"))

	for _, key := range res.Keys {

		var formatString string
		var writeString string
		stringItem := res.Items[key]

		if stringItem.formatted {
			formatString = ` formatted="false"`
		}

		value := stringItem.value

		// We fix apostrophe errors, by adding a \ in front of it
		// this slash is only added when it does not yet exist
		if stringItem.apostropheFix && value[0] != '"' {
			var newValue string

			strSlice := strings.Split(value, "'")
			splits := len(strSlice)
			for i, v := range strSlice {

				if v == "" {
					continue
				}

				if i == splits-1 {
					newValue = newValue + v
					break
				}

				lastChar := len(v) - 1
				if v[lastChar] == 92 {
					newValue = newValue + v + "'"
					continue
				}
				newValue = newValue + v + `\'`
			}
			value = newValue
		}

		writeString = fmt.Sprintf(`  <string name="%s"%s>%s</string>`+"\n", stringItem.name, formatString, value)

		f.WriteString(writeString)

	}
	io.WriteString(f, fmt.Sprintf("</resources>\n"))

}
