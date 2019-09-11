var express = require('express');
var app = express();

var swaggerUi = require('swagger-ui-express');
var swaggerDocument = require('./swagger.json');

var EnigmaController = require('./EnigmaController');

app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocument));
app.use('/enigma', EnigmaController);

module.exports = app;
