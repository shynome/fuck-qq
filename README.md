## 简介

使用 gui automation 工具实现将动态推送. 强阻塞, 一次只能运行一个任务

## 依赖

```sh
apt install xclip wmctrl xdotool
```

- `wmctrl` 找到对应的 QQ 独立窗口
- `xdotool` 操作对应窗口
- `copyq` 输入文本

## 注意

- 打开独立窗口后需要输入@一次, 避免首次@全体时选不中@全体
