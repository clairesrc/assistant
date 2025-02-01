# assistant

## what
heads-up display/dashboard using local Ollama instance for AI-assisted summaries of personal data.

web-based display interface.

this assumes you already have ollama running locally.

some features might include:
- reads your unread emails and summarizes them, flags any that might require you to respond, provides a button to mark-as-read the rest
- presents daily information about your calendar events, weather in a friendly conversational tone
- provides a summary of today's news

once the display feature is working, the plan is to add a voice interface using Whisper STT and Coqui TTS.

none of this is working yet, i am still working on the basic scaffolding and display.

## run
copy .env.example to .env and fill in the values

```bash
docker compose up
```
