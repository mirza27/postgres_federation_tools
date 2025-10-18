package mapping

// planner bertujuan untuk memetakan topic yang dibutuhkan
// planner memmuat beberapa topic
// dan setiap topic tersebut mempunyai beberapa entity

import (
	"fmt"
	"slices"
)

// Planner memetakan topik->entities yang relevan + info join topik lain.
type Planner struct {
	Entities        []Entity
	TopicToEntities map[string][]Entity
	TopicList       []string
}

func NewPlanner(root *Root) *Planner {
	mp := map[string][]Entity{}
	for _, e := range root.Entities {

		// setiap sumber punya Topic
		for _, s := range e.Sources {
			if s.Topic != "" {
				mp[s.Topic] = append(mp[s.Topic], e)
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
		if s.Topic != "" && !slices.Contains(out, s.Topic) {
			out = append(out, s.Topic)
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
		fmt.Printf("%d. Entity: %s (%s)\n", i+1, e.Entity, e.Kind)
		fmt.Printf("   TargetTable: %s\n", e.TargetTable)
		fmt.Printf("   Sources:\n")
		for _, s := range e.Sources {
			fmt.Printf("     - alias=%s from=%s topic=%s\n", s.Alias, s.From, s.Topic)
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
