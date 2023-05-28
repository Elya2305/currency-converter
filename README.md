# BTC/UAH converter

### Description
This project is a simple currency converter that allows you to get current price for 1 BTC in UAH and subscribe to receive email notifications with current rate.
The project is written in scope of [Software Engineering School](https://www.genesis-for-univ.com/genesis-software-engineering-school-3?utm_source=email_campaing&utm_medium=email&utm_campaign=se3&utm_content=free) 
### Endpoints
GET `/rate` - returns current price for 1 BTC in UAH

POST `/subscribe` - subscribes user to receive email notifications that are sent using /sendEmails

POST `/sendEmails` - sends emails to all subscribed users with current price for 1 BTC in UAH

### Pre-requisites
You'll need to provide the following environment variables:
* EMAIL_FROM - email from which emails will be sent
* EMAIL_PASSWORD - password for email from which emails will be sent. To generate a password for gmail, follow this [link](https://support.google.com/accounts/answer/185833?hl=en)

### How to run locally

To run the project locally:

* docker build -t currency-converter .
* docker run \
  -e EMAIL_FROM='\<email from>' \
  -e EMAIL_PASSWORD='\<password>' \
  -p 9090:9090 currency-converter

### Provider

https://docs.coinapi.io/market-data/rest-api/exchange-rates
(100 requests / day)
