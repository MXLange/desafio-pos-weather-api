# Clima por CEP

API em Go que recebe um CEP de 8 digitos, consulta a cidade na ViaCEP e retorna a temperatura atual pela WeatherAPI em Celsius, Fahrenheit e Kelvin.

## URL Cloud Run

Exemplo de chamada:

```text
https://temperature-api-9654929303.europe-west1.run.app/temperature?cep=01001000
```

## Variaveis de ambiente

```env
API_PORT=8080
WEATHER_API_KEY=sua-chave-weatherapi
WEATHER_BASE_URL=https://api.weatherapi.com
API_CEP_BASE_URL=https://viacep.com.br
```

## Rodar localmente

Faça o export das envs do example.env no terminal e coloque sua api key

```sh
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
