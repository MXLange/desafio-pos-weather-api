# Clima por CEP

API em Go que recebe um CEP de 8 digitos, consulta a cidade na ViaCEP e retorna a temperatura atual pela WeatherAPI em Celsius, Fahrenheit e Kelvin.

## URL Cloud Run

Substitua pelo endereco gerado apos o deploy no Google Cloud Run:

```text
TODO: https://SEU-SERVICO-REGIAO.run.app
```

## Deploy no Cloud Run

```sh
gcloud run deploy desafio-pos-weather-api \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=sua-chave-weatherapi,WEATHER_BASE_URL=https://api.weatherapi.com,API_CEP_BASE_URL=https://viacep.com.br
```

Depois do deploy, copie a URL exibida pelo comando e atualize a secao `URL Cloud Run`.

## Variaveis de ambiente

```env
API_PORT=8080
WEATHER_API_KEY=sua-chave-weatherapi
WEATHER_BASE_URL=https://api.weatherapi.com
API_CEP_BASE_URL=https://viacep.com.br
```

No Cloud Run, a aplicacao tambem aceita a variavel `PORT` injetada pela plataforma quando `API_PORT` nao estiver definida.

## Rodar localmente

```sh
cp example.env .env
go run ./cmd/server
```

Endpoint:

```sh
curl "http://localhost:8080/temperature?cep=01001000"
```

Resposta de sucesso:

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

## Rodar com Docker

```sh
docker build -t weather-api .
docker run --rm -p 8080:8080 --env-file .env weather-api
```

## Testes

```sh
go test ./...
```

## Contrato de erros

| Condicao | Status | Resposta |
| --- | --- | --- |
| CEP sem 8 digitos ou com caracteres invalidos | 422 | `invalid zipcode` |
| CEP nao encontrado | 404 | `can not find zipcode` |
