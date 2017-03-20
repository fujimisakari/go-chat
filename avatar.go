package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarはAvatarインスタンスがアバターのURLを
// 返すことができない場合に発生するエラーです
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

// Avatarはユーザーのプロフィール画像を表す型です。
type Avatar interface {
	// AvatarURLは指定されたクライアントのアバターのURLを返します。
	// 問題が発生した場合にはエラーを返します。
	// 特に、URLを取得できなかった場合にはErrNoAvatarURLを返します。
	AvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) AvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.AvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) AvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (_ GravatarAvatar) AvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) AvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
