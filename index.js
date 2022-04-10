const { Telegraf, Markup } = require('telegraf');
require('dotenv').config();
const axios = require('axios');

const bot = new Telegraf(process.env.BOT_TOKEN);

const TIME_URL = 'https://yandex.com/time/sync.json?geo=213';
const COMMANDS = {
  'Точное московское время': getTime
};

function getTime() {
  return new Promise((resolve, reject) => {
    axios.get(TIME_URL).then(res => {
      const time = res.data && res.data.time || null;
      if (!time) reject('We have an error in process');

      resolve(new Date(res.data.time).toLocaleTimeString());
    }).catch(err => {
      console.log(err);
      reject('We have an error in process');
    });
  })
};

function getMenu() {
  return Markup.keyboard(['Точное московское время']);
};

bot.start((ctx) => {
  ctx.reply('Select ', getMenu());
})

bot.command('time', ctx => {
  ctx.reply('Moscow time is here');
});

bot.on('message', async (ctx) => {
  if (!COMMANDS[ctx.message.text]) {
    ctx.reply('404 Not Found');
    return;
  };

  ctx.reply(await COMMANDS[ctx.message.text]());
})

bot.launch();

// Enable graceful stop
process.once('SIGINT', () => bot.stop('SIGINT'));
process.once('SIGTERM', () => bot.stop('SIGTERM'));