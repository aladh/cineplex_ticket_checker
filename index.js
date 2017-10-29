const https = require('https');
const AWS = require('aws-sdk');
const ses = new AWS.SES();

const BASE_URL = 'https://www.cineplex.com/Movie/';
const EMAIL_ADDRESS = '';
const MOVIES = [];

function getResponse(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (resp) => {
        const { statusCode } = resp;
        let data = '';

        if (statusCode !== 200) reject(`Status code: ${statusCode}`);
        resp.on('data', chunk => data += chunk);
        resp.on('end', () => resolve(data));
      })
      .on('error', reject)
  })
}

function sendEmail({ subject, body }) {
  const to = [EMAIL_ADDRESS];

  ses.sendEmail({
      Source: to[0],
      Destination: { ToAddresses: to },
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
    }
    , function (err, _data) {
      if (err) {
        throw err
      }

      console.log('Email sent:');
    });
}

function ticketsAvailable(movie) {
  return new Promise((resolve, reject) => {
    console.log('Checking ' + movie);

    getResponse(BASE_URL + movie)
      .then((html) => {
        let onSale = html.includes('quick-tickets');
        console.log(movie + ' ' + onSale);

        resolve(onSale)
      })
      .catch((err) => {
        console.log(err);

        reject(err)
      })
  })
}

function emailIfAvailable(movies) {
  return new Promise((_, reject) => {
    movies.forEach((movie) => {
      ticketsAvailable(movie)
        .then((available) => {
          if (available) {
            sendEmail({
              subject: `Cineplex tickets on sale for ${movie}`,
              body: `${BASE_URL}${movie}`
            })
          }
        })
        .catch((err) => {
          sendEmail({
            subject: `Cineplex checker failed for ${movie}`,
            body: `${err} ${BASE_URL}${movie}`
          });

          reject(err)
        })
    })
  })
}

exports.handler = (event, context, callback) => {
  emailIfAvailable(MOVIES)
    .catch(callback)
};