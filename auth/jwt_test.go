package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"example.com/sample/go_todo_app/clock"
	"example.com/sample/go_todo_app/entity"
	"example.com/sample/go_todo_app/store"
	"example.com/sample/go_todo_app/testutil/fixture"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}

	want = []byte("-----BEGIN RSA PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	ctx := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(20)
	u := fixture.User(&entity.User{ID: wantID})
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}

	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}

func TestJWTer_GetToken(t *testing.T) {
	t.Parallel()

	c := clock.NewFixedClocker()
	want, err := jwt.NewBuilder().JwtID(uuid.New().String()).Issuer("go_todo_app").Subject("access_token").IssuedAt(c.Now().UTC()).Expiration(c.Now().Add(30*time.Minute)).Claim(RoleKey, "test").Claim(UserNameKey, "test_user").Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	// アルゴリズムを確認
	alg, ok := pkey.Get(jwk.AlgorithmKey)
	if !ok || alg == "" {
		// アルゴリズムが設定されていない場合、明示的に設定
		if err := pkey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
			t.Fatal(err)
		}
	}
	// アルゴリズムを jwk.Key から取得
	keyAlg := pkey.Algorithm()

	// 取得したアルゴリズムを使用して署名
	signed, err := jwt.Sign(want, jwt.WithKey(keyAlg, pkey))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	userID := entity.UserID(20)

	ctx := context.Background()
	moq := &StoreMock{}
	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return userID, nil
	}

	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		`https://example.com/go-toto-app`,
		nil,
	)
	req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("want no err, but got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetToken() got = %v, want %v", got, want)
	}
}

type FixedTomorrowClocker struct{}

func (c FixedTomorrowClocker) Now() time.Time {
	return clock.FixedClocker{}.Now().Add(24 * time.Hour)
}

func TestJWTer_GetToken_NG(t *testing.T) {
	t.Parallel()

	c := clock.NewFixedClocker()
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).Issuer("go_todo_app").Subject("access_token").IssuedAt(c.Now()).Expiration(c.Now().Add(30*time.Minute)).Claim(RoleKey, "test").Claim(UserNameKey, "test_user").Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	// アルゴリズムを確認
	alg, ok := pkey.Get(jwk.AlgorithmKey)
	if !ok || alg == "" {
		// アルゴリズムが設定されていない場合、明示的に設定
		if err := pkey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
			t.Fatal(err)
		}
	}
	// アルゴリズムを jwk.Key から取得
	keyAlg := pkey.Algorithm()

	// 取得したアルゴリズムを使用して署名
	signed, err := jwt.Sign(tok, jwt.WithKey(keyAlg, pkey))
	if err != nil {
		t.Fatal(err)
	}

	type moq struct {
		userID entity.UserID
		err    error
	}
	tests := map[string]struct {
		c   clock.Clocker
		moq moq
	}{
		"expire": {
			c: FixedTomorrowClocker{},
		},
		"notFoundInStrore": {
			c: clock.FixedClocker{},
			moq: moq{
				err: store.ErrNotFound,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			moq := &StoreMock{}
			moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
				return tt.moq.userID, tt.moq.err
			}
			sut, err := NewJWTer(moq, tt.c)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(
				http.MethodGet,
				`https://example.com/go-toto-app`,
				nil,
			)
			req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
			got, err := sut.GetToken(ctx, req)
			if err == nil {
				t.Errorf("want error, but got nil")
			}
			if got != nil {
				t.Errorf("want nil, but got %v", got)
			}
		})
	}
}
