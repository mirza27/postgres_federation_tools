package main

import (
	"context"
	"db_migrate_server/internal/config"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeederConn struct {
	config config.Config
	Pool   *pgxpool.Pool
}

func NewSeederConn(config config.Config) (*SeederConn, error) {

	sourceDSN := "postgres://" + config.SourceDatabaseUser + ":" + config.SourceDatabasePass + "@" + config.SourceDatabaseHost + ":" +
		strconv.Itoa(config.SourceDatabasePort) + "/" + config.SourceDatabaseName

	pool, err := pgxpool.New(context.Background(), sourceDSN)
	if err != nil {
		return nil, err
	}

	return &SeederConn{
		Pool:   pool,
		config: config,
	}, nil
}

func (sd SeederConn) chunkInsertCountry(countries []Country, chunkSize int) {
	if len(countries) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(countries)
	}

	for start := 0; start < len(countries); start += chunkSize {
		end := start + chunkSize
		if end > len(countries) {
			end = len(countries)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, c := range countries[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (cntry_id) DO NOTHING`,
				c.cntry_id, c.cntry_name_es, c.cntry_name_en, c.cntry_name_gl, c.cntry_aneca_id,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

// insertSingleCountry inserts one country in its own transaction.
func (sd SeederConn) insertSingleCountry(ctx context.Context, c Country) error {
	tx, err := sd.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (cntry_id) DO NOTHING`,
		c.cntry_id, c.cntry_name_es, c.cntry_name_en, c.cntry_name_gl, c.cntry_aneca_id,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (sd SeederConn) chunkInsertResearcher(researchers []Researcher, chunkSize int) {
	if len(researchers) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(researchers)
	}

	for start := 0; start < len(researchers); start += chunkSize {
		end := start + chunkSize
		if end > len(researchers) {
			end = len(researchers)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, r := range researchers[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public.researcher (
					rese_id, rese_nif, rese_name, rese_first_surname, rese_sec_surname, rese_birth_day,
					rese_mail, rese_web, rese_phone, rese_univ, rese_sign, rese_short_curr_es,
					rese_short_curr_en, rese_gend, rese_pos
				 ) VALUES (
					$1, $2, $3, $4, $5, $6::date, $7, $8, $9, $10, $11, $12, $13, $14, $15
				 )
				 ON CONFLICT (rese_id) DO NOTHING`,
				r.rese_id, r.rese_nif, r.rese_name, r.rese_first_surname, r.rese_sec_surname,
				r.rese_birth_day, r.rese_mail, r.rese_web, r.rese_phone, r.rese_univ, r.rese_sign,
				r.rese_short_curr_es, r.rese_short_curr_en, 0, 0,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

// insertSingleResearcher inserts one researcher in its own transaction.
func (sd SeederConn) insertSingleResearcher(ctx context.Context, r Researcher) error {
	tx, err := sd.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.researcher (
			rese_id, rese_nif, rese_name, rese_first_surname, rese_sec_surname, rese_birth_day,
			rese_mail, rese_web, rese_phone, rese_univ, rese_sign, rese_short_curr_es,
			rese_short_curr_en, rese_gend, rese_pos
		 ) VALUES (
			$1, $2, $3, $4, $5, $6::date, $7, $8, $9, $10, $11, $12, $13, $14, $15
		 )
		 ON CONFLICT (rese_id) DO NOTHING`,
		r.rese_id, r.rese_nif, r.rese_name, r.rese_first_surname, r.rese_sec_surname,
		r.rese_birth_day, r.rese_mail, r.rese_web, r.rese_phone, r.rese_univ, r.rese_sign,
		r.rese_short_curr_es, r.rese_short_curr_en, 0, 0,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (sd SeederConn) chunkInsertPublication(publications []Publication, chunkSize int) {
	if len(publications) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(publications)
	}

	for start := 0; start < len(publications); start += chunkSize {
		end := start + chunkSize
		if end > len(publications) {
			end = len(publications)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, p := range publications[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public."publication" (
					pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt,
					pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry
				 ) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
				 )
				 ON CONFLICT (pub_id) DO NOTHING`,
				p.pub_id, p.pub_isbn, p.pub_title, p.pub_year, p.pub_loc, p.pub_edit, p.pub_subj,
				p.pub_draft, p.pub_prsnt, p.pub_minute, p.pub_source, p.pub_book, p.pub_link,
				p.pub_repo, p.pub_visib, p.pub_cntry,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

func (sd SeederConn) chunkInsertPublicationWriter(publicationWriters []Publication_writer, chunkSize int) {
	if len(publicationWriters) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(publicationWriters)
	}

	for start := 0; start < len(publicationWriters); start += chunkSize {
		end := start + chunkSize
		if end > len(publicationWriters) {
			end = len(publicationWriters)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, w := range publicationWriters[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public.publication_writer (rese_id, pub_id, pwrit_is_edit, pwrit_ord)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (rese_id, pub_id) DO NOTHING`,
				w.rese_id, w.pub_id, w.pwrit_is_edit, w.pwrit_ord,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

// insertSinglePublicationWriter inserts one publication_writer in its own transaction.
func (sd SeederConn) insertSinglePublicationWriter(ctx context.Context, w Publication_writer) error {
	tx, err := sd.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.publication_writer (rese_id, pub_id, pwrit_is_edit, pwrit_ord)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (rese_id, pub_id) DO NOTHING`,
		w.rese_id, w.pub_id, w.pwrit_is_edit, w.pwrit_ord,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (sd SeederConn) chunkInsertBook(books []Book, chunkSize int) {
	if len(books) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(books)
	}

	for start := 0; start < len(books); start += chunkSize {
		end := start + chunkSize
		if end > len(books) {
			end = len(books)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, b := range books[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public.book (book_id, book_isbn, book_edit, book_vol, book_editors)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (book_id) DO NOTHING`,
				b.book_id, b.book_isbn, b.book_edit, b.book_vol, b.book_editors,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

func (sd SeederConn) chunkInsertConference(conferences []Conference, chunkSize int) {
	if len(conferences) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(conferences)
	}

	for start := 0; start < len(conferences); start += chunkSize {
		end := start + chunkSize
		if end > len(conferences) {
			end = len(conferences)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, c := range conferences[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public.conference (
					conf_id, conf_name, conf_ref, conf_publ, conf_edit, conf_start_date, conf_end_date,
					conf_start_page, conf_end_page, conf_type, conf_serie, conf_organi, conf_particip
				 ) VALUES (
					$1, $2, $3, $4, $5, $6::date, $7::date, $8, $9, $10, $11, $12, $13
				 )
				 ON CONFLICT (conf_id) DO NOTHING`,
				c.conf_id, c.conf_name, c.conf_ref, c.conf_publ, c.conf_edit, c.conf_start_date,
				c.conf_end_date, c.conf_start_page, c.conf_end_page, c.conf_type, c.conf_serie,
				c.conf_organi, c.conf_particip,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

func (sd SeederConn) chunkInsertJournal(journals []Journal, chunkSize int) {
	if len(journals) == 0 {
		return
	}

	ctx := context.Background()
	if chunkSize <= 0 {
		chunkSize = len(journals)
	}

	for start := 0; start < len(journals); start += chunkSize {
		end := start + chunkSize
		if end > len(journals) {
			end = len(journals)
		}

		tx, err := sd.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}

		for _, j := range journals[start:end] {
			_, err = tx.Exec(ctx,
				`INSERT INTO public.journal (
					jrnl_id, jrnl_journ_title, jrnl_vol, jrnl_numb, jrnl_start_page, jrnl_end_page,
					jrnl_month, jrnl_scope, jrnl_indexed, jrnl_pos_numb, jrnl_tot_pos
				 ) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
				 )
				 ON CONFLICT (jrnl_id) DO NOTHING`,
				j.jrnl_id, j.jrnl_journ_title, j.jrnl_vol, j.jrnl_numb, j.jrnl_start_page,
				j.jrnl_end_page, j.jrnl_month, j.jrnl_scope, j.jrnl_indexed, j.jrnl_pos_numb,
				j.jrnl_tot_pos,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				fmt.Println("error inserting", err)
				panic(err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}
}

// insertPublicationWithBook inserts publication and book in one transaction.
func (sd SeederConn) insertPublicationWithBook(ctx context.Context, p Publication, b Book) error {
	tx, err := sd.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public."publication" (
			pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt,
			pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry
		 ) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		 )
		 ON CONFLICT (pub_id) DO NOTHING`,
		p.pub_id, p.pub_isbn, p.pub_title, p.pub_year, p.pub_loc, p.pub_edit, p.pub_subj,
		p.pub_draft, p.pub_prsnt, p.pub_minute, p.pub_source, p.pub_book, p.pub_link,
		p.pub_repo, p.pub_visib, p.pub_cntry,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.book (book_id, book_isbn, book_edit, book_vol, book_editors)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (book_id) DO NOTHING`,
		b.book_id, b.book_isbn, b.book_edit, b.book_vol, b.book_editors,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

// insertPublicationWithJournal inserts publication and journal in one transaction.
func (sd SeederConn) insertPublicationWithJournal(ctx context.Context, p Publication, j Journal) error {
	tx, err := sd.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public."publication" (
			pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt,
			pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry
		 ) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		 )
		 ON CONFLICT (pub_id) DO NOTHING`,
		p.pub_id, p.pub_isbn, p.pub_title, p.pub_year, p.pub_loc, p.pub_edit, p.pub_subj,
		p.pub_draft, p.pub_prsnt, p.pub_minute, p.pub_source, p.pub_book, p.pub_link,
		p.pub_repo, p.pub_visib, p.pub_cntry,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.journal (
			jrnl_id, jrnl_journ_title, jrnl_vol, jrnl_numb, jrnl_start_page, jrnl_end_page,
			jrnl_month, jrnl_scope, jrnl_indexed, jrnl_pos_numb, jrnl_tot_pos
		 ) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		 )
		 ON CONFLICT (jrnl_id) DO NOTHING`,
		j.jrnl_id, j.jrnl_journ_title, j.jrnl_vol, j.jrnl_numb, j.jrnl_start_page,
		j.jrnl_end_page, j.jrnl_month, j.jrnl_scope, j.jrnl_indexed, j.jrnl_pos_numb,
		j.jrnl_tot_pos,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

// insertPublicationWithConference inserts publication and conference in one transaction.
func (sd SeederConn) insertPublicationWithConference(ctx context.Context, p Publication, c Conference) error {
	tx, err := sd.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public."publication" (
			pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt,
			pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry
		 ) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		 )
		 ON CONFLICT (pub_id) DO NOTHING`,
		p.pub_id, p.pub_isbn, p.pub_title, p.pub_year, p.pub_loc, p.pub_edit, p.pub_subj,
		p.pub_draft, p.pub_prsnt, p.pub_minute, p.pub_source, p.pub_book, p.pub_link,
		p.pub_repo, p.pub_visib, p.pub_cntry,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.conference (
			conf_id, conf_name, conf_ref, conf_publ, conf_edit, conf_start_date, conf_end_date,
			conf_start_page, conf_end_page, conf_type, conf_serie, conf_organi, conf_particip
		 ) VALUES (
			$1, $2, $3, $4, $5, $6::date, $7::date, $8, $9, $10, $11, $12, $13
		 )
		 ON CONFLICT (conf_id) DO NOTHING`,
		c.conf_id, c.conf_name, c.conf_ref, c.conf_publ, c.conf_edit, c.conf_start_date,
		c.conf_end_date, c.conf_start_page, c.conf_end_page, c.conf_type, c.conf_serie,
		c.conf_organi, c.conf_particip,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}
