package tweet

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/nbd-wtf/go-nostr"
)

func PostNostr(text string) error {
	keyPath := filepath.Join(os.Getenv("USERPROFILE"), ".nostor")
	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("秘密鍵ファイルの読み込みに失敗しました: %v", err)
	}

	// 秘密鍵から公開鍵を算出
	pubKey, _ := nostr.GetPublicKey(string(key))

	// 2. 投稿内容（Event）の作成
	ev := nostr.Event{
		PubKey:    pubKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      nil,
		Content:   text,
	}

	// 3. イベントに署名
	err = ev.Sign(pubKey)
	if err != nil {
		return err
	}

	// 4. リレーに接続して送信
	ctx := context.Background()
	relayURL := "wss://relay.yabu.me"
	relay, err := nostr.RelayConnect(ctx, relayURL)
	if err != nil {
		return err
	}

	// イベントをリレーにPublish
	return relay.Publish(ctx, ev)
}