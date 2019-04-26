package main

import "testing"

func TestAuthAvatar(t *testing.T)  {
	var authAvatar AuthAvatar
	client := new(client)
	url,err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL{
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLは "+ "ErrNoAvatarURLを返すべきです")
	}
	testUrl := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url":testUrl}
	url,err = authAvatar.GetAvatarURL(client)
	if err != nil{
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLは" + "エラーを返すべきではありません")
	}else {
		if url != testUrl{
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T)  {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url,err := gravatarAvatar.GetAvatarURL(client)
	if err != nil{
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"{
		t.Errorf("GravatarAvatar.GetAvatarURLが%sというあやまった値を返しました",url)
	}
}