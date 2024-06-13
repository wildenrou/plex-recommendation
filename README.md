# Plex Recommendations
Feed your Plex recently viewed movies into a RAG powered LLM chain to get recommendations based on your watch history.

## Connecing to your Plex
You'll need two things:
1. The IP address that your Plex runs on
2. Your Plex key
You can find the Plex key in the `preferences.xml` file in your installation path.
I run Plex on a Synology NAS using Package Center, so I found mine in 
`/PlexMediaServer/AppData/Plex Media Server/Preferences.xml`. 

Provide both of these to the `PLEX_ADDRESS` and `PLEX_TOKEN` environment variables.

## Connecting to your LLM
This recommendation engine connects to Ollama. You can bring your own or 
run it on the cloud. Just provide `OLLAMA_ADDRESS`, `OLLAMA_EMBEDDING_MODEL`, and `OLLAMA_LANGUAGE_MODEL` as 
environment variables. 

### Grounding your LLM
You can find the RAG prompt in `internal/pkg/langchain/generate.go`. This is 
written to my specific needs. If your needs are not my needs, adjust the 
grounding prompt accordingly. Note that if you have a very large media collection
that it may exceed the context window of the model you are using. You can 
adjust the amount of titles you retreive by adjusting the limits passed into
the media getters in `internal/pkg/plex/api.go`. 

## Building and Running
### Compiling from source
Download this repository and build the app using 
`go build -o recommendations ./cmd/main.go`. This builds a binary called
`recommendations` that you can run with `./recommendations`. 

### Docker Compose
There is a `docker-compose.yml` file in the root of this repo. You can run
the Recommendation Engine and the associated Weaviate database by issuing the
`TAG_VERSION=foo docker compose up --build` command. Make sure you have 
Docker and Docker Compose installed. 

### Testing this app
There are several tests in the internal packages that were almost all written by 
an LLM. You can test this program using `go test ./...` from the root of this repo. These tests are automatically run when you build with Docker.