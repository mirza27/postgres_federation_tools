package main

import (
	"db_migrate_server/internal/config"
	"fmt"
	"math/rand"

	"github.com/0x6flab/namegenerator"
)

func main() {
	generator := namegenerator.NewGenerator()

	researcher_number := 40
	max_paper_authors := 5
	publication_number := 100000

	researcher_list := []Researcher{}
	publication_list := []Publication{}
	journal_list := []Journal{}
	conference_list := []Conference{}
	book_list := []Book{}
	publication_writer_list := []Publication_writer{}

	// add researchers
	for i := 1; i <= researcher_number; i++ {

		writer_name := generator.Generate()

		researcher := Researcher{
			rese_id:            i,
			rese_nif:           randomString(9, false, true, false),
			rese_name:          writer_name,
			rese_first_surname: writer_name[len(writer_name)/2:],
			rese_sec_surname:   writer_name[:len(writer_name)/2],
			rese_birth_day:     randomDate(),
			rese_mail:          writer_name + "@mail.com",
			rese_web:           "www." + randomString(5, false, true, false) + ".com",
			rese_phone:         fmt.Sprintf("%d", randomInt(1000000000)),
			rese_univ:          "University " + randomString(5, false, false, false),
			rese_sign:          randomString(15, false, true, false),
			rese_short_curr_es: "Short CV ES " + randomString(20, false, false, false),
			rese_short_curr_en: "Short CV EN " + randomString(20, false, false, false),
		}

		researcher_list = append(researcher_list, researcher)
	}

	for i := 1; i <= publication_number; i++ {

		publication_type := rand.Intn(3) // 0: journal, 1: conference, 2: book
		writer_number := rand.Intn(max_paper_authors) + 1

		// list publication
		publication := Publication{
			pub_id:     i,
			pub_isbn:   randomString(13, false, false, true),
			pub_title:  randomString(20, true, false, false),
			pub_year:   randomInt(25) + 2000,
			pub_loc:    "Location " + randomString(10, false, false, false),
			pub_edit:   "Editorial " + randomString(10, false, false, false),
			pub_draft:  "",
			pub_prsnt:  "",
			pub_minute: "",
			pub_source: "Source " + randomString(10, false, false, false),
			pub_book:   "",
			pub_link:   "www." + randomString(5, false, true, false) + ".com",
			pub_repo:   "www.repo.com/" + randomString(10, false, true, false),
			pub_visib:  randomBool(),
			pub_cntry:  getRandomCountry().cntry_id,
		}
		publication_list = append(publication_list, publication)

		// list publication writers
		for j := 1; j <= writer_number; j++ {
			pub_writer := Publication_writer{
				rese_id:       getRandomResearcher(researcher_list).rese_id,
				pub_id:        i,
				pwrit_is_edit: randomBool(),
				pwrit_ord:     j,
			}
			publication_writer_list = append(publication_writer_list, pub_writer)
		}

		// list publication by type
		if publication_type == 0 { // journal
			journal := Journal{
				jrnl_id:          i,
				jrnl_journ_title: randomString(15, true, false, false),
				jrnl_vol:         randomString(20, false, false, true),
				jrnl_numb:        randomString(15, false, false, true),
				jrnl_start_page:  randomInt(120),
				jrnl_end_page:    randomInt(200),
				jrnl_month:       randomInt(12),
				jrnl_scope:       randomInt(5),
				jrnl_indexed:     randomBool(),
				jrnl_pos_numb:    randomInt(100),
				jrnl_tot_pos:     randomInt(200),
			}
			journal_list = append(journal_list, journal)
		}

		if publication_type == 1 { // conference

			conference := Conference{
				conf_id:         i,
				conf_name:       "Conference " + randomString(15, false, false, false),
				conf_ref:        "Ref " + randomString(10, false, false, false),
				conf_publ:       "Publisher " + randomString(10, false, false, false),
				conf_edit:       "Editor " + randomString(10, false, false, false),
				conf_start_date: randomDate(),
				conf_end_date:   randomDate(),
				conf_start_page: randomInt(120),
				conf_end_page:   randomInt(200),
				conf_type:       randomInt(3),
				conf_serie:      "Serie " + randomString(10, false, false, false),
				conf_organi:     "Organizer " + randomString(10, false, false, false),
				conf_particip:   randomInt(500),
			}
			conference_list = append(conference_list, conference)
		}

		if publication_type == 2 { // book
			book := Book{
				book_id:      i,
				book_isbn:    randomString(13, false, false, true),
				book_edit:    "Editorial " + randomString(10, false, false, false),
				book_vol:     randomString(10, false, false, false),
				book_editors: "Editors " + randomString(20, false, false, false),
			}

			book_list = append(book_list, book)
		}
	}

	// run chunk insert
	chunk_size := 100

	config := config.Load()
	seederConn, err := NewSeederConn(*config)
	if err != nil {
		panic(err)
	}
	defer seederConn.Pool.Close()

	seederConn.chunkInsertCountry(CountryList, chunk_size)
	seederConn.chunkInsertResearcher(researcher_list, chunk_size)
	seederConn.chunkInsertPublication(publication_list, chunk_size)
	seederConn.chunkInsertJournal(journal_list, chunk_size)
	seederConn.chunkInsertConference(conference_list, chunk_size)
	seederConn.chunkInsertBook(book_list, chunk_size)
	seederConn.chunkInsertPublicationWriter(publication_writer_list, chunk_size)

}
