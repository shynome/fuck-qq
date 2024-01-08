package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/shynome/err0"
	"github.com/shynome/err0/try"
	"golang.org/x/sync/errgroup"
)

var testGroups []int64

func TestMain(m *testing.M) {
	groups := os.Getenv("FUCK_QQ_TEST_GROUPS")
	if groups != "" {
		groups := strings.Split(groups, ",")
		for _, g := range groups {
			g := try.To1(strconv.ParseInt(g, 10, 64))
			testGroups = append(testGroups, g)
		}
	}
	m.Run()
}

func TestMultiSend(t *testing.T) {
	eg := new(errgroup.Group)
	for _, group := range testGroups {
		group := group
		eg.Go(func() (err error) {
			defer err0.Then(&err, nil, nil)
			data := map[string]any{
				"group_id": group,
				"message":  `[CQ:at,qq=all] hello`,
			}
			body := try.To1(json.Marshal(data))
			req := try.To1(http.NewRequest(http.MethodPost, "http://127.0.0.1:5700/onebot/send_group_msg?check=1", bytes.NewReader(body)))
			req.Header.Add("Content-Type", "application/json")
			resp := try.To1(http.DefaultClient.Do(req))
			defer resp.Body.Close()
			rbody := try.To1(io.ReadAll(resp.Body))
			t.Log(rbody)
			return nil
		})
	}
	try.To(eg.Wait())
}
