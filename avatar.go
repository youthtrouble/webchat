package main

import "errors"

//ErrNoAvatarURL is the error that occurs when the avatar instance is unable to provide a URL
var ErrNoAvatarURL = errors.New("Unable to get an avatar URL")

type Avatar interface {
	//GetAvatarURL gets the avatar url or returns an error
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}



