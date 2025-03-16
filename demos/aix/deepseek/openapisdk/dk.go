package openaisdk

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Chat() {
	client := openai.NewClient(
		option.WithAPIKey("sk-6b7ce843c9ef4187a7c8f42928841eba"), // defaults to os.LookupEnv("OPENAI_API_KEY")
		option.WithBaseURL("https://api.deepseek.com"),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
		}),
		Model: openai.F("deepseek-chat"),
	})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
