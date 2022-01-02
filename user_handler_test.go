package userapp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/nsf/jsondiff"
)

func init() {
	dbEndpoint := "http://dynamodb-local:8000"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           "local",
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Endpoint:   aws.String(dbEndpoint),
			DisableSSL: aws.Bool(true),
		},
	}))
	gdb = dynamo.New(sess)
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name           string
		input          func(t *testing.T)
		wantStatusCode int
		want           string
	}{
		{
			name: "複数件のユーザの取得",
			input: func(t *testing.T) {
				err := gdb.CreateTable(usersTable, User{}).Provision(1, 1).Run()
				if err != nil {
					t.Errorf("dynamo create table %s: %v", usersTable, err)
				}
				inputUsers := []User{{UserID: "001", UserName: "gopher"}, {UserID: "002", UserName: "rubyist"}}
				for _, u := range inputUsers {
					if err := gdb.Table(usersTable).Put(u).Run(); err != nil {
						t.Errorf("dynamo input user %v: %v", u, err)
					}
				}
			},
			wantStatusCode: 200,
			want:           "./testdata/want_get_users_1.json",
		},
		{
			name: "ユーザ0件",
			input: func(t *testing.T) {
				err := gdb.CreateTable(usersTable, User{}).Provision(1, 1).Run()
				if err != nil {
					t.Errorf("dynamo create table %s: %v", usersTable, err)
				}
			},
			wantStatusCode: 200,
			want:           "./testdata/want_get_users_2.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input(t)
			t.Cleanup(func() {
				if err := gdb.Table(usersTable).DeleteTable().Run(); err != nil {
					t.Fatalf("dynamo delete table %s: %v", usersTable, err)
				}
			})

			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.Request, _ = http.NewRequest("GET", "/users", nil)
			GetUsers(context)

			want, err := ioutil.ReadFile(tt.want)
			if err != nil {
				t.Fatalf("want file read: %v", err)
			}

			if response.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("status got %v, but want %v", response.Result().StatusCode, tt.wantStatusCode)
			}

			opt := jsondiff.DefaultConsoleOptions()
			if d, s := jsondiff.Compare(response.Body.Bytes(), want, &opt); d != jsondiff.FullMatch {
				t.Errorf("unmatch, got=%s, want=%s, diff=%s", response.Body.String(), string(want), s)
			}
		})
	}
}

func TestPostUsers(t *testing.T) {
	tests := []struct {
		name           string
		input          func(t *testing.T)
		wantStatusCode int
		want           string
	}{
		{
			name: "ユーザの取得成功",
			input: func(t *testing.T) {
				err := gdb.CreateTable(usersTable, User{}).Provision(1, 1).Run()
				if err != nil {
					t.Errorf("dynamo create table %s: %v", usersTable, err)
				}
			},
			wantStatusCode: 200,
			want:           "./testdata/want_post_users_1.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input(t)
			t.Cleanup(func() {
				if err := gdb.Table(usersTable).DeleteTable().Run(); err != nil {
					t.Fatalf("dynamo delete table %s: %v", usersTable, err)
				}
			})

			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			body := bytes.NewBufferString("{\"user_id\":\"003\",\"name\":\"ferris\"}")
			context.Request, _ = http.NewRequest("POST", "/users", body)
			PostUsers(context)

			want, err := ioutil.ReadFile(tt.want)
			if err != nil {
				t.Fatalf("want file read: %v", err)
			}

			if response.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("status got %v, but want %v", response.Result().StatusCode, tt.wantStatusCode)
			}

			opt := jsondiff.DefaultConsoleOptions()
			if d, s := jsondiff.Compare(response.Body.Bytes(), want, &opt); d != jsondiff.FullMatch {
				t.Errorf("unmatch, got=%s, want=%s, diff=%s", response.Body.String(), string(want), s)
			}
		})
	}
}
