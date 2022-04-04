const { Telegraf } = require('telegraf');
require('dotenv').config();

const bot = new Telegraf(process.env.BOT_TOKEN)

bot.start((ctx) => {
  let message = `Lets drop my friend`
  ctx.reply(message)
});