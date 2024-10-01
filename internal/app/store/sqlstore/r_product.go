package sqlstore

import (
	"errors"
	"fmt"
	"warehouse/internal/app/models"

	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
)

func (r *Repo) ProductAdd(category_id, model_id, user_id int, product *models.ProductInfo) error {

	err := r.store.db.QueryRow(`
	insert into product.products (catd_id, model_id, c_user, serial, status_id) 
	values ($1, $2, $3, $4, 1) returning id
	`, category_id, model_id, user_id, product.Serial).Scan(&product.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) ProductStatusListOutcome() (interface{}, error) {

	type Report struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Status   string `json:"status"`
	}

	rows, err := r.store.db.Query(`
	select m."name", count(p.id) as quantity, s."name" as status
	from product.products p, model_info.models m, model_info.status s 
	where p.status_id = 1 
	and m.id = p.model_id 
	and s.id = p.status_id 
	and p.t_status = true
	group by m."name", s."name"`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []Report{}

	for rows.Next() {
		var key Report
		if err := rows.Scan(&key.Name, &key.Quantity, &key.Status); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *Repo) ProductUpdateStatusToOutcome(user_id int, products []models.ProductInfo) (interface{}, error) {
	var BadSerial []string

	for i := 0; i < len(products); i++ {
		err := r.store.db.QueryRow(`
		update product.products 
		set status_id = 2,
		u_user = $1,
		u_time = 'now'
		where serial  = $2
		returning id
		`, user_id, products[i].Serial).Scan(&products[i].ID)

		if err != nil {
			BadSerial = append(BadSerial, products[i].Serial)
		}
	}

	if len(BadSerial) > 0 {
		return BadSerial, errors.New("seriya nomerda muammo")
	}

	return nil, nil
}

func (r *Repo) ProductUpdateStatusToIncome(user_id int, products []models.ProductInfo) (interface{}, error) {
	var BadSerial []string

	for i := 0; i < len(products); i++ {
		err := r.store.db.QueryRow(`
		update product.products 
		set status_id = 3,
		u_user = $1,
		u_time = 'now'
		where serial  = $2
		returning id
		`, user_id, products[i].Serial).Scan(&products[i].ID)

		if err != nil {
			BadSerial = append(BadSerial, products[i].Serial)
		}
	}

	if len(BadSerial) > 0 {
		return BadSerial, errors.New("seriya nomerda muammo")
	}

	return nil, nil
}

func (r *Repo) ProductStatusListIncome() (interface{}, error) {

	type Report struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Status   string `json:"status"`
	}

	rows, err := r.store.db.Query(`
	select m."name", count(p.id) as quantity, s."name" as status
	from product.products p, model_info.models m, model_info.status s 
	where p.status_id = 3 
	and m.id = p.model_id 
	and s.id = p.status_id 
	and p.t_status = true
	group by m."name", s."name"`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []Report{}

	for rows.Next() {
		var key Report
		if err := rows.Scan(&key.Name, &key.Quantity, &key.Status); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *Repo) ProductSerialInfo(serial string) (interface{}, error) {

	type Report struct {
		ID       int
		Serial   string `json:"serial"`
		Category string `json:"category"`
		Model    string `json:"model"`
		Status   string `json:"status"`
		GsCode   string `json:"gscode"`
		Time     string `json:"time"`
	}
	data := Report{}
	err := r.store.db.QueryRow(`
	select p.id, p.serial, c."name" as p_category, m."name" as p_model, 
	s."name" as p_status, gc.gs_code, 
	case when p.c_time is null then ' ' else to_char(p.c_time, 'DD-MM-YYYY HH24-MI')  end as p_time
	from product.gs_code gc, product.products p, model_info.categories c, model_info.models m, model_info.status s 
	where p.serial = $1
	and p.catd_id = c.id 
	and p.model_id = m.id 
	and p.status_id = s.id
	and gc.product_id = p.id`, serial).Scan(&data.ID, &data.Serial, &data.Category, &data.Model, &data.Status, &data.GsCode, &data.Time)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("bu serial nomer bo'yicha ma'lumot topilmadi")
		}
		return nil, err
	}

	return data, nil

}
func (r *Repo) ProductGeneratorBarcode(serial string) error {

	gscode := ""

	err := r.store.db.QueryRow(`select gc.gs_code  
	from product.gs_code gc, product.products p 
	where gc.product_id = p.id 
	and p.serial = $1`, serial).Scan(&gscode)
	if err != nil {
		return err
	}
	qrCode, err := datamatrix.Encode(gscode)
	if err != nil {
		return err
	}
	// qrCode, _ := qr.Encode(gscode, qr.H, qr.Unicode)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	// create the output fileC:\API\admin_ui\dist\assets
	file, _ := os.Create(fmt.Sprintf("C:/API/admin_ui/dist/assets/%s.png", serial))
	// file, _ := os.Create(fmt.Sprintf("D:/premier/import/ui/src/assets/%s.png", serial))
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return nil

}

func (r *Repo) GetPrintInfo(product *models.ProductInfo) error {

	err := r.store.db.QueryRow(`
	select p.id, c."name" as p_category, m."name" as p_model, gc.gs_code
	from model_info.categories c, model_info.models m, product.products p, product.gs_code gc 
	where c.id = p.catd_id 
	and m.id = p.model_id 
	and gc.product_id  = p.id 
	and p.serial = $1`, &product.Serial).Scan(&product.ID, &product.Category, &product.Model, &product.GsCode)
	if err != nil {
		return err
	}
	return nil

}

func (r *Repo) GetPrinterInfo(printerId int) (string, string, error) {

	printIP := ""
	printerName := ""
	if err := r.store.db.QueryRow(`
		select p.serverip, p."name" 
		from users.printers p 
		where p.id = $1
		`, printerId).Scan(&printIP, &printerName); err != nil {
		return printIP, printerName, err
	}

	fmt.Println("printerId: ", printerId)
	fmt.Println("printIP: ", printIP)
	fmt.Println("printerName: ", printerName)
	return printIP, printerName, nil
}

func (r *Repo) GetPrinterList() (interface{}, error) {

	type Report struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	rows, err := r.store.db.Query(`
	select p.id, p."name" 
	from users.printers p 
	order by p.id `)
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

func (r *Repo) ProductGetLast(count int) (interface{}, error) {

	type Report struct {
		ID       int    `json:"id"`
		Serial   string `json:"serial"`
		Category string `json:"category"`
		Model    string `json:"model"`
		Time     string `json:"time"`
	}

	rows, err := r.store.db.Query(`
	select p.id, p.serial, c."name" as category, m."name" as model, to_char(p.c_time, 'YYYY-MM-DD HH24:MI') as c_time
	from product.products p, model_info.models m, model_info.categories c  
	where c.id = p.catd_id
	and m.id = p.model_id 
	order by p.id desc 
	limit $1`, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []Report{}

	for rows.Next() {
		var key Report
		if err := rows.Scan(&key.ID, &key.Serial, &key.Category, &key.Model, &key.Time); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *Repo) ProductGeID(serial string) (int, error) {
	id := 0
	err := r.store.db.QueryRow(`
			select p.id from product.products p 
		where p.serial = $1`, serial).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *Repo) ProductGSCodeClear(product_id int) error {
	_, err := r.store.db.Exec(`
		update product.gs_code 
		set gs_status = true,
		u_time = null,
		u_user_id = null,
		product_id = null
		where product_id = $1
		`, product_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) ProductDelete(product_id int) error {
	_, err := r.store.db.Exec(`
		delete from product.products 
		where id = $1
		`, product_id)
	if err != nil {
		return err
	}

	return nil
}
