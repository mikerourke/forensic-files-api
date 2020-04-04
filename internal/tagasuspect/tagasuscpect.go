package tagasuspect

import (
	"fmt"
	"log"

	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"gopkg.in/jdkato/prose.v2"
)

func Test() {
	ep, err := whodunit.NewEpisodeFromName("08-11-a-wrong-foot")
	if err != nil {
		return
	}
	t := killigraphy.NewTranscript(ep)
	// Create a new document with the default configuration:
	doc, err := prose.NewDocument(t.Read())
	if err != nil {
		log.Fatal(err)
	}

	for _, tok := range doc.Tokens() {
		fmt.Println(tok.Text, tokenTagMap[tok.Tag])
	}

	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
	}
}
