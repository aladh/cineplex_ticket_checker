const { get } = require('https');
// noinspection NpmUsedModulesInstalled
const AWS = require('aws-sdk');
// noinspection JSUnresolvedFunction
const ses = new AWS.SES();

const BASE_URL = 'https://www.cineplex.com/Movie/';
const THEATRE_IDS = JSON.parse(process.env.THEATRE_IDS);
const MOVIES = [];

function fetch(url) {
  return new Promise((resolve, reject) => {
    get(url, (res) => {
      const { statusCode } = res;

      if (statusCode !== 200) {
        reject(`Bad status code: ${statusCode}`);
        res.destroy();
        return
      }

      let data = '';
      res.on('data', (chunk) => data += chunk);
      res.on('end', () => resolve(data))
    })
      .on('error', reject)
  })
}

function sendEmail({ to, subject, body }) {
  // noinspection JSUnresolvedFunction, JSUnusedLocalSymbols
  ses.sendEmail({
    Source: to,
    Destination: { ToAddresses: [to] },
    Message: {
      Subject: {
        Data: subject
      },
      Body: {
        Text: {
          Data: body
        }
      }
    }
  }, (err, _data) => {
    if (err) throw err
  })
}

async function isAvailable(movie) {
  const html = await fetch(BASE_URL + movie);

  return THEATRE_IDS
    .map((id) => html.includes(id))
    .includes(true)
}

exports.handler = () => {
  MOVIES.forEach(async (movie) => {
    try {
      await isAvailable(movie) && sendEmail({
        to: process.env.EMAIL_ADDRESS,
        subject: `Cineplex tickets available for ${movie}`,
        body: `${BASE_URL}${movie}`
      })
    } catch (e) {
      sendEmail({
        to: process.env.EMAIL_ADDRESS,
        subject: `Cineplex availability checker failed for ${movie}`,
        body: `Error: ${e} for ${BASE_URL}${movie}`
      })
    }
  })
};