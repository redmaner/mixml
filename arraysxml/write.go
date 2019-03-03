package arraysxml

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

	io.WriteString(f, fmt.Sprintf("<?xml version='1.0' encoding='UTF-8'?>\n"))
	io.WriteString(f, fmt.Sprintf("<resources>\n"))

	for _, key := range res.Keys {

		entry := res.Entries[key]
		buf := bytes.NewBufferString("")

		io.WriteString(buf, fmt.Sprintf(`    <%s name="%s">`+"\n", entry.form, entry.name))

		for _, item := range entry.items {
			io.WriteString(buf, fmt.Sprintf(`        <item>%s</item>`+"\n", item))
		}
		io.WriteString(buf, fmt.Sprintf(`    </%s>`+"\n", entry.form))

		io.WriteString(f, buf.String())

	}
	io.WriteString(f, fmt.Sprintf("</resources>\n"))
}
