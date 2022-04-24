module.exports = {
    hello: (client, channel, tags) => {
        client.say(channel, `Hi ${tags['display-name']}!`)
    },
    dice: (client, channel, tags, args) => {
        const diceSize = parseInt(args[0]) || 6
        const result = Math.ceil(Math.random() * diceSize)
        client.say(channel, `${tags['display-name']} rolled a ${diceSize} sided dice, and got ${result}`)
    },
    discord: (client, channel) => {
        client.say(channel, 'Need a place to squee over things and hang out with super awesome cool people? Then The Realm of Squees is the place for you! https://discord.gg/HJjfKwtcZu')
    },
    twitter: (client, channel) => {
        client.say(channel, 'Jenness Games Community Twitter Account: https://twitter.com/JennessGames')
    },
    specs: (client, channel) => {
        client.say(channel, 'Curious what Jenness is rocking for her setup? Be sure to check out the list in her panels or you can check out her PC Partpicker at: https://pcpartpicker.com/list/zhrHfP')
    },
    brb: (client, channel) => {
        client.say(channel, 'Hey y\'all! Jenness will brb. Be sure to stay hydrated and stretch in between sitting for long amounts of time! We want to stay healthy around here! Much love! oxox')
    },
}