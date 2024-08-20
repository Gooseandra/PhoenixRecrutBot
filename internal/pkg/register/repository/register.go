package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"phenixRecrutBot/internal/pkg/errors"
	"phenixRecrutBot/internal/pkg/register/models"
	"time"
)

type RegisterRepo struct {
	db pgxpool.Pool
}

func NewRegisterRepo(pool pgxpool.Pool) RegisterRepo { return RegisterRepo{db: pool} }

func (r RegisterRepo) SetNewForm(form models.Form, id int64) error {
	query := `INSERT INTO forms (name, surname, phone, vk, tg, confirmed_name, time)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (tg) 
DO UPDATE 
SET 
    name = EXCLUDED.name,
    surname = EXCLUDED.surname,
    phone = EXCLUDED.phone,
    vk = EXCLUDED.vk,
    confirmed_name = EXCLUDED.confirmed_name,
    time = EXCLUDED.time;
`
	_, err := r.db.Exec(context.Background(), query, form.Name, form.Surname, form.Phone, form.Vk, form.Tg, form.ConfirmedName, time.Now())
	if err != nil {
		errors.Warn(id, fmt.Sprintf("can`t insert form with arguments Name:%s, Surname:%s, Phone:%s, Vk:%s, Tg:%s, Confirmation:%s",
			form.Name, form.Surname, form.Phone, form.Vk, form.Tg, form.ConfirmedName))
		return err
	}
	return nil
}

func (r RegisterRepo) List() ([]models.Form, error) {
	query := `select name, surname, phone, vk, tg from forms`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	var form models.Form
	var forms []models.Form
	for rows.Next() {
		err = rows.Scan(&form.Name, &form.Surname, &form.Phone, &form.Vk, &form.Tg)
		if err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}
	return forms, nil
}
