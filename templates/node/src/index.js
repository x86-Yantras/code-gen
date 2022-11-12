const express = require("express");
const newLogger = require("./lib/logger/logger");
const requestId = require("./lib/middlewares/requestId/requestId");
const cors = require("cors");
const requestLogger = require("./lib/middlewares/logRequest/logRequest");

const userAPI = require("./services/user/adapters/api/http/routes");

const app = express();

const appName = dapp;

app.use(express.json());
app.use(requestId);
app.use(cors());

const logger = newLogger(appName);

app.use(requestLogger(logger));

apiPort = process.env.API_PORT || 3000;

const connectDB = require("./lib/db/mongo");

connectDB().catch((err) => logger.error(err));

app.use("/users", userAPI());

app.listen(apiPort, () => {
  logger.info(`api started at ${apiPort}`);
});
