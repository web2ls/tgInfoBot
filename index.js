const { Telegraf } = require('telegraf');
require('dotenv').config();

const bot = new Telegraf(process.env.BOT_TOKEN)

bot.start((ctx) => {
  let message = `I can help you!`;
  ctx.reply(message);
});

bot.command('time', ctx => {
  ctx.reply('Moscow time is here');
})

bot.launch();