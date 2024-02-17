package entities

import "announce-api/utils"

func (u *InputSignInUser) Validate() error {
	if err := utils.ValidateStruct(u); err != nil {
		return err
	}
	return nil
}

func (u *InputSignUpUser) Validate() error {
	if err := utils.ValidateStruct(u); err != nil {
		return err
	}
	return nil
}

func (a *InputAnnouncement) Validate() error {
	if err := utils.ValidateStruct(a); err != nil {
		return err
	}
	return nil
}

func (a *AnnouncementForDB) Validate() error {
	if err := utils.ValidateStruct(a); err != nil {
		return err
	}
	return nil
}
