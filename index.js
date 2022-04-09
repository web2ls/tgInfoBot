const { Telegraf, Markup } = require('telegraf');
require('dotenv').config();

const bot = new Telegraf(process.env.BOT_TOKEN)

bot.start((ctx) => {
  ctx.reply('Select ', Markup
    .keyboard(['Точное московское время'])
    .oneTime()
    .resize()
  )
});

bot.command('time', ctx => {
  ctx.reply('Moscow time is here');
});

bot.on('message', (ctx) => ctx.reply('message'));

bot.launch();

// Enable graceful stop
process.once('SIGINT', () => bot.stop('SIGINT'));
process.once('SIGTERM', () => bot.stop('SIGTERM'));