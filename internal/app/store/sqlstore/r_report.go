package sqlstore

func (r *Repo) ReportGetByModelsC(date1, date2, serial string, category_id, model_id, status_id int) (interface{}, error) {

	type ByModel struct {
		Category string `json:"category"`
		Model    string `json:"model"`
		Count    string `json:"count"`
		Status   string `json:"status"`
	}
	var byModel []ByModel

	rows, err := r.store.db.Query(`
	select c."name" as category, m."name" as model, COUNT(*), s."name" as status_name
	FROM product.products p, model_info.models m, model_info.categories c, model_info.status s
	where p.c_time between to_timestamp((case when $1  in('')  then '2022-01-01' else $1 end), 'YYYY-MM-DD HH24:MI') and to_timestamp((case when $2 in('') then (to_char(now(), 'YYYY-MM-DD HH24:MI')) else $2 end), 'YYYY-MM-DD HH24:MI')
	and c.id = p.catd_id 
	and m.id = p.model_id
	and s.id = p.status_id
	and (p.catd_id = $3 or (case when $3 in(0) then null else $3 end) is null )
	and (p.model_id = $4 or (case when $4 in(0) then null else $4 end) is null )
	and (p.status_id = $5 or (case when $5 in(0) then null else $5 end) is null )
	and (p.serial = $6 or (case when $6 in('') then null else $6 end) is null )
	group by m."name", p.model_id, c."name", s."name"`, date1, date2, category_id, model_id, status_id, serial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp ByModel
		if err := rows.Scan(&comp.Category,
			&comp.Model,
			&comp.Count,
			&comp.Status); err != nil {
			return byModel, err
		}
		byModel = append(byModel, comp)
	}
	if err = rows.Err(); err != nil {
		return byModel, err
	}

	return byModel, nil
}

func (r *Repo) ReportGetBySerialsC(date1, date2, serial string, category_id, model_id, status_id int) (interface{}, error) {

	type ByModel struct {
		Category string `json:"category"`
		Model    string `json:"model"`
		Serial   string `json:"serial"`
		Status   string `json:"status"`
		Time     string `json:"time"`
		GsCode   string `json:"gscode"`
	}
	var byModel []ByModel

	rows, err := r.store.db.Query(`
	select c."name" as category, m."name" as model, p.serial, s."name" as status_name, to_char(p.c_time, 'DD-MM-YYYY HH24-MI') as r_time, gc.gs_code 
	FROM product.products p, model_info.models m, model_info.categories c, model_info.status s, product.gs_code gc
	where p.c_time between to_timestamp((case when $1  in('')  then '2022-01-01' else $1 end), 'YYYY-MM-DD HH24:MI') and to_timestamp((case when $2 in('') then (to_char(now(), 'YYYY-MM-DD HH24:MI')) else $2 end), 'YYYY-MM-DD HH24:MI')
	and c.id = p.catd_id 
	and m.id = p.model_id
	and s.id = p.status_id
	and gc.product_id = p.id 
	and (p.catd_id = $3 or (case when $3 in(0) then null else $3 end) is null )
	and (p.model_id = $4 or (case when $4 in(0) then null else $4 end) is null )
	and (p.status_id = $5 or (case when $5 in(0) then null else $5 end) is null )
	and (p.serial = $6 or (case when $6 in('') then null else $6 end) is null )
	order by p.catd_id, p.model_id, p.serial, p.status_id`, date1, date2, category_id, model_id, status_id, serial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp ByModel
		if err := rows.Scan(&comp.Category,
			&comp.Model,
			&comp.Serial,
			&comp.Status,
			&comp.Time,
			&comp.GsCode); err != nil {
			return byModel, err
		}
		byModel = append(byModel, comp)
	}
	if err = rows.Err(); err != nil {
		return byModel, err
	}

	return byModel, nil
}

