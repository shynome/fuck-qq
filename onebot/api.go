package onebot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/docker/go-units"
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/shynome/err0"
	"github.com/shynome/err0/try"
	"github.com/shynome/fuck-qq/onebot/msg"
)

var rw sync.Locker = &sync.Mutex{}

var fileLoader = ttlcache.LoaderFunc[string, []byte](func(c *ttlcache.Cache[string, []byte], k string) *ttlcache.Item[string, []byte] {
	var err error
	logger := slog.With(
		"store", "file loader",
		"link", k,
	)
	defer err0.Then(&err, nil, func() {
		logger.Warn("cache failed")
	})
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req := try.To1(http.NewRequestWithContext(ctx, http.MethodGet, k, http.NoBody))
	resp := try.To1(http.DefaultClient.Do(req))
	defer resp.Body.Close()
	if code := resp.StatusCode; code != http.StatusOK {
		err := fmt.Errorf("status code expect 200, got %d", code)
		try.To(err)
	}
	reader := io.LimitReader(resp.Body, 50*units.MiB)
	buf := try.To1(io.ReadAll(reader))
	return c.Set(k, buf, ttlcache.DefaultTTL)
})
var fileCache = ttlcache.New[string, []byte](
	ttlcache.WithTTL[string, []byte](5*time.Minute),
	ttlcache.WithLoader[string, []byte](fileLoader),
)

func initAPI(g *echo.Group) {
	g.Any("/send_group_msg", func(c echo.Context) (err error) {
		defer err0.Then(&err, nil, nil)

		rw.Lock()
		defer rw.Unlock()

		var params SendGroupMsgParams
		try.To(c.Bind(&params))
		ctx := c.Request().Context()

		group := fmt.Sprintf("%d", params.Group)
		try.To(active(ctx, group))
		try.To(clear(ctx))

		textOnly := false
		if text := params.AutoEscape; text != nil {
			textOnly = *text
		}
		logger := slog.With(
			"qq_num", params.Group,
		)
		if textOnly {
			try.To(write(ctx, "text/plain", []byte(params.Message)))
			try.To(send(ctx))
		} else {
			items := msg.ParseString(params.Message)
			for _, m := range items {
				switch m.Type {
				case "at":
					for _, pair := range m.Data {
						switch {
						case pair.K == "qq" && pair.V == "all":
							try.To(atAll(ctx))
						}
					}
				case "text":
					for _, pair := range m.Data {
						switch {
						case pair.K == "text":
							try.To(write(ctx, "text/plain", []byte(pair.V)))
						}
					}
				case "image":
					for _, pair := range m.Data {
						switch {
						case pair.K == "file":
							img := pair.V
							item := fileCache.Get(img)
							if item == nil {
								msg := fmt.Sprintf("get img failed, link: %s", img)
								logger.Debug(msg)
								img := fmt.Sprintf("[img=%s]", img)
								try.To(write(ctx, "text/plain", []byte(img)))
								continue
							}
							try.To(write(ctx, "image/png", item.Value()))
						}
					}
				}
			}
			try.To(send(ctx))
		}
		logger.Debug("发送完成")
		return c.JSON(http.StatusOK, map[string]any{
			"message_id": 0,
		})
	})
}

func active(ctx context.Context, group string) error {
	return exec.CommandContext(ctx, "wmctrl", "-R", group).Run()
}

func clear(ctx context.Context) error {
	if err := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "ctrl+a").Run(); err != nil {
		return err
	}
	if err := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "BackSpace").Run(); err != nil {
		return err
	}
	return nil
}

func atAll(ctx context.Context) error {
	if err := exec.CommandContext(ctx, "xdotool", "type", "--delay", "100", "@").Run(); err != nil {
		return err
	}
	if err := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "Return").Run(); err != nil {
		return err
	}
	if err := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "ctrl+a").Run(); err != nil {
		return err
	}
	if err := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "ctrl+c").Run(); err != nil {
		return err
	}
	if err := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "Right").Run(); err != nil {
		return err
	}
	var buf bytes.Buffer
	clipboard := exec.CommandContext(ctx, "copyq", "clipboard")
	clipboard.Stdout = &buf
	if err := clipboard.Run(); err != nil {
		return err
	}
	if content := buf.String(); !strings.HasSuffix(content, "@全体成员 ") {
		return fmt.Errorf("@全体成员失败. 实际内容为: %s", content)
	}
	return nil
}

func write(ctx context.Context, mime string, buf []byte) error {
	copy := exec.CommandContext(ctx, "copyq", "copy", mime, "-")
	copy.Stdin = bytes.NewBuffer(buf)
	if err := copy.Run(); err != nil {
		return err
	}
	paste := exec.CommandContext(ctx, "copyq", "paste")
	if err := paste.Run(); err != nil {
		return err
	}
	return nil
}

func send(ctx context.Context) error {
	enter := exec.CommandContext(ctx, "xdotool", "getactivewindow", "key", "Return")
	return enter.Run()
}

type SendGroupMsgParams struct {
	Group      int64  `json:"group_id" form:"group_id" query:"group_id"`
	Message    string `json:"message" form:"message" query:"message"`
	AutoEscape *bool  `json:"auto_escape" form:"auto_escape" query:"auto_escape"` // 消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 message 字段是字符串时有效
}
