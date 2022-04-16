require('dotenv').config()

const channels = process.env.PANDA_CHANNELS.split(';')

module.exports = {
  options: { debug: process.env.NODE_ENV === 'development' },
  identity: {
    username: process.env.PANDA_USERNAME,
    password: `${process.env.PANDA_OATH_TOKEN}`
  },
  channels
}