func (r *Repo) ReportGetByModelsU(date1, date2, serial string, category_id, model_id, status_id int) (interface{}, error) {

	type ByModel struct {
		ID       string `json:"id"`
		Category string `json:"category"`
		Model    string `json:"model"`
		Count    string `json:"count"`
		Status   string `json:"status"`
		Time     string `json:"time"`
		GsCode   string `json:"gscode"`
	}
	var byModel []ByModel

	rows, err := r.store.db.Query(`
	select p.id, c."name" as category, m."name" as model, COUNT(*), s."name" as status_name, 
	case when p.u_time is null then ' ' else to_char(p.u_time, 'DD-MM-YYYY HH24-MI') as r_time end,
	gc.gs_code
	FROM product.products p, model_info.models m, model_info.categories c, model_info.status s
	where p.u_time between to_timestamp((case when $1  in('')  then '2022-01-01' else $1 end), 'YYYY-MM-DD HH24:MI') and to_timestamp((case when $2 in('') then (to_char(now(), 'YYYY-MM-DD HH24:MI')) else $2 end), 'YYYY-MM-DD HH24:MI')
	and c.id = p.catd_id 
	and m.id = p.model_id
	and s.id = p.status_id
	and (p.catd_id = $3 or (case when $3 in(0) then null else $3 end) is null )
	and (p.model_id = $4 or (case when $4 in(0) then null else $4 end) is null )
	and (p.status_id = $5 or (case when $5 in(0) then null else $5 end) is null )
	and (p.serial = $6 or (case when $6 in('') then null else $6 end) is null )
	group by m."name", p.model_id, c."name", s."name"`, date1, date2, category_id, model_id, status_id, serial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp ByModel
		if err := rows.Scan(&comp.ID,
			&comp.Category,
			&comp.Model,
			&comp.Count,
			&comp.Status,
			&comp.Time,
			&comp.GsCode); err != nil {
			return byModel, err
		}
		byModel = append(byModel, comp)
	}
	if err = rows.Err(); err != nil {
		return byModel, err
	}

	return byModel, nil
}

func (r *Repo) ReportGetBySerialsU(date1, date2, serial string, category_id, model_id, status_id int) (interface{}, error) {

	type ByModel struct {
		ID       string `json:"id"`
		Category string `json:"category"`
		Model    string `json:"model"`
		Serial   string `json:"serial"`
		Status   string `json:"status"`
		Time     string `json:"time"`
		GsCode   string `json:"gscode"`
	}
	var byModel []ByModel

	rows, err := r.store.db.Query(`
	select p.id, c."name" as category, m."name" as model, p.serial, s."name" as status_name, 
	case when p.c_time is null then ' ' else to_char(p.c_time, 'DD-MM-YYYY HH24-MI') as r_time end, 
	gc.gs_code 
	FROM product.products p, model_info.models m, model_info.categories c, model_info.status s, product.gs_code gc
	where p.u_time between to_timestamp((case when $1  in('')  then '2022-01-01' else $1 end), 'YYYY-MM-DD HH24:MI') and to_timestamp((case when $2 in('') then (to_char(now(), 'YYYY-MM-DD HH24:MI')) else $2 end), 'YYYY-MM-DD HH24:MI')
	and c.id = p.catd_id 
	and m.id = p.model_id
	and s.id = p.status_id
	and gc.product_id = p.id 
	and (p.catd_id = $3 or (case when $3 in(0) then null else $3 end) is null )
	and (p.model_id = $4 or (case when $4 in(0) then null else $4 end) is null )
	and (p.status_id = $5 or (case when $5 in(0) then null else $5 end) is null )
	and (p.serial = $6 or (case when $6 in('') then null else $6 end) is null )
	order by p.catd_id, p.model_id, p.serial, p.status_id`, date1, date2, category_id, model_id, status_id, serial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp ByModel
		if err := rows.Scan(&comp.ID,
			&comp.Category,
			&comp.Model,
			&comp.Serial,
			&comp.Status,
			&comp.Time,
			&comp.GsCode); err != nil {
			return byModel, err
		}
		byModel = append(byModel, comp)
	}
	if err = rows.Err(); err != nil {
		return byModel, err
	}

	return byModel, nil
}
