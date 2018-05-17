package model

import (
	"strings"

	"github.com/pkg/errors"
)

type MailingListSubscriber struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MailManager struct {
	bootstrap
}

func newMailManager(m bootstrap) *MailManager {
	return &MailManager{m}
}

func (m *MailManager) GetMailingList(key string) ([]*MailingListSubscriber, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	var res []*MailingListSubscriber
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

func (m *MailManager) AddToMailingList(key string, rec MailingListSubscriber) error {
	if !strings.Contains(rec.Email, "@") {
		return errors.New("invalid email address")
	}
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(`INSERT INTO app.mailing_lists (key, name, email) VALUES($1, $2, $3) ON CONFLICT ON CONSTRAINT "mailing_lists_pkey" DO NOTHING`, key, rec.Name, rec.Email)
	if err != nil {
		return err
	}
	return nil
}
