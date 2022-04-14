require('dotenv').config()

module.exports = {
  options: { debug: true },
  identity: {
    username: process.env.PANDA_USERNAME,
    password: `oauth:${process.env.PANDA_OATH_TOKEN}`
  },
  channels: ['sephory']
}
