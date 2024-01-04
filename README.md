## 简介

使用 gui automation 工具实现将动态推送. 强阻塞, 一次只能运行一个任务

## 运行

```sh
docker run \
  --name fuck-qq \
  --restart always \
  -d \
  --shm-size=512m \
  -p 6901:6901 \
  -p 5700:5700 \
 -e VNC_PW=password \
 shynome/fuck-qq:v0.0.6
```

## 测试是否正常运行

```sh
curl 'http://127.0.0.1:5700/onebot/send_group_msg?group_id=qq群号&message=hello'
```

## 依赖

```sh
apt install wmctrl xdotool copyq
```

- `wmctrl` 找到对应的 QQ 独立窗口
- `xdotool` 操作对应窗口
- `copyq` 输入文本

## 注意

- 打开独立窗口后需要输入@一次, 避免首次@全体时选不中@全体
- 内存占用: 600M QQ 基础, 每个独立窗口+100M, 一个用户 800M 成本

## 如何定价

每次用户发送请求 1op/6s, 能忍受的最大延时为 30s, 所以一台机器最多容纳 5 名用户,
800+(5x2)\*100 = 2G, 2G 主机费用 5 人分担, 使用 NAT 型主机费用为 240/年, 所以分摊到
每人的成本为 48/年, 40% 出售率要收回成本, 所以定价为 0.4\*x=48, 即 x = 48/0.4= 120 元/年
