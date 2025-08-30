# Local Logging Setup

Prosty system logowania lokalnego dla Astras API używający Dozzle.

## Przegląd

System logowania zapewnia:
- ✅ **Strukturalny format JSON** dla wszystkich logów
- 🐳 **Dozzle** - prosty viewer logów kontenerów Docker
- 📁 **Przeglądarka plików** - web interface dla plików logów
- 🔄 **Real-time** wyświetlanie logów

## Szybki start

```bash
# Uruchom Dozzle
./scripts/local-logging.sh start

# Otwórz dashboard
./scripts/local-logging.sh dashboard

# Zatrzymaj serwisy
./scripts/local-logging.sh stop
```

## Dostępne interfejsy

Po uruchomieniu `./scripts/local-logging.sh start`:

- **🐳 Dozzle**: http://localhost:8080 - logi kontenerów Docker w real-time
- **📁 File viewer**: http://localhost:8081 - przeglądarka plików logów z interaktywnym JSON viewer

## Komendy

```bash
# Uruchom logging
./scripts/local-logging.sh start

# Wyświetl logi w terminalu  
./scripts/local-logging.sh logs

# Otwórz web dashboard
./scripts/local-logging.sh dashboard

# Sprawdź status serwisów
./scripts/local-logging.sh status

# Zatrzymaj serwisy
./scripts/local-logging.sh stop

# Usuń dane logów
./scripts/local-logging.sh clean
```

## Format logów

Wszystkie logi używają strukturalnego formatu JSON:

```json
{
  "@timestamp": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "service": "kid-service",
  "message": "Processing request", 
  "request_id": "req_123456",
  "http_method": "GET",
  "http_path": "/kids",
  "status_code": 200,
  "duration": 45,
  "environment": "local"
}
```

## Integracja z SAM Local

System automatycznie wykrywa środowisko SAM Local:

```bash
# 1. Uruchom logging
./scripts/local-logging.sh start

# 2. Uruchom SAM z siecią logging
sam local start-api --env-vars env.json --docker-network astras-logging --port 3000

# 3. Zobacz logi w czasie rzeczywistym
./scripts/local-logging.sh dashboard
```

## Workflow developmentu

```bash
# 1. Start logging
./scripts/local-logging.sh start

# 2. Start database
docker-compose up -d

# 3. Build services
mage build:all

# 4. Start SAM Local
sam local start-api --env-vars env.json --docker-network astras-logging --port 3000

# 5. Test endpoints
curl http://localhost:3000/kids

# 6. View logs
./scripts/local-logging.sh dashboard
```

## Architektura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SAM Local     │    │   Log Files     │    │   Dozzle        │
│   Lambda        │───▶│   /logs/*.log   │───▶│   :8080         │
│   Functions     │    │   (JSON)        │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                               │   File Viewer   │
                                               │   :8081         │
                                               └─────────────────┘
```

## Rozwiązywanie problemów

### Brak logów

1. Sprawdź status serwisów:
   ```bash
   ./scripts/local-logging.sh status
   ```

2. Sprawdź czy SAM używa poprawnej sieci:
   ```bash
   sam local start-api --docker-network astras-logging
   ```

3. Sprawdź uprawnienia plików:
   ```bash
   ls -la logs/
   ```

### Konflikty portów

```bash
# Zatrzymaj logging
./scripts/local-logging.sh stop

# Lub zatrzymaj wszystkie kontenery
docker stop $(docker ps -q)
```

## Lokalizacja plików

- **Logi**: `./logs/astras-*.log`
- **Skrypt**: `./scripts/local-logging.sh`
- **Konfiguracja**: `./config/docker-compose.dozzle.yml`
- **SAM env**: `./env.json`