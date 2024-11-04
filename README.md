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

### Migrating Data 
On initial boot, the system will detect if your Plex library is stored in the vector
database. If it is not, your media will be retreived. You need to provide the default
library to download media from via the `PLEX_DEFAULT_LIBRARY_SECTION` environment 
variable. You will need to query Plex yourself to get this, but for me, my movies are
in section 3. I will fall back to this section if you do not provide one.

## Connecting to your LLM
This recommendation engine connects to Ollama. You can bring your own or 
run it on the cloud. Just provide `OLLAMA_ADDRESS`, `OLLAMA_EMBEDDING_MODEL`, and `OLLAMA_LANGUAGE_MODEL` as 
environment variables. 

### Grounding your LLM
You can find the RAG prompt in `backend/internal/pkg/langchain/generate.go`. This is 
written to my specific needs. If your needs are not my needs, adjust the 
grounding prompt accordingly. Note that if you have a very large media collection
that it may exceed the context window of the model you are using. You can 
adjust the amount of titles you retreive by adjusting the limits passed into
the media getters in `backend/internal/pkg/plex/api.go`. 

## Building and Running
### Compiling from source
Download this repository and build the app using 
`go build -o recommendations .backend/cmd/main.go`. This builds a binary called
`recommendations` that you can run with `./recommendations`. 

### Cached Responses
A Postgres database is attached at `./pg-data` and is used to cache recommendations. 
The titles provided from recently viewed get base-64 encoded and stored along
with the normalized recomendation provided by the LLM. This cache is looked up 
before requesting LLM recommendation and is returned as is if one is found.

#### Default Postgres Environment
The default Postgres environment is as follows:
```
POSTGRES_PASSWORD=postgres
POSTGRES_USER=postgres
POSTGRES_DB=caches
POSTGRES_PORT=5432
POSTGRES_HOST=postgres
```
You may override this by providing your own values to these environment
variables


### Docker Compose
There is a `docker-compose.yml` file in the root of this repo. You can run
the Recommendation Engine and the associated Weaviate database by issuing the
`TAG_VERSION=foo docker compose up --build` command. Make sure you have 
Docker and Docker Compose installed. 

### Testing this app
There are several tests in the internal packages that were almost all written by 
an LLM. You can test this program using `go test ./...` from the root of this repo. These tests are automatically run when you build with Docker.

### Open Telemetry 
[Open Telemetry](https://opentelemetry.io/docs/what-is-opentelemetry/) tracing is instrumented in the backend. To use this out of the
box, set `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=otlp://jaeger:4317` in the environment and ensure that the Jaeger service
is started from Docker Compose. Traces are collected in Jaeger and accessible at `http://localhost:16686`. Metrics are 
available to extend by passing the `telemetry.WithMeter` option `telemetry.InitOtel()`. You will have to instrument 
metrics yourself. 

#### Disabling Telemetry
Telemetry can be disabled by setting `DISABLE_TELEMETRY=true` in your environment. Note that if using `docker compose up`, the 
Jaeger container will still start. 