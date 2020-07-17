curl -X POST http://checkout/proxy -d @checkout.json -i
curl -X POST http://frontend/proxy -d @frontend.json -i
curl -X POST http://currency/proxy -d @currency.json -i
curl -X POST http://email/proxy -d @email.json -i
curl -X POST http://payment/proxy -d @payment.json -i