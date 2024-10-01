package sqlstore

import (
	"fmt"
	"warehouse/internal/app/models"
)

type Repo struct {
	store *Store
}

func (r *Repo) Create(u *models.User) error {
	if err := r.store.db.QueryRow(`insert into users.users ("name", "password") values ($1, $2) returning id`, u.UserName, u.EncryptedPassword).Scan(&u.ID); err != nil {
		return err
	}
	_, err := r.store.db.Exec(`insert into users.routes (user_id) values ($1)`, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdatePassword(u *models.User) error {
	_, err := r.store.db.Exec(`update users.users set "password" = $1 where "name" = $2`, u.EncryptedPassword, u.UserName)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindByUserName(u *models.User) error {
	fmt.Println(u.UserName)
	if err := r.store.db.QueryRow(`select u."password"  from users.users u 
	where u."name" = $1`, u.UserName).Scan(&u.EncryptedPassword); err != nil {
		return err
	}
	return nil
}

func (r *Repo) CheckRole(route string, id int) (bool, error) {
	// select r."%s" as result from routes r where "user" = '%s'
	result := false
	err := r.store.db.QueryRow(fmt.Sprintf(`
	select r."%s" as result from users.routes r where r.user_id = '%d' 
	 `, route, id)).Scan(&result)
	if err != nil {
		if err.Error() == `pq: column r.`+route+` does not exist` {
			r.store.db.Exec(`
			ALTER TABLE users.routes ADD "` + route + `" bool NOT NULL DEFAULT false; `)
		}
		return false, err
	}
	return result, nil
}

func (r *Repo) GetUserID(username string) (int, error) {
	id := 0
	if err := r.store.db.QueryRow(`select u.id from users.users u where u."name" = $1`, username).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r *Repo) GetUserCheckPoint(id int) (int, error) {
	check_id := 0
	if err := r.store.db.QueryRow(`select b.checkpoint_id from brigadir b where b.user_id = $1`, id).Scan(&check_id); err != nil {
		return check_id, err
	}
	return check_id, nil
}
