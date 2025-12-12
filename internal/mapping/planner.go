package mapping

// planner bertujuan untuk memetakan topic yang dibutuhkan
// planner memmuat beberapa topic
// dan setiap topic tersebut mempunyai beberapa entity

import (
	"fmt"
	"slices"
	"strings"
)

// Planner memetakan topik->entities yang relevan + info join topik lain.
type Planner struct {
	Entities        []Entity
	TopicToEntities map[string][]Entity
	TopicList       []string
}

func NewPlanner(root *Root) *Planner {

	mp := map[string][]Entity{}

	// save topic_name as key -> array of entities
	// example :
	// db_events_orders : [ orders-to-pesanan, ]
	// db_events_publication : [ publication-to-publication&paper, ]

	for _, e := range root.Entities {

		for _, s := range e.Sources {

			topic := SourceTopic(s)
			if topic != "" {
				mp[topic] = append(mp[topic], e)
			}
		}
	}

	var topicList []string
	for t := range mp {
		topicList = append(topicList, t)
	}

	return &Planner{Entities: root.Entities, TopicToEntities: mp, TopicList: topicList}
}

// ExpectedTopics mengembalikan list topik yg dibutuhkan entity e (untuk join-wait).
func ExpectedTopics(e Entity) []string {
	out := []string{}
	for _, s := range e.Sources {
		if topic := SourceTopic(s); topic != "" && !slices.Contains(out, topic) {
			out = append(out, topic)
		}
	}
	return out
}

// Print menampilkan isi Planner dalam format yang mudah dibaca di console.
func (p *Planner) Print() {
	fmt.Println("==== Planner ====")
	fmt.Println("Total Entities :", len(p.Entities))

	// Tampilkan daftar entity
	fmt.Println("-- Entities --")
	for i, e := range p.Entities {
		fmt.Printf("%d. Entity: %s\n", i+1, e.Entity)
		fmt.Printf("   TargetTable: %s\n", e.TargetTable)
		fmt.Printf("   Sources:\n")
		for _, s := range e.Sources {
			fmt.Printf("     - alias=%s from=%s topic=%s\n", s.Alias, s.From, SourceTopic(s))
		}
	}

	// Tampilkan mapping topic -> entities
	fmt.Println("-- TopicToEntities --")
	for topic, ents := range p.TopicToEntities {
		fmt.Printf("Topic: %s\n", topic)
		for _, e := range ents {
			fmt.Printf("   - %s\n", e.Entity)
		}
	}

	fmt.Println("-- TopicToEntities --")
	fmt.Println(p.TopicList)

}

// get key 'from' from each item in sources[]
func SourceTopic(s EntitySource) string {

	name := strings.TrimSpace(s.From)
	if name == "" {
		return ""
	}

	// get base name
	parts := strings.Split(name, ".")
	base := parts[len(parts)-1]
	base = strings.ToLower(strings.TrimSpace(base))
	if base == "" {
		return ""
	}

	// add prefix db_events_
	return fmt.Sprintf("db_events_%s", base)
}
