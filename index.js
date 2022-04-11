const { Telegraf, Markup } = require('telegraf');
require('dotenv').config();
const axios = require('axios');

const bot = new Telegraf(process.env.BOT_TOKEN);

const TIME_URL = 'https://yandex.com/time/sync.json?geo=213';
const CURRENCY_URL = 'https://www.cbr-xml-daily.ru/daily_json.js';
const COMMANDS = {
  'Точное московское время': getTime,
  'Курсы валют': getCurrencyRates,
};

function getTime() {
  return new Promise((resolve, reject) => {
    axios.get(TIME_URL).then(res => {
      const time = res.data && res.data.time || null;
      if (!time) reject('We have an error in process');

      resolve(new Date(res.data.time).toLocaleTimeString());
    }).catch(err => {
      reject(false);
    });
  })
};

function getCurrencyRates() {
  return new Promise(async (resolve, reject) => {
    try {
      const currency = await getCurrency();
      const USD = currency.Valute.USD.Value;
      const EUR = currency.Valute.EUR.Value;

      resolve(`USD: ${USD.toFixed(2)} , EUR: ${EUR.toFixed(2)}`);
    } catch (err) {
      reject(false);
    }
  })
};

function getCurrency() {
  return new Promise((resolve, reject) => {
    axios.get(CURRENCY_URL).then(res => {
      resolve(res.data);
    }).catch(err => {
      reject(false);
    })
  })
}

function getMenu() {
  return Markup.keyboard(['Точное московское время', 'Курсы валют']);
};

bot.start((ctx) => {
  ctx.reply('Select ', getMenu());
})

bot.command('time', ctx => {
  ctx.reply('Moscow time is here');
});

bot.on('message', async (ctx) => {
  if (!COMMANDS[ctx.message.text]) {
    ctx.reply('Sorry, but i dont know what you want...');
    return;
  };

  try {
    ctx.reply(await COMMANDS[ctx.message.text]());
  } catch (err) {
    ctx.reply('We have an error in process');
  }
})

bot.launch();

// Enable graceful stop
process.once('SIGINT', () => bot.stop('SIGINT'));
process.once('SIGTERM', () => bot.stop('SIGTERM'));