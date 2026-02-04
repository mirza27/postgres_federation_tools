package main

import (
	"context"
	"db_migrate_server/internal/config"
	"fmt"
	"os"
	"time"

	"github.com/0x6flab/namegenerator"
)

func runStreamSeed() {

	config := config.Load()
	seederConn, err := NewSeederConn(*config)
	if err != nil {
		panic(err)
	}
	defer seederConn.Pool.Close()

	seederConn.chunkInsertCountry(CountryList, 100)
	// comment unused func
	// insertScenario1(*seederConn)
	insertScenario2(*seederConn)
	// insertScenario3(*seederConn)
	// insertScenario4(*seederConn)
	// insertScenario5(*seederConn)

}

// insert researchers
func insertScenario1(db SeederConn) {

	start_id := 1
	end_id := 100
	count_each_second := 10

	generator := namegenerator.NewGenerator()

	researcher_list := []Researcher{}

	for i := start_id; i <= end_id; i++ {

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

	if count_each_second <= 0 {
		count_each_second = 1
	}

	logFile, err := os.OpenFile("scenario1.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	ctx := context.Background()
	for i, r := range researcher_list {
		ts := time.Now().Format(time.RFC3339Nano)
		fmt.Fprintf(logFile, "%s inserting researcher id %d\n", ts, r.rese_id)

		if err := db.insertSingleResearcher(ctx, r); err != nil {
			fmt.Fprintf(logFile, "%s error inserting researcher id %d: %v\n", ts, r.rese_id, err)
			panic(err)
		}

		// throttle per second based on count_each_second
		if (i+1)%count_each_second == 0 && i+1 < len(researcher_list) {
			time.Sleep(time.Second)
		}
	}

}

func insertScenario2(db SeederConn) {

	start_id := 1
	end_id := 100
	count_each_second := 10
	generator := namegenerator.NewGenerator()

	publication_writer_list := []Publication_writer{}
	researcher_list := []Researcher{}

	for i := start_id; i <= end_id; i++ {

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

	for i := start_id; i <= end_id; i++ {

		pub_writer := Publication_writer{
			rese_id:       getRandomResearcher(researcher_list).rese_id,
			pub_id:        i,
			pwrit_is_edit: randomBool(),
			pwrit_ord:     i,
		}
		publication_writer_list = append(publication_writer_list, pub_writer)
	}

	// insert researchers first to avoid foreign key constraint error
	db.chunkInsertResearcher(researcher_list, 100)

	if count_each_second <= 0 {
		count_each_second = 1
	}

	logFile, err := os.OpenFile("scenario2.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	ctx := context.Background()
	for i, r := range publication_writer_list {
		ts := time.Now().Format(time.RFC3339Nano)
		fmt.Fprintf(logFile, "%s inserting publication id %d\n", ts, r.pub_id)

		if err := db.insertSinglePublicationWriter(ctx, publication_writer_list[i]); err != nil {
			fmt.Fprintf(logFile, "%s error inserting publication/book id %d: %v\n", ts, r.pub_id, err)
			panic(err)
		}

		// throttle per second based on count_each_second
		if (i+1)%count_each_second == 0 && i+1 < len(publication_writer_list) {
			time.Sleep(time.Second)
		}
	}

}

func insertScenario3(db SeederConn) {

	start_id := 1
	end_id := 100
	count_each_second := 10

	publication_list := []Publication{}
	book_list := []Book{}

	for i := start_id; i <= end_id; i++ {

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

		book := Book{
			book_id:      i,
			book_isbn:    randomString(13, false, false, true),
			book_edit:    "Editorial " + randomString(10, false, false, false),
			book_vol:     randomString(10, false, false, false),
			book_editors: "Editors " + randomString(20, false, false, false),
		}

		book_list = append(book_list, book)
	}

	if count_each_second <= 0 {
		count_each_second = 1
	}

	logFile, err := os.OpenFile("scenario3.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	ctx := context.Background()
	for i, r := range publication_list {
		ts := time.Now().Format(time.RFC3339Nano)
		fmt.Fprintf(logFile, "%s inserting publication id %d\n", ts, r.pub_id)

		if err := db.insertPublicationWithBook(ctx, r, book_list[i]); err != nil {
			fmt.Fprintf(logFile, "%s error inserting publication/book id %d: %v\n", ts, r.pub_id, err)
			panic(err)
		}

		// throttle per second based on count_each_second
		if (i+1)%count_each_second == 0 && i+1 < len(publication_list) {
			time.Sleep(time.Second)
		}
	}

}

func insertScenario4(db SeederConn) {

	start_id := 1
	end_id := 100
	count_each_second := 10

	publication_list := []Publication{}
	conference_list := []Conference{}

	for i := start_id; i <= end_id; i++ {
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

	if count_each_second <= 0 {
		count_each_second = 1
	}

	logFile, err := os.OpenFile("scenario4.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	ctx := context.Background()
	for i, r := range publication_list {
		ts := time.Now().Format(time.RFC3339Nano)
		fmt.Fprintf(logFile, "%s inserting publication id %d\n", ts, r.pub_id)

		if err := db.insertPublicationWithConference(ctx, r, conference_list[i]); err != nil {
			fmt.Fprintf(logFile, "%s error inserting publication/conference id %d: %v\n", ts, r.pub_id, err)
			panic(err)
		}

		// throttle per second based on count_each_second
		if (i+1)%count_each_second == 0 && i+1 < len(publication_list) {
			time.Sleep(time.Second)
		}
	}
}

func insertScenario5(db SeederConn) {

	start_id := 1
	end_id := 100
	count_each_second := 10

	publication_list := []Publication{}
	journal_list := []Journal{}

	for i := start_id; i <= end_id; i++ {
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

	if count_each_second <= 0 {
		count_each_second = 1
	}

	logFile, err := os.OpenFile("scenario5.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	ctx := context.Background()
	for i, r := range publication_list {
		ts := time.Now().Format(time.RFC3339Nano)
		fmt.Fprintf(logFile, "%s inserting publication id %d\n", ts, r.pub_id)

		if err := db.insertPublicationWithJournal(ctx, r, journal_list[i]); err != nil {
			fmt.Fprintf(logFile, "%s error inserting publication/conference id %d: %v\n", ts, r.pub_id, err)
			panic(err)
		}

		// throttle per second based on count_each_second
		if (i+1)%count_each_second == 0 && i+1 < len(publication_list) {
			time.Sleep(time.Second)
		}
	}
}
