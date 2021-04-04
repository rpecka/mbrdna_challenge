# MBRDNA Chatbot

A chatbot server and client for remotely accessing a Mercedes-Benz connected vehicle.

![Screen Shot 2021-04-04 at 7 45 30 PM](https://user-images.githubusercontent.com/13580126/113524693-53d7c680-957e-11eb-88d1-4e1a228e37d5.png)

# Running
```
git clone git@github.com:rpecka/mbrdna_challenge.git
cd mbrdna_challenge
```
Run the server:
```
go run ./cmd/server/...
```
Run the client:
```
go run ./cmd/client/...
```

### Credentials
I've hardcoded my own Houndify credentials into the project so that you can run with my rules and intent configurations.

I also hardcoded my Mercedes-Benz Connected Vehicle (MBCV) credentials in case you don't want to use your own, but you will be prompted to provide your own on startup if you wish.

### OAuth
Since I chose the command line option, the process for getting an OAuth token is the following:
1. After you have chosen the credentials you would like to use on the command line, open the printed URL in your browser and authorize the app.
2. You will then be redirected to `localhost/...`. Copy the URL you were redirected to and paste it onto the command line.
3. The client will use the authorization code form the URL to get an auth token which will be used for the remainder of the session.
