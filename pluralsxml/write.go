package pluralsxml

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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
	defer f.Close()

	io.WriteString(f, fmt.Sprintf("<?xml version='1.0' encoding='UTF-8'?>\n"))

	if res.Comment != "" {
		io.WriteString(f, res.Comment)
	}

	io.WriteString(f, fmt.Sprintf("<resources>\n"))

	for _, key := range res.Keys {

		entry := res.Entries[key]
		buf := bytes.NewBufferString("")

		io.WriteString(buf, fmt.Sprintf(`    <plurals name="%s">`+"\n", entry.name))

		for _, item := range entry.items {
			io.WriteString(buf, fmt.Sprintf(`        <item quantity="%s">%s</item>`+"\n", item[0], item[1]))
		}
		io.WriteString(buf, fmt.Sprintf(`    </plurals>`+"\n"))

		io.WriteString(f, buf.String())

	}
	io.WriteString(f, fmt.Sprintf("</resources>\n"))
}
