# Chat

This project contains a online chat room that allow registered users to send text messages between each other and also to use a StockBot to obtain the current stock amount from [stooq.com](https://stooq.com).

The goal of this exercise is to create a simple browser-based chat application using Go.
This application should allow several users to talk in a chatroom and also to get stock quotes
from an API using a specific command.

## How to run
The project is dockerized so, use the following command to run the application and dependencies:
```
docker-compose up
```

Wait until the message from `chatroom-app-1` appears:
```
$ chatroom-app-1       | [GIN-debug] Listening and serving HTTP on :8080
```

Now, you can access the chat: http://localhost:8080/sign

## Accomplished Features:
- Allow registered users to log in and talk with other users in a chatroom;
- Allow users to post messages as commands into the chatroom with the following format: `/stock=stock_code` (ex: `/stock=aapl.us`);
- Create a decoupled bot that will call **stooq.com API** passing `stock_code` as a parameter. https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv where **_aapl.us_** is the stock_code;
- The bot should parse the received CSV file. The message will be a stock quote using the following format: **_“APPL.US quote is $93.42 per share”_**. The post owner will be the bot;
- Have the chat messages ordered by their timestamps.

## Missing Features:
- The bot should send the message back into the chatroom using a message broker like RabbitMQ;
- Show only the last 50 messages on chat.
- Unit test the functionality you prefer.

## Considerations

The project was a great study case and it was pretty enjoyable to code it.
It still not completed, considering that is missing tests and documentation (things that are mandatory to any project), but also has some features that makes this project easily refactorable to become a actual production project.

For who'll read this project, thank you for your time and attention.
