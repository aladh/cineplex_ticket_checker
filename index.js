const { get } = require('https');
// noinspection NpmUsedModulesInstalled
const AWS = require('aws-sdk');
// noinspection JSUnresolvedFunction
const ses = new AWS.SES();

const BASE_URL = 'https://www.cineplex.com/Movie/';
const EMAIL_ADDRESS = '';
const THEATRE_IDS = [];
const MOVIES = [];

function fetch(url) {
  return new Promise((resolve, reject) => {
    get(url, (res) => {
      const { statusCode } = res;

      if (statusCode === 200) {
        let data = '';
        res.on('data', (chunk) => data += chunk);
        res.on('end', () => resolve(data))
      } else {
        reject(`Status code: ${statusCode}`);
        res.destroy()
      }
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
    if (err) throw err;
  })
}

async function checkAvailability(movie) {
  const html = await fetch(BASE_URL + movie);
  const available =
    THEATRE_IDS
      .map((id) => html.includes(id))
      .includes(true);

  console.log(movie + ' ' + available);

  return available
}

exports.handler = () => {
  MOVIES.forEach(async (movie) => {
    try {
      const available = await checkAvailability(movie);
      available && sendEmail({ to: EMAIL_ADDRESS, subject: `Cineplex tickets available for ${movie}`, body: `${BASE_URL}${movie}` })
    } catch (e) {
      sendEmail({ to: EMAIL_ADDRESS, subject: `Cineplex availability checker failed for ${movie}`, body: `${e} ${BASE_URL}${movie}` })
    }
  })
};