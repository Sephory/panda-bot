const tmi = require('tmi.js')

const clientConfig = require('./config.js')

const client = tmi.Client(clientConfig)
client.connect()

client.on('message', (channel, tags, message, self) => {
  if (self || !message.startsWith('!')) return

  if (message.toLowerCase() == '!hello')
    client.say(channel, `Hi ${tags['display-name']}!`)

  const args = message.slice(1).split(' ')
  const command = args.shift().toLowerCase()
})
