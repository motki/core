package model

import (
	"errors"
	"strings"
)

type MailingListSubscriber struct {
	Name  string
	Email string
}

func (m *Manager) GetMailingList(key string) ([]*MailingListSubscriber, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	res := []*MailingListSubscriber{}
	rows, err := db.Query("SELECT name, email FROM app.mailing_lists WHERE key = $1", key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rec := &MailingListSubscriber{}
		err := rows.Scan(&rec.Name, &rec.Email)
		if err != nil {
			return nil, err
		}
		res = append(res, rec)
	}
	return res, nil
}

func (m *Manager) AddToMailingList(key string, rec MailingListSubscriber) error {
	if !strings.Contains(rec.Email, "@") {
		return errors.New("invalid email address")
	}
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`INSERT INTO app.mailing_lists (key, name, email) VALUES($1, $2, $3) ON CONFLICT ON CONSTRAINT "mailing_lists_pkey" DO NOTHING`, key, rec.Name, rec.Email)
	if err != nil {
		return err
	}
	return nil
}