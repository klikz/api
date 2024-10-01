package sqlstore

import (
	"errors"
)

func (r *Repo) CategoryGetAll() (interface{}, error) {

	type DataStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	rows, err := r.store.db.Query(`
		select c.id, c."name" from model_info.categories c 
		where c.status = true
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datas := []DataStruct{}

	for rows.Next() {
		var data DataStruct
		if err := rows.Scan(&data.ID, &data.Name); err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}
	if err = rows.Err(); err != nil {
		return datas, err
	}

	return datas, nil

}

func (r *Repo) CategoryAdd(name string, user_id int) error {
	_, err := r.store.db.Exec(`
	insert into model_info.categories ("name", c_user) values ($1, $2)
	`, name, user_id)
	if err != nil {
		return err
	}

	return nil

}

func (r *Repo) CategoryDelete(id, user_id int) error {
	_, err := r.store.db.Exec(`
	update model_info.categories set status = false, u_time = 'now()', u_user = $2
	where id = $1
	`, id, user_id)
	if err != nil {
		return err
	}

	return nil

}

func (r *Repo) ModelsGetAll(id int) (interface{}, error) {

	type DataStruct struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
	}

	rows, err := r.store.db.Query(`
	select m.id, m."name", c."name" as category
	from model_info.models m, model_info.categories c  
	where m.status = true
	and c.id = m.cat_id 
	and (c.id = $1 or (case when $1 in(0) then null else $1 end) is null )
	order by m.cat_id, m."name"
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datas := []DataStruct{}

	for rows.Next() {
		var data DataStruct
		if err := rows.Scan(&data.ID, &data.Name, &data.Category); err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}
	if err = rows.Err(); err != nil {
		return datas, err
	}

	return datas, nil

}

func (r *Repo) ModelAdd(name string, user_id, category_id int) error {

	_, err := r.store.db.Exec(`
	insert into model_info.models ("name", cat_id, c_user) values ($1, $2, $3)
	`, name, category_id, user_id)
	if err != nil {
		return err
	}

	return nil

}

func (r *Repo) ModelsDelete(id, user_id int) error {

	_, err := r.store.db.Exec(`
	update model_info.models set status = false, u_time = 'now()', u_user = $2
	where id = $1
	`, id, user_id)
	if err != nil {
		return err
	}

	return nil

}

func (r *Repo) InsertGsCode(gscode []string, model, user_id int) ([]string, error) {
	var badCode []string
	falseCode := false
	for _, code := range gscode {
		_, err := r.store.db.Exec(`insert into product.gs_code (gs_code, model_id, c_user_id) values ($1, $2, $3)`, code, model, user_id)
		if err != nil {
			badCode = append(badCode, code)
			falseCode = true
		}
	}
	if falseCode {
		return badCode, errors.New("code repeated")
	}

	return badCode, nil
}

func (r *Repo) GsCodeGetList() (interface{}, error) {

	type Report struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}

	rows, err := r.store.db.Query(`
	select m."name", count(g.id) as quantity 
	from product.gs_code g, model_info.models m  
	where g.gs_status = true and m.id = g.model_id 
	group by m."name"`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []Report{}

	for rows.Next() {
		var key Report
		if err := rows.Scan(&key.Name, &key.Quantity); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *Repo) GsCodeCheckCount(id, count int) error {

	check_couont := 0
	err := r.store.db.QueryRow(`
		select count(*) from product.gs_code gc 
	where gc.model_id = $1`, id).Scan(&check_couont)
	if err != nil {

		return err
	}
	if count > check_couont {
		return errors.New("gscode yetarli emas")
	}

	return nil
}

func (r *Repo) GsCodeGetGetId(id int) (int, string, error) {

	gc_id := 0
	gscodeString := ""
	err := r.store.db.QueryRow(`
	select gc.id, gc.gs_code  
	from product.gs_code gc 
	where gc.gs_status = true 
	and gc.product_id  is null
	and gc.model_id = $1
	limit 1`, id).Scan(&gc_id, &gscodeString)
	if err != nil {

		return gc_id, gscodeString, err
	}

	return gc_id, gscodeString, nil
}

func (r *Repo) GsCodeUpdate(product_id, user_id, model_id int) error {

	_, err := r.store.db.Exec(`
	update product.gs_code t
	set product_id = $1, 
	gs_status = false, 
	u_time = 'now()',
	u_user_id = $2
	where gs_status = true 
	and model_id = $3
	and t.id in (select gc.id from product.gs_code gc 
	where gc.model_id = $3
	and gc.gs_status = true
	limit 1 offset 0 )`, product_id, user_id, model_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GsCodeUpdateByID(product_id, user_id, gscodeId int) error {
	// println(fmt.Sprintf(`
	// update product.gs_code t
	// set t.product_id = %d,
	// t.gs_status = false,
	// t.u_time = 'now()',
	// t.u_user_id = %d
	// where t.id = %d`, product_id, user_id, gscodeId))

	_, err := r.store.db.Exec(`
	update product.gs_code 
	set product_id = $1,
	gs_status = false ,
	u_time = 'now()',
	u_user_id = $2
	where id = $3`, product_id, user_id, gscodeId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) StatusGetAll() (interface{}, error) {

	type Report struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	rows, err := r.store.db.Query(`
	select s.id, s."name" from model_info.status s 
	where s.status = true`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []Report{}

	for rows.Next() {
		var key Report
		if err := rows.Scan(&key.ID, &key.Name); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
