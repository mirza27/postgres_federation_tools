package main

var CountryList = []Country{
	{cntry_id: "ES", cntry_name_es: "España", cntry_name_en: "Spain", cntry_name_gl: "España", cntry_aneca_id: 724},
	{cntry_id: "US", cntry_name_es: "Estados Unidos", cntry_name_en: "United States", cntry_name_gl: "Estados Unidos", cntry_aneca_id: 840},
	{cntry_id: "FR", cntry_name_es: "Francia", cntry_name_en: "France", cntry_name_gl: "Francia", cntry_aneca_id: 250},
	{cntry_id: "DE", cntry_name_es: "Alemania", cntry_name_en: "Germany", cntry_name_gl: "Alemaña", cntry_aneca_id: 276},
	{cntry_id: "IT", cntry_name_es: "Italia", cntry_name_en: "Italy", cntry_name_gl: "Italia", cntry_aneca_id: 380},
	{cntry_id: "GB", cntry_name_es: "Reino Unido", cntry_name_en: "United Kingdom", cntry_name_gl: "Reino Unido", cntry_aneca_id: 826},
	{cntry_id: "CN", cntry_name_es: "China", cntry_name_en: "China", cntry_name_gl: "China", cntry_aneca_id: 156},
	{cntry_id: "JP", cntry_name_es: "Japón", cntry_name_en: "Japan", cntry_name_gl: "Xapón", cntry_aneca_id: 392},
	{cntry_id: "IN", cntry_name_es: "India", cntry_name_en: "India", cntry_name_gl: "India", cntry_aneca_id: 356},
	{cntry_id: "BR", cntry_name_es: "Brasil", cntry_name_en: "Brazil", cntry_name_gl: "Brasil", cntry_aneca_id: 76},
}

type Country struct {
	cntry_id       string
	cntry_name_es  string
	cntry_name_en  string
	cntry_name_gl  string
	cntry_aneca_id int
}

type Researcher struct {
	rese_id            int
	rese_nif           string
	rese_name          string
	rese_first_surname string
	rese_sec_surname   string
	rese_birth_day     string
	rese_mail          string
	rese_web           string
	rese_phone         string
	rese_univ          string
	rese_sign          string
	rese_short_curr_es string
	rese_short_curr_en string
}

type Publication_writer struct {
	rese_id       int
	pub_id        int
	pwrit_is_edit bool
	pwrit_ord     int
}

type Publication struct {
	pub_id     int
	pub_isbn   string
	pub_title  string
	pub_year   int
	pub_loc    string
	pub_edit   string
	pub_subj   int
	pub_draft  string
	pub_prsnt  string
	pub_minute string
	pub_source string
	pub_book   string
	pub_link   string
	pub_repo   string
	pub_visib  bool
	pub_cntry  string
}

type Journal struct {
	jrnl_id          int
	jrnl_journ_title string
	jrnl_vol         string
	jrnl_numb        string
	jrnl_start_page  int
	jrnl_end_page    int
	jrnl_month       int
	jrnl_scope       int
	jrnl_indexed     bool
	jrnl_pos_numb    int
	jrnl_tot_pos     int
}

type Conference struct {
	conf_id         int
	conf_name       string
	conf_ref        string
	conf_publ       string
	conf_edit       string
	conf_start_date string
	conf_end_date   string
	conf_start_page int
	conf_end_page   int
	conf_type       int
	conf_serie      string
	conf_organi     string
	conf_particip   int
}

type Book struct {
	book_id      int
	book_isbn    string
	book_edit    string
	book_vol     string
	book_editors string
}

type Book_chapter struct {
	chapt_id         int
	chapt_book_tit   string
	chapt_vol        string
	chapt_start_page int
	chapt_end_page   int
	chapt_edit       string
	chapt_editors    string
}
