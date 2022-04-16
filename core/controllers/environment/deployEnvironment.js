"use strict";
const { handleRequest } = require('../helpers');
const environmentService = require('../../db/models/environment/environment.service');

async function deployEnvironment(req, res) {

  const extractOutput = async (outputs) => (outputs);

  const handle = async () => {
    const { accountId, userId, environmentName, provider, credentials, headers } = res.locals;
    const { username } = req.user;
    return await environmentService.deployEnvironment(accountId, userId, environmentName, provider, credentials, headers, username)
  }
  await handleRequest({ req, res, extractOutput, handle });
}

exports.handler = deployEnvironment;
