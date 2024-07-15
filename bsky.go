package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bluesky-social/indigo/api/atproto"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	encoding "github.com/mattn/go-encoding"
	"golang.org/x/net/html/charset"
)

// テキストから抽出されたタグやリンクを表す構造体です。
type entry struct {
	start int64
	end   int64
	text  string
}

// タグとURLを識別するための正規表現パターン
var (
	tagRe = regexp.MustCompile(`\B#\S+`)
	urlRe = regexp.MustCompile(`https?://\S+`)
)

type config struct {
	Identifier string
	Password   string
}

func postBsky(text string) error {
	// Bluesky APIにアクセスするためのクライアントを初期化します。
	cli := &xrpc.Client{
		Host: "https://bsky.social",
	}

	// セッションを作成するための入力データを準備します。

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	input := &atproto.ServerCreateSession_Input{
		Identifier: cfg.Identifier, // Blueskyのハンドル名(xxxxx.bsky.social)
		Password:   cfg.Password,   // Blueskyのパスワード
	}
	// セッション作成のリクエストを送信し、結果を受け取ります。
	output, err := atproto.ServerCreateSession(context.TODO(), cli, input)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	// 認証情報をクライアントに設定します。
	cli.Auth = &xrpc.AuthInfo{
		AccessJwt:  output.AccessJwt,
		RefreshJwt: output.RefreshJwt,
		Handle:     output.Handle,
		Did:        output.Did,
	}

	// ここで投稿データを作成します
	post := &bsky.FeedPost{
		Text:      text,                            // 投稿テキスト
		CreatedAt: time.Now().Format(time.RFC3339), // 投稿日時
		Langs:     []string{"ja"},                  // 言語設定
	}

	// テキストからタグを抽出し、投稿データに追加します。
	for _, entry := range extractTagsBytes(text) {
		post.Facets = append(post.Facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
						Tag: entry.text,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: entry.start,
				ByteEnd:   entry.end,
			},
		})
	}

	// テキストからリンクを抽出し、投稿データに追加します。
	for _, entry := range extractLinksBytes(text) {
		post.Facets = append(post.Facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Link: &bsky.RichtextFacet_Link{
						Uri: entry.text,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: entry.start,
				ByteEnd:   entry.end,
			},
		})
		if post.Embed == nil {
			post.Embed = &bsky.FeedPost_Embed{}
		}
		if post.Embed.EmbedExternal == nil {
			addLink(cli, post, entry.text)
		}
	}

	inp := &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       cli.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: post,
		},
	}
	// 投稿データをBlueskyに送信します。
	ctx := context.Background()
	_, err = atproto.RepoCreateRecord(ctx, cli, inp)
	if err != nil {
		return fmt.Errorf("Error posting to bluesky: %w", err)
	}
	return nil
}

// 投稿テキストからタグを抽出する関数です。
func extractTagsBytes(text string) []entry {
	var result []entry
	matches := tagRe.FindAllStringSubmatchIndex(text, -1)
	for _, m := range matches {
		result = append(result, entry{
			text:  strings.TrimPrefix(text[m[0]:m[1]], "#"),
			start: int64(len(text[0:m[0]])),
			end:   int64(len(text[0:m[1]]))},
		)
	}
	return result
}

// 投稿テキストからリンクを抽出する関数です。
func extractLinksBytes(text string) []entry {
	var result []entry
	matches := urlRe.FindAllStringSubmatchIndex(text, -1)
	for _, m := range matches {
		result = append(result, entry{
			text:  text[m[0]:m[1]],
			start: int64(len(text[0:m[0]])),
			end:   int64(len(text[0:m[1]]))},
		)
	}
	return result
}

// 投稿データに外部リンクの詳細（タイトル、説明、サムネイル画像）を追加する関数です。
func addLink(xrpcc *xrpc.Client, post *bsky.FeedPost, link string) {
	res, _ := http.Get(link)
	if res == nil {
		return
	}
	defer res.Body.Close()

	br := bufio.NewReader(res.Body)
	var reader io.Reader = br

	data, err := br.Peek(1024)
	if err == nil {
		enc, name, _ := charset.DetermineEncoding(data, res.Header.Get("content-type"))
		if enc != nil {
			reader = enc.NewDecoder().Reader(br)
		} else if len(name) > 0 {
			enc := encoding.GetEncoding(name)
			if enc != nil {
				reader = enc.NewDecoder().Reader(br)
			}
		}
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		post.Embed.EmbedExternal = &bsky.EmbedExternal{
			External: &bsky.EmbedExternal_External{
				Uri: link,
			},
		}
		return
	}

	title := doc.Find(`title`).Text()
	if title == "" {
		title, _ = doc.Find(`meta[property="og:title"]`).Attr("content")
		if title == "" {
			title = link
		}
	}
	description, _ := doc.Find(`meta[property="description"]`).Attr("content")
	if description == "" {
		description, _ = doc.Find(`meta[property="og:description"]`).Attr("content")
		if description == "" {
			description = link
		}
	}
	post.Embed.EmbedExternal = &bsky.EmbedExternal{
		External: &bsky.EmbedExternal_External{
			Description: description,
			Title:       title,
			Uri:         link,
		},
	}

	imgURL, _ := doc.Find(`meta[property="og:image"]`).Attr("content")
	if imgURL == "" || post.Embed.EmbedExternal == nil {
		return
	}
	resp, err := http.Get(imgURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp2, err := comatproto.RepoUploadBlob(context.TODO(), xrpcc, bytes.NewReader(b))
	if err != nil {
		return
	}
	post.Embed.EmbedExternal.External.Thumb = &lexutil.LexBlob{
		Ref:      resp2.Blob.Ref,
		MimeType: http.DetectContentType(b),
		Size:     resp2.Blob.Size,
	}
}

func loadConfig() (*config, error) {
	var configFile = path.Join(os.Getenv("USERPROFILE"), ".bluesky.json")
	f, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load config file: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var cfg config
	dec.Decode(&cfg)
	return &cfg, nil
}
