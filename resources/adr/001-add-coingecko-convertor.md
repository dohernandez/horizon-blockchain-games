# Add CoinGecko convertor

* Status: accepted
* Deciders: Darien Hernandez
* Date: 2024-11-07

## Context and Problem Statement

To convert the token currencies into USD, there is a need to extract the conversion data using the [CoinGecko API](https://www.coingecko.com/en/api).

## Considered Options

* Given the CoinGecko API does not provide easily the usd rate based on the given currency symbol.
* Given it is not sure the amount of request per second the service will have in the future.
* Given the pipeline will be run on a daily basis.

## Decision Outcome

* Hardcode the mapping between the currency symbol and coin id in the CoinGecko API, to avoid the need to make a request to the API to `GET coins/list?status=active` endpoint first to get the coin id and then find the correct coin id making request to the API to `GET coins/<coin_id>` endpoint to match the currency address because more than one coin can have the same symbol.
* Use in-memory cache to avoid making the same request to the CoinGecko API to get usd price for the same currency symbol.

### Positive Consequences

* Simplify the logic to convert to usd.
* Avoid making unnecessary requests to the CoinGecko API.

### Negative Consequences

* The mapping between the currency symbol and coin id in the CoinGecko API can be outdated missing new coins or having coins that are not active anymore.

## Links

- [https://docs.coingecko.com/reference/simple-price](https://docs.coingecko.com/reference/simple-price)
- [https://docs.coingecko.com/reference/coins-list](https://docs.coingecko.com/reference/coins-list)
- [https://docs.coingecko.com/reference/coins-id](https://docs.coingecko.com/reference/coins-id)
