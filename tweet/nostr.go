package tweet

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func PostNostr(text string) error {
	// 1. 秘密鍵ファイルの読み込みとクレンジング
	keyPath := filepath.Join(os.Getenv("USERPROFILE"), ".nostr")
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("秘密鍵ファイルの読み込みに失敗しました: %w", err)
	}

	// 改行コードや余計な空白を除去
	rawKey := strings.TrimSpace(string(keyBytes))

	var sk string

	// 2. nsec形式かHex形式かを判別してデコード
	if strings.HasPrefix(rawKey, "nsec1") {
		var prefix string
		var value any
		// nsec1... の場合はNIP-19デコードを行う
		prefix, value, err = nip19.Decode(rawKey)
		if err != nil {
			return fmt.Errorf("nsecのデコードに失敗しました: %w", err)
		}

		// デコード結果が秘密鍵（string型）であることを確認
		var ok bool
		sk, ok = value.(string)
		if !ok || prefix != "nsec" {
			return fmt.Errorf("不正な秘密鍵形式です (prefix: %s)", prefix)
		}
	} else {
		// すでにHex形式の場合はそのまま使用
		sk = rawKey
	}

	// 3. 秘密鍵から公開鍵を算出
	pubKey, err := nostr.GetPublicKey(sk)
	if err != nil {
		return fmt.Errorf("公開鍵の算出に失敗しました: %w", err)
	}

	// 4. 投稿内容（Event）の作成
	ev := nostr.Event{
		PubKey:    pubKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      nil,
		Content:   text,
	}

	// 5. イベントに署名 (デコードした秘密鍵 sk を渡す)
	err = ev.Sign(sk)
	if err != nil {
		return fmt.Errorf("イベントの署名に失敗しました: %w", err)
	}

	// 6. リレーに接続して送信
	ctx := context.Background()
	relayURL := "wss://relay-jp.nostr.wirednet.jp"
	relay, err := nostr.RelayConnect(ctx, relayURL)
	if err != nil {
		return fmt.Errorf("リレーへの接続に失敗しました: %w", err)
	}
	defer relay.Close() // 接続を適切に閉じる

	// イベントをリレーにPublish
	if err := relay.Publish(ctx, ev); err != nil {
		return fmt.Errorf("イベントの送信に失敗しました: %w", err)
	}

	return nil
}