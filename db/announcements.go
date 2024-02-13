package db

import "github.com/jmoiron/sqlx"

type AnnouncementManager struct {
	db *sqlx.DB
}

func (m *AnnouncementManager) GetList() {}

func (m *AnnouncementManager) GetOneById() {}

func (m *AnnouncementManager) CreateAnnounce() {}
func (m *AnnouncementManager) UpdateAnnounce() {}
func (m *AnnouncementManager) HideAnnounce()   {}
func (m *AnnouncementManager) DeleteAnnounce() {}
