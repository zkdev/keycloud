package main

import (
	"encoding/hex"
	"fmt"
	"webauthn/webauthn"
)

type Storage struct {
	users          map[string]*User
	authenticators map[string]*Authenticator
}

func (s *Storage) AddAuthenticator(user webauthn.User, authenticator webauthn.Authenticator) error {
	authr := &Authenticator{
		ID:           authenticator.WebAuthID(),
		CredentialID: authenticator.WebAuthCredentialID(),
		PublicKey:    authenticator.WebAuthPublicKey(),
		AAGUID:       authenticator.WebAuthAAGUID(),
		SignCount:    authenticator.WebAuthSignCount(),
	}
	key := hex.EncodeToString(authr.ID)

	u, ok := s.users[string(user.WebAuthID())]
	if !ok {
		return fmt.Errorf("user not found")
	}

	if _, ok := s.authenticators[key]; ok {
		return fmt.Errorf("authenticator already exists")
	}

	authr.User = u

	u.Authenticators[key] = authr
	s.authenticators[key] = authr

	return nil
}

func (s *Storage) GetAuthenticator(id []byte) (webauthn.Authenticator, error) {
	authr, ok := s.authenticators[hex.EncodeToString(id)]
	if !ok {
		return nil, fmt.Errorf("authenticator not found")
	}
	return authr, nil
}

func (s *Storage) GetUser(webauthnID string) (user *User) {
	u, ok := s.users[webauthnID]
	if !ok {
		return nil
	} else {
		return u
	}
}

func (s *Storage) GetAuthenticators(user webauthn.User) ([]webauthn.Authenticator, error) {
	u, ok := s.users[string(user.WebAuthID())]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	var authrs []webauthn.Authenticator
	for _, v := range u.Authenticators {
		authrs = append(authrs, v)
	}
	return authrs, nil
}

/*
DUMMY IMPLEMENTATION
*/
// TODO implement session functions
func (s *Storage) GetSessionKeyForUser(user *User) []byte {
	return []byte("abc")
}

func (s *Storage) SetSessionKeyForUser(user *User, b []byte) error {
	return nil
}

func (s *Storage) DeleteSessionKeyForUser(user *User) error {
	return nil
}

func (s *Storage) AddUser(u *User) (*User, error) {
	s.users[string(u.WebAuthID())] = u
	queryAddUser(connectDatabase(), u)
	return u, nil
}

func (s *Storage) RemoveUser(u *User) (bool, error) {
	s.users[string(u.WebAuthID())] = u
	queryRemoveUser(connectDatabase(), u)
	// TODO return bool depending on user got removed
	return true, nil
}

func (s *Storage) UpdateUser(u *User) (*User, error) {
	s.users[string(u.WebAuthID())] = u
	queryUpdateUser(connectDatabase(), u)
	return u, nil
}

func (s *Storage) AddPassword(u *User, st string, p *Password) (*Password, error) {
	s.users[string(u.WebAuthID())] = u
	queryAddPassword(connectDatabase(), u, p)
	return p, nil
}

func (s *Storage) GetPassword(u *User, st string) (*Password, error) {
	// TODO remove mockup url
	url := "http://keycloud.zeekay.dev/"
	s.users[string(u.WebAuthID())] = u
	rs := queryGetPassword(connectDatabase(), u, url)
	// TODO parse rs to passwd
	passwd := "mockup_passwd"
	id := "mockup_id"

	return &Password{
		Password: passwd,
		Id:       id,
	}, nil
}

func (s *Storage) UpdatePassword(u *User, st string, p *Password) (*Password, error) {
	s.users[string(u.WebAuthID())] = u
	queryUpdatePassword(connectDatabase(), u, p)
	return p, nil
}

func (s *Storage) DeletePassword(u *User, p *Password) (bool, error) {
	s.users[string(u.WebAuthID())] = u
	queryDeletePassword(connectDatabase(), p)
	// TODO return bool depending on passwd got removed
	return true, nil
}
