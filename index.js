const https = require('https');
const AWS = require('aws-sdk');
const ses = new AWS.SES();

const BASE_URL = 'https://www.cineplex.com/Movie/';
const EMAIL_ADDRESS = '';
const THEATRE_IDS = [];
const MOVIES = [];

function fetch(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, resp => {
        const { statusCode } = resp;
        if (statusCode !== 200) reject(`Status code: ${statusCode}`);

        let data = '';
        resp.on('data', chunk => data += chunk);
        resp.on('end', () => resolve(data));
      })
      .on('error', reject)
  })
}

function sendEmail({ to, subject, body }) {
  // noinspection JSUnusedLocalSymbols
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
  });
}

function checkAvailability(movie) {
  return new Promise((resolve, reject) => {
    fetch(BASE_URL + movie)
      .then((html) => {
        let available =
          THEATRE_IDS
            .map(id => html.includes(id))
            .includes(true);

        console.log(movie + ' ' + available);

        resolve(available)
      })
      .catch(reject)
  })
}

exports.handler = () => {
  MOVIES.forEach(movie => {
    checkAvailability(movie)
      .then(available => {
        if (available) {
          sendEmail({
            to: EMAIL_ADDRESS,
            subject: `Cineplex tickets available for ${movie}`,
            body: `${BASE_URL}${movie}`
          })
        }
      })
      .catch(e => {
        sendEmail({
          to: EMAIL_ADDRESS,
          subject: `Cineplex availability checker failed for ${movie}`,
          body: `${e} ${BASE_URL}${movie}`
        })
      })
  })
};