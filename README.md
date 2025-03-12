# OpenAI Computer Use Example using Go

This is a simple example of utilizing OpenAI's computing capabilities in Go, leveraging go-rod to control a browser.

https://platform.openai.com/docs/guides/tools-computer-use

https://github.com/openai/openai-cua-sample-app

https://github.com/go-rod/rod


## Usage

```bash
export OPENAI_API_KEY="your-api-key-here"
go run ./example
```


### with flag
```bash
go run ./example -url "https://duckduckgo.com/" -prompt "Find out the winner of the Academy Award for Best Picture in 2025 and tell me the title." -timeout "3m"
```


## License

MIT License
