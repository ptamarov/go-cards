package models

// TemplateDate holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	CardData  map[string]string
	Progress  int
	CSRFToken string
	Answer    string
	Count     int
}
