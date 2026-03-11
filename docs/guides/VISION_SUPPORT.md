# Vision Support Guide

Send images to AI models for analysis and understanding.

## Overview

PicoClaw supports multi-modal AI with vision capabilities, allowing you to send images to compatible models.

## Supported Models

### OpenAI
- GPT-4V (gpt-4-vision-preview)
- GPT-4o (gpt-4o)
- GPT-4o-mini (gpt-4o-mini)

### Anthropic
- Claude 3 Opus (claude-3-opus)
- Claude 3 Sonnet (claude-3-sonnet)
- Claude 3 Haiku (claude-3-haiku)
- Claude 3.5 Sonnet (claude-3-5-sonnet)

## Usage

### Direct File Path
```bash
picoclaw agent
> Describe this image: /path/to/image.jpg
> What's in this chart: ~/Downloads/chart.png
```

### Media Reference
```bash
> Analyze media://image-id
```

### Multiple Images
```bash
> Compare these images: /path/to/image1.jpg and /path/to/image2.jpg
```

## Configuration

### Max File Size
```json
{
  "agents": {
    "defaults": {
      "max_media_size": 20971520
    }
  }
}
```

Default: 20MB (20971520 bytes)

### Environment Variable
```bash
PICOCLAW_AGENTS_DEFAULTS_MAX_MEDIA_SIZE=20971520
```

## Supported Formats

### Image Types
- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)

### MIME Types
- image/jpeg
- image/png
- image/gif
- image/webp

## How It Works

### 1. File Detection
PicoClaw detects image paths in messages.

### 2. Encoding
Images are encoded to base64 using streaming (memory efficient).

### 3. Multipart Content
Message sent with both text and image:
```json
{
  "role": "user",
  "content": [
    {"type": "text", "text": "Describe this image"},
    {"type": "image", "source": {"type": "base64", "data": "..."}}
  ]
}
```

### 4. Response
Model analyzes image and responds with description.

## Examples

### Image Description
```bash
> Describe /path/to/photo.jpg
```

Response:
```
The image shows a sunset over mountains with orange and purple colors...
```

### Chart Analysis
```bash
> Analyze this chart: /path/to/sales-chart.png
```

Response:
```
This bar chart shows sales data for Q1-Q4. Key insights:
- Q4 had highest sales at $500K
- Q1 was lowest at $200K
- Upward trend throughout the year
```

### Code Screenshot
```bash
> What's wrong with this code: /path/to/error-screenshot.png
```

Response:
```
The code has a syntax error on line 5:
- Missing closing parenthesis
- Should be: result = calculate(x, y)
```

## Performance

### Memory Efficiency
- Streaming base64 encoding
- No full file load into memory
- Efficient for large images

### Size Limits
- Default: 20MB
- Configurable per agent
- Automatic rejection if too large

## Error Handling

### File Not Found
```
Error: Image file not found: /path/to/image.jpg
```

### File Too Large
```
Error: Image exceeds max size (20MB): /path/to/large.jpg
```

### Unsupported Format
```
Error: Unsupported image format: .bmp
```

### Model Not Compatible
```
Error: Model does not support vision: gpt-3.5-turbo
```

## Best Practices

### Image Quality
- Use clear, high-resolution images
- Avoid blurry or low-quality images
- Crop to relevant area

### File Size
- Compress large images
- Use appropriate resolution
- Balance quality vs size

### Prompts
- Be specific about what to analyze
- Ask clear questions
- Provide context if needed

## Troubleshooting

### Image Not Recognized
Check:
- File path is correct
- File format is supported
- File size is within limit

### Poor Analysis
Try:
- Higher resolution image
- Better lighting/contrast
- More specific prompt

### Memory Issues
Reduce:
- Image file size
- Max media size setting
- Number of concurrent requests

## API Reference

### Media Field
```go
type Message struct {
    Role    string
    Content string
    Media   []MediaRef
}

type MediaRef struct {
    Type     string // "image"
    MimeType string // "image/jpeg"
    Data     string // base64 data
}
```

### Resolution
```go
resolveMediaRefs(message) -> Message
// Converts file paths to base64 data
```

## See Also

- [Model Configuration](../reference/MODEL_CONFIGURATION.md)
- [v0.2.1 Features](V0.2.1_FEATURES.md)
- [Configuration Guide](../reference/CONFIGURATION.md)

---

**Version**: v0.2.1  
**Last Updated**: 2026-03-09
