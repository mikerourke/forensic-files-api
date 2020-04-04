package tagasuspect

import "gopkg.in/jdkato/prose.v2"

var tokenTagMap = map[string]string{
	"(":    "left round bracket",
	")":    "right round bracket",
	",":    "comma",
	":":    "colon",
	".":    "period",
	"''":   "closing quotation mark",
	"``":   "opening quotation mark",
	"#":    "number sign",
	"$":    "currency",
	"CC":   "conjunction, coordinating",
	"CD":   "cardinal number",
	"DT":   "determiner",
	"EX":   "existential there",
	"FW":   "foreign word",
	"IN":   "conjunction, subordinating or preposition",
	"JJ":   "adjective",
	"JJR":  "adjective, comparative",
	"JJS":  "adjective, superlative",
	"LS":   "list item marker",
	"MD":   "verb, modal auxiliary",
	"NN":   "noun, singular or mass",
	"NNP":  "noun, proper singular",
	"NNPS": "noun, proper plural",
	"NNS":  "noun, plural",
	"PDT":  "predeterminer",
	"POS":  "possessive ending",
	"PRP":  "pronoun, personal",
	"PRP$": "pronoun, possessive",
	"RB":   "adverb",
	"RBR":  "adverb, comparative",
	"RBS":  "adverb, superlative",
	"RP":   "adverb, particle",
	"SYM":  "symbol",
	"TO":   "infinitival to",
	"UH":   "interjection",
	"VB":   "verb, base form",
	"VBD":  "verb, past tense",
	"VBG":  "verb, gerund or present participle",
	"VBN":  "verb, past participle",
	"VBP":  "verb, non-3rd person singular present",
	"VBZ":  "verb, 3rd person singular present",
	"WDT":  "wh-determiner",
	"WP":   "wh-pronoun, personal",
	"WP$":  "wh-pronoun, possessive",
	"WRB":  "wh-adverb",
}

func getModel() *prose.Model {
	entities := []prose.LabeledEntity{
		{
			Start: 10,
			End:   15,
			Label: "WEAPON",
		},
	}

	train := []prose.EntityContext{
		{
			Text:   "This is a knife",
			Spans:  entities,
			Accept: true,
		},
	}

	return prose.ModelFromData("PRODUCT", prose.UsingEntities(train))
}
