package services

import "announce-api/db"

type AnnouncementManager struct {
	client *db.DatabaseClient
}

func (m *AnnouncementManager) GetList()        {}
func (m *AnnouncementManager) GetOneById()     {}
func (m *AnnouncementManager) CreateAnnounce() {}
func (m *AnnouncementManager) UpdateAnnounce() {}
func (m *AnnouncementManager) HideAnnounce()   {}
func (m *AnnouncementManager) DeleteAnnounce() {}