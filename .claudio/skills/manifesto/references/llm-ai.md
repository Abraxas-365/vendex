# LLM, Embeddings, Vector Store, Agents & Memory

## LLM Client

```go
import (
    "github.com/Abraxas-365/manifesto/internal/ai/llm"
    "github.com/Abraxas-365/manifesto/internal/ai/providers/aiopenai"
)

provider := aiopenai.NewOpenAIProvider(apiKey)
client := llm.NewClient(provider)
```

### Chat
```go
messages := []llm.Message{
    llm.SystemMessage("You are a helpful assistant."),
    llm.UserMessage("Hello!"),
}
resp, err := client.Chat(ctx, messages, llm.WithModel("gpt-4.1"))
```

### Streaming
```go
stream, err := client.ChatStream(ctx, messages)
for token := range stream {
    fmt.Print(token)
}
```

### Multimodal messages
```go
msg := llm.UserMessageWithImages("Describe this image", imageURL)
msg := llm.UserMessageWithFiles("Summarize this doc", fileData)
```

### Tools / function calling
```go
tool := llm.Tool{
    Name:        "get_weather",
    Description: "Get current weather for a location",
    Parameters:  weatherParamsSchema, // JSON schema
}
resp, err := client.Chat(ctx, messages, llm.WithTools([]llm.Tool{tool}))
// Inspect resp.ToolCalls to handle invocations
```

---

## Providers

All providers implement `llm.Provider`. Swap without changing business logic:

```go
aiopenai.NewOpenAIProvider(apiKey)         // OpenAI (default: gpt-4.1)
aianthropix.NewAnthropicProvider(apiKey)   // Anthropic Claude
aigemini.NewGeminiProvider(apiKey)         // Google Gemini
aiaws.NewBedrockProvider(cfg)              // AWS Bedrock
```

Import paths: `internal/ai/providers/<provider-name>/`

---

## Embeddings

```go
import "github.com/Abraxas-365/manifesto/internal/ai/embedding"

// Single
vec, err := embedder.Embed(ctx, "text to embed")

// Batch
vecs, err := embedder.EmbedBatch(ctx, []string{"doc1", "doc2"})
```

---

## Vector Store

```go
import "github.com/Abraxas-365/manifesto/internal/ai/vstore"

// Upsert
err := vsClient.Upsert(ctx, vectors, vstore.WithNamespace("docs"))

// Query
results, err := vsClient.Query(ctx, queryVec,
    vstore.WithTopK(10),
    vstore.WithNamespace("docs"),
    vstore.WithFilter(map[string]any{"tenant": tenantID}),
)

// Delete
err := vsClient.Delete(ctx, ids, vstore.WithNamespace("docs"))
```

Supports filtering, namespacing, batch ops, and hybrid search via options.

---

## Agent Loop (agentx)

```go
import "github.com/Abraxas-365/manifesto/internal/ai/llm/agentx"

agent := agentx.New(client, tools, agentx.WithMaxIterations(10))
result, err := agent.Run(ctx, userMessage)
```

---

## Memory (memoryx)

Three memory types available:

```go
import "github.com/Abraxas-365/manifesto/internal/ai/llm/memoryx"

mem := memoryx.NewConversationMemory()     // full history
mem := memoryx.NewContextualMemory(k)      // last k turns
mem := memoryx.NewSummarizingMemory(llm)   // auto-summarizes older turns
```

Use memory with a client by injecting it before chat calls:
```go
messages = mem.Load(ctx)
messages = append(messages, llm.UserMessage(input))
resp, err := client.Chat(ctx, messages)
mem.Save(ctx, input, resp.Content)
```

---

## Harness Integration (hada-commerce specific)

hada-commerce uses `github.com/Abraxas-365/harness` for the agent tool loop, not agentx directly.
See `.claudio/CLAUDE.md` → "Agent Integration Layer" for the tool implementation pattern.

Tools live in `internal/agent/`. Each wraps a domain service and implements `tools.Tool`.
