# P.A.N.D.A.

## Protective Automaton for Nefarious Discourse and Antagonsim

![P.A.N.D.A.](./panda-bot.png)

### Getting Started

Running P.A.N.D.A. requires [Node.js](https://nodejs.org/)

Clone the project to your computer, and enter the project folder:
```
git clone https://github.com/Sephory/panda-bot.git
cd panda-bot
```

Install the required dependencies by running:

```
npm install
```

Create a file named `.env` in the project folder, open it and add the following text:
```
PANDA_USERNAME=<username>
PANDA_OATH_TOKEN=<token>
PANDA_CHANNELS=<channels>
```

Replace the placeholders with your own configuration:

- Replace `<username>` with the Twitch username you would like the bot to log in with
- Replace `<token>` with an authentication token for your bot, which can be acquired at [https://twitchapps.com/tmi/](https://twitchapps.com/tmi/)
- Replace `<channels>` with one or more Twitch channels you would like the bot to connect to.  For multiple channels, separate them with a `;`

After saving your configuration, run the following command to start P.A.N.D.A. and have it connect to your chosen servers:
```
npm run dev
```
