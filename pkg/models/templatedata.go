package models

import "github.com/ptamarov/go-cards/app/card"

// TemplateDate holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	BoolMap   map[string]bool
	CardData  card.MemoryCard
}
