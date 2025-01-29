# Weather Cep

## Descrição

Weather Cep é um serviço em Go que permite a consulta de temperaturas usando CEPs brasileiros. Dado um CEP válido de 8 dígitos, o sistema identifica a cidade associada e retorna a temperatura atual em graus Celsius, Fahrenheit e Kelvin.

O serviço é implementado em Go e publicado no Google Cloud Run, garantindo escalabilidade e fácil acesso.

## Funcionalidades

- Recebe um CEP de 8 dígitos.
- Valida o CEP.
- Pesquisa a localização do CEP utilizando a API [viaCep](https://viacep.com.br/).
- Consulta as temperaturas da localização usando a API [WeatherAPI](https://www.weatherapi.com/).
- Converte e formata as temperaturas em Celsius, Fahrenheit e Kelvin.
- Retorna as temperaturas em um formato JSON amigável.

## Requisitos

- CEP válido com 8 dígitos.
- Código HTTP 200 em caso de sucesso.
- Código HTTP 422 e mensagem "invalid zipcode" se o CEP tiver formato incorreto.
- Código HTTP 404 e mensagem "can not find zipcode" se o CEP não for encontrado.
- Docker e Docker Compose para testes locais.

## Documentação da API

### Endpoint

`GET /weather?postalCode={postalCode}`

### Parâmetros

| Parâmetro  | Tipo   | Descrição                |
|------------|--------|--------------------------|
| postalCode | string | CEP válido com 8 dígitos |

### Exemplo de Resposta

#### Sucesso (200)

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### CEP Inválido (422)

```text
invalid zipcode
```

#### CEP Não Encontrado (404)

```text
can not find zipcode
```

## Acesso ao Serviço

Você pode acessar o sistema em: [https://go-expert-cep-475279029123.us-central1.run.app/weather?postalCode=14500000](https://go-expert-cep-475279029123.us-central1.run.app/weather?postalCode=14500000) 