# ADK Weather Report Agent

一个基于 Google ADK (Agent Development Kit) 框架开发的智能天气报告代理，能够查询城市天气信息并分析用户对天气的情感反馈。

## 项目简介

这是一个使用 Go 语言开发的 AI 代理应用，集成了以下功能：

- **天气查询**：通过 OpenWeatherMap API 获取指定城市的实时天气信息
- **情感分析**：分析用户对天气反馈的情感倾向（积极/消极/中性）
- **自然语言交互**：使用 Grok (x.ai) 大语言模型进行智能对话

## 主要功能

### 1. 天气报告查询 (`get_weather_report`)
- 根据用户提供的城市名称查询实时天气
- 返回城市名称、天气描述和温度信息
- 支持错误处理和友好的错误提示

### 2. 情感分析 (`analyze_sentiment`)
- 分析用户对天气反馈的文本情感
- 返回情感类型（positive/negative/neutral）和置信度
- 基于关键词匹配进行情感判断

### 3. 智能对话
- 使用 Grok-4-1-fast 模型进行自然语言理解
- 自动识别用户意图并调用相应的工具
- 提供流畅的对话体验

## 技术栈

- **语言**: Go 1.25.0
- **框架**: Google ADK (Agent Development Kit) v0.2.0
- **LLM**: Grok (x.ai) via OpenAI-compatible API
- **天气API**: OpenWeatherMap
- **依赖管理**: Go Modules

## 环境要求

- Go 1.25.0 或更高版本
- 有效的 API 密钥：
  - `XAI_API_KEY`: x.ai (Grok) API 密钥
  - `OWM_API_KEY`: OpenWeatherMap API 密钥

## 安装步骤

### 1. 克隆项目

```bash
git clone <repository-url>
cd adk-weatherreport
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

创建 `.env` 文件或设置环境变量：

```bash
export XAI_API_KEY="your-xai-api-key"
export OWM_API_KEY="your-openweathermap-api-key"
```

或者使用 `.env` 文件（需要安装支持工具）：

```bash
# .env
XAI_API_KEY=your-xai-api-key
OWM_API_KEY=your-openweathermap-api-key
```

### 4. 获取 API 密钥

#### XAI API 密钥 (Grok)
1. 访问 [x.ai](https://x.ai)
2. 注册账号并获取 API 密钥

#### OpenWeatherMap API 密钥
1. 访问 [OpenWeatherMap](https://openweathermap.org/api)
2. 注册免费账号
3. 在控制台获取 API Key

## 运行项目

### 基本运行

```bash
go run main.go
```

### 使用 ADK Launcher

项目使用 Google ADK 的 launcher 框架，支持多种运行模式：

```bash
# 查看帮助信息
go run main.go --help

# 运行代理（具体命令取决于 ADK launcher 配置）
go run main.go [launcher-options]
```

## 项目结构

```
adk-weatherreport/
├── main.go                 # 程序入口，初始化代理和启动器
├── weatherreport.go        # 核心业务逻辑：天气查询和情感分析
├── llm/                    # LLM 适配器
│   ├── grok.go            # Grok (x.ai) 模型适配器
│   ├── openai.go          # OpenAI 模型适配器
│   ├── openrouter.go      # OpenRouter 模型适配器
│   └── request2params.go  # 请求参数转换工具
├── go.mod                  # Go 模块依赖
├── go.sum                  # 依赖校验和
└── README.md              # 项目文档
```

## 核心组件说明

### WeatherSentimentAgent

主要的代理实现，包含：
- **模型配置**: 使用 Grok-4-1-fast 模型
- **工具集成**: 
  - `get_weather_report`: 天气查询工具
  - `analyze_sentiment`: 情感分析工具
- **指令系统**: 智能识别用户意图并调用相应工具

### 工具函数

#### `getWeatherReport`
- **功能**: 查询指定城市的天气信息
- **参数**: `city` (城市名称)
- **返回**: 天气报告（状态、报告内容）

#### `analyzeSentiment`
- **功能**: 分析文本情感
- **参数**: `text` (待分析文本)
- **返回**: 情感类型和置信度

## 使用示例

启动代理后，可以通过对话方式使用：

```
用户: "What's the weather in Beijing?"
代理: [调用 get_weather_report 工具]
代理: "The weather in Beijing is clear sky with a temperature of 15.5 degrees Celsius."

用户: "That's good!"
代理: [调用 analyze_sentiment 工具]
代理: "I'm glad you're happy with the weather!"
```

## 开发说明

### 修改 LLM 模型

在 `weatherreport.go` 的 `NewWeatherSentimentAgent` 函数中，可以切换不同的 LLM 模型：

```go
// 使用 Grok
model, err := llm.NewGrokModel(ctx, "grok-4-1-fast", &genai.ClientConfig{
    APIKey: os.Getenv("XAI_API_KEY"),
})

// 或使用 OpenAI
model, err := llm.NewOpenAIModel(ctx, "gpt-4", &genai.ClientConfig{
    APIKey: os.Getenv("OPENAI_API_KEY"),
})

// 或使用 OpenRouter
model, err := llm.NewOpenRouterModel(ctx, "anthropic/claude-3-opus", &genai.ClientConfig{
    APIKey: os.Getenv("OPENROUTER_API_KEY"),
})
```

### 扩展功能

可以添加更多工具来扩展代理功能：

1. 在 `weatherreport.go` 中定义新的工具函数
2. 使用 `functiontool.New` 创建工具实例
3. 将工具添加到 `llmagent.Config` 的 `Tools` 数组中

## 故障排除

### 常见问题

1. **API 密钥错误**
   - 确保环境变量已正确设置
   - 检查 API 密钥是否有效

2. **天气查询失败**
   - 验证 OpenWeatherMap API 密钥
   - 检查城市名称是否正确
   - 确认网络连接正常

3. **LLM 调用失败**
   - 检查 XAI API 密钥
   - 确认账户有足够的配额
   - 查看日志中的详细错误信息

## 许可证

查看 [LICENSE](LICENSE) 文件了解许可证信息。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 相关链接

- [Google ADK 文档](https://github.com/google/adk)
- [OpenWeatherMap API](https://openweathermap.org/api)
- [x.ai (Grok)](https://x.ai)