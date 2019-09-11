var express = require('express');
var router = express.Router();
var bodyParser = require('body-parser');

router.use(bodyParser.urlencoded({
    extended: true
}));
router.use(bodyParser.json());

var RTA_IRDA = require("./FabricHelper");

// Request requestVehicleInfo
router.post('/requestVehicleInfo', function (req, res) {
    RTA_IRDA.requestVehicleInfo(req, res);
});

// Issue requestVehicleInfo
router.post('/responseVehicleInfo', function (req, res) {
    RTA_IRDA.responseVehicleInfo(req, res);
});

// Get VehicleInfo history
router.post('/getVehicleHistory', function (req, res) {
    RTA_IRDA.getVehicleHistory(req, res);
});

module.exports = router;