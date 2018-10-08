# api-server

# run
```docker-compose -f docker-compose.yml up --build```
# api
- ```/currencies/currencies/getList```: Get list tokens
- ```/currencies/get_ethrate_buy/```: Get rate buy token base on ETH
```EX: http://localhost:8000/currencies/get_ethrate_buy?&id=0xdd974D5C2e2928deA5F71b9825b8b646686BD200&id=0xd26114cd6EE289AccF82350c8d8487fedB8A0C07&qty=123-123-123&qty=234-234```