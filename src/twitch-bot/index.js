const tmi = require('tmi.js')

const clientConfig = require('./config.js')

const commands = require('./commands/index.js')

const client = tmi.Client(clientConfig)

client.connect()

client.on('message', (channel, tags, message, self) => {
  if (self || !message.startsWith('!')) return

  const args = message.slice(1).split(' ')
  const commandName = args.shift().toLowerCase()

  const command = commands[commandName]

  if (command) {
    command(client, channel, tags, args)
  } else {
    client.say(channel, `${tags['display-name']}, that's not a valid command, dummy!`)
  }
})
