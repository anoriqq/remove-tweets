# 🐦 remove-tweets

過去ツイとリツイートを消すやつ

## Usage

```sh
go install github.com/anoriqq/remove-tweets/cmd/rmtw@latest

# screenname: Target username for removing tweets
# maxid:      Remove tweets before the tweet of this tweet ID
# untile:     remove tweets up to this date time
rmtw -screenname anoriqq -maxid 1252213125127917568 -until 2020-04-01T00:00:00Z
```

